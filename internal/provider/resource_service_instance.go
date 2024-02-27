package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SAP/terraform-provider-cloudfoundry/internal/provider/managers"
	cfv3client "github.com/cloudfoundry-community/go-cfclient/v3/client"
	cfv3resource "github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/samber/lo"
)

type serviceInstanceResource struct {
	cfClient *cfv3client.Client
}

var (
	_ resource.Resource                   = &serviceInstanceResource{}
	_ resource.ResourceWithConfigure      = &serviceInstanceResource{}
	_ resource.ResourceWithImportState    = &serviceInstanceResource{}
	_ resource.ResourceWithValidateConfig = &serviceInstanceResource{}
)

const (
	managedSerivceInstance      = "managed"
	userProvidedServiceInstance = "user-provided"
)

func NewServiceInstanceResource() resource.Resource {
	return &serviceInstanceResource{}
}

func (r *serviceInstanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_instance"
}

func (r *serviceInstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Creates a service instance in a cloudfoundry space.

__Further documentation:__
https://docs.cloudfoundry.org/devguide/services`,

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the service instance",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the service instance. Either managed or user-provided.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("managed", "user-provided"),
				},
			},
			"space": schema.StringAttribute{
				MarkdownDescription: "The ID of the space in which to create the service instance",
				Required:            true,
			},
			"service_plan": schema.StringAttribute{
				MarkdownDescription: "The ID of the service plan from which to create the service instance",
				Optional:            true,
			},
			"parameters": schema.StringAttribute{
				MarkdownDescription: "A JSON object that is passed to the service broker for managed service instance.",
				Optional:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
			"tags": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Set of tags used by apps to identify service instances. They are shown in the app VCAP_SERVICES env.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"credentials": schema.StringAttribute{
				MarkdownDescription: "A JSON object that is made available to apps bound to this service instance of type user-provided.",
				Optional:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
			"syslog_drain_url": schema.StringAttribute{
				MarkdownDescription: "URL to which logs for bound applications will be streamed; only shown when type is user-provided.",
				Optional:            true,
			},
			"route_service_url": schema.StringAttribute{
				MarkdownDescription: "URL to which requests for bound routes will be forwarded; only shown when type is user-provided.",
				Optional:            true,
			},
			"maintenance_info": schema.ListAttribute{
				MarkdownDescription: "Information about the version of this service instance; only shown when type is managed",
				ElementType:         maintenanceInfoAttrTypes,
				Optional:            true,
				Computed:            true,
			},
			"upgrade_available": schema.BoolAttribute{
				MarkdownDescription: "Whether or not an upgrade of this service instance is available on the current Service Plan; details are available in the maintenance_info object; Only shown when type is managed",
				Optional:            true,
				Computed:            true,
			},
			"dashboard_url": schema.StringAttribute{
				MarkdownDescription: "The URL to the service instance dashboard (or null if there is none); only shown when type is managed.",
				Computed:            true,
			},
			"last_operation": schema.ListAttribute{
				MarkdownDescription: "The last operation of this service instance.",
				Computed:            true,
				ElementType:         lastOperationAttrTypes,
			},
			idKey:          guidSchema(),
			labelsKey:      resourceLabelsSchema(),
			annotationsKey: resourceAnnotationsSchema(),
			createdAtKey:   createdAtSchema(),
			updatedAtKey:   updatedAtSchema(),
		},
	}
}

func (r *serviceInstanceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	session, ok := req.ProviderData.(*managers.Session)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *managers.Session, got: %T. Please report this issue to the provider developers", req.ProviderData),
		)
		return
	}
	r.cfClient = session.CFClient
}

func (r *serviceInstanceResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config serviceInstanceType
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.Type.ValueString() == userProvidedServiceInstance && !config.ServicePlan.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("service_plan"),
			"Conflicting attribute service instance",
			"Service plan is not allowed for user-provided service instance",
		)
		return
	}

	if config.Type.ValueString() == managedSerivceInstance && config.ServicePlan.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("service_plan"),
			"Missing attribute service instance",
			"Service plan is required for managed service instance",
		)
		return

	}

	// If Service Instance is of type managed only parameters is allowed to pass
	if !config.Parameters.IsNull() && config.Type.ValueString() != "managed" {
		resp.Diagnostics.AddAttributeError(
			path.Root("type"),
			"Parameters can only passed to service instance of type managed",
			"Parameters json object can only be passed to managed serivce instance",
		)
		return
	}

	// If Service instance of type user-provided then credentials , syslog_drain_url and route_service_url allowed
	if !config.SyslogDrainURL.IsNull() || !config.RouteServiceURL.IsNull() || !config.Credentials.IsNull() {
		if config.Type.ValueString() != "user-provided" {
			resp.Diagnostics.AddAttributeError(
				path.Root("type"),
				"Mistmatch attribute passed to user provided service instance",
				"Allowed attributes for serivce instance of type user provided: credentials, syslog_drain_url, route_service_url",
			)
		}
	}

}

func (r *serviceInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serviceInstanceType
	var serviceInstance *cfv3resource.ServiceInstance
	var err error
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	switch plan.Type.ValueString() {
	case managedSerivceInstance:
		createServiceInstance := cfv3resource.ServiceInstanceManagedCreate{
			Type: plan.Type.ValueString(),
			Name: plan.Name.ValueString(),
			Relationships: cfv3resource.ServiceInstanceRelationships{
				ServicePlan: &cfv3resource.ToOneRelationship{
					Data: &cfv3resource.Relationship{
						GUID: plan.ServicePlan.ValueString(),
					},
				},
				Space: &cfv3resource.ToOneRelationship{
					Data: &cfv3resource.Relationship{
						GUID: plan.Space.ValueString(),
					},
				},
			},
			Metadata: cfv3resource.NewMetadata(),
		}
		if !plan.Parameters.IsNull() {
			var params json.RawMessage
			err := json.Unmarshal([]byte(plan.Parameters.ValueString()), &params)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error in unmarshalling parameters",
					"Unable to unmarshal json parameters of service instance"+plan.Name.ValueString()+": "+err.Error(),
				)
				return
			}
			createServiceInstance.Parameters = &params
		}
		if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
			tags, diags := toTagsList(ctx, plan.Tags)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			createServiceInstance.Tags = tags
		}
		labelsDiags := plan.Labels.ElementsAs(ctx, &createServiceInstance.Metadata.Labels, false)
		resp.Diagnostics.Append(labelsDiags...)

		annotationsDiags := plan.Annotations.ElementsAs(ctx, &createServiceInstance.Metadata.Annotations, false)
		resp.Diagnostics.Append(annotationsDiags...)

		jobID, err := r.cfClient.ServiceInstances.CreateManaged(ctx, &createServiceInstance)
		if err != nil {
			resp.Diagnostics.AddError(
				"API Error in creating managed service instance",
				"Unable to create service instance "+plan.Name.ValueString()+": "+err.Error(),
			)
		}
		if pollJob(ctx, *r.cfClient, jobID) != nil {
			resp.Diagnostics.AddError(
				"Unable to verify service instance creation",
				"Service Instance verification failed for "+plan.Name.ValueString()+": "+err.Error(),
			)
		}
		serviceInstance, err = r.cfClient.ServiceInstances.Single(ctx, &cfv3client.ServiceInstanceListOptions{
			Names: cfv3client.Filter{
				Values: []string{
					plan.Name.ValueString(),
				},
			},
			SpaceGUIDs: cfv3client.Filter{
				Values: []string{
					plan.Space.ValueString(),
				},
			},
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error get service instance after creation",
				"Unable to fetch created service instance"+plan.Name.ValueString()+": "+err.Error(),
			)
		}

		plan, diags = mapResourceServiceInstanceValuesToType(ctx, serviceInstance, plan.Parameters)
		resp.Diagnostics.Append(diags...)

	case userProvidedServiceInstance:

		createServiceInstance := cfv3resource.ServiceInstanceUserProvidedCreate{
			Type: plan.Type.ValueString(),
			Name: plan.Name.ValueString(),
			Relationships: cfv3resource.ServiceInstanceRelationships{
				Space: &cfv3resource.ToOneRelationship{
					Data: &cfv3resource.Relationship{
						GUID: plan.Space.ValueString(),
					},
				},
			},
		}
		if !plan.Credentials.IsNull() {
			var credentials json.RawMessage
			err := json.Unmarshal([]byte(plan.Credentials.ValueString()), &credentials)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error in unmarshalling credentials",
					"Unable to unmarshal json credentials of service instance"+plan.Name.ValueString()+": "+err.Error(),
				)
				return
			}
			createServiceInstance.Credentials = &credentials
		}
		if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
			tags, diags := toTagsList(ctx, plan.Tags)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			createServiceInstance.Tags = tags
		}

		if !plan.SyslogDrainURL.IsNull() {
			createServiceInstance.SyslogDrainURL = plan.SyslogDrainURL.ValueStringPointer()
		}
		if !plan.RouteServiceURL.IsNull() {
			createServiceInstance.RouteServiceURL = plan.RouteServiceURL.ValueStringPointer()
		}

		_, err = r.cfClient.ServiceInstances.CreateUserProvided(ctx, &createServiceInstance)
		if err != nil {
			resp.Diagnostics.AddError(
				"API Error in creating user-provided service instance",
				"Unable to create service instance "+plan.Name.ValueString()+": "+err.Error(),
			)
		}
		serviceInstance, err = r.cfClient.ServiceInstances.Single(ctx, &cfv3client.ServiceInstanceListOptions{
			Names: cfv3client.Filter{
				Values: []string{
					plan.Name.ValueString(),
				},
			},
			SpaceGUIDs: cfv3client.Filter{
				Values: []string{
					plan.Space.ValueString(),
				},
			},
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error get service instance after creation",
				"Unable to fetch created service instance"+plan.Name.ValueString()+": "+err.Error(),
			)
		}
		plan, diags = mapResourceServiceInstanceValuesToType(ctx, serviceInstance, plan.Credentials)
		resp.Diagnostics.Append(diags...)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

}

func (r *serviceInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data serviceInstanceType

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	svcInstances, err := r.cfClient.ServiceInstances.ListAll(ctx, &cfv3client.ServiceInstanceListOptions{
		Names: cfv3client.Filter{
			Values: []string{
				data.Name.ValueString(),
			},
		},
	})
	if err != nil {
		handleReadErrors(ctx, resp, err, "service_instance", data.ID.ValueString())
		return
	}
	svcInstance, found := lo.Find(svcInstances, func(svcInstance *cfv3resource.ServiceInstance) bool {
		return svcInstance.Name == data.Name.ValueString()
	})
	if !found {
		resp.Diagnostics.AddError(
			"Unable to find service instance in list",
			fmt.Sprintf("Given name %s not in the list of service instances.", data.Name.ValueString()),
		)
		return
	}
	data, diags = mapResourceServiceInstanceValuesToType(ctx, svcInstance, data.Parameters)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *serviceInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan, previousState serviceInstanceType
	var diags diag.Diagnostics
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &previousState)...)

	switch plan.Type.ValueString() {
	case managedSerivceInstance:

		updateServiceInstance := cfv3resource.ServiceInstanceManagedUpdate{
			Name: plan.Name.ValueStringPointer(),
		}
		// Check if the service plan is different from the previous state
		if plan.ServicePlan.ValueString() != previousState.ServicePlan.ValueString() {
			ok, err := isServiceInstanceUpgradable(ctx, previousState.ID.ValueString(), *r.cfClient)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error in checking service instance upgradability",
					"Unable to check service instance upgradability"+plan.Name.ValueString()+": "+err.Error(),
				)
				return
			}
			if !ok {
				resp.Diagnostics.AddError(
					"Service instance not upgradable",
					"Service instance "+plan.Name.ValueString()+" is not upgradable",
				)
				return
			}
			updateServiceInstance.Relationships = &cfv3resource.ServiceInstanceRelationships{
				ServicePlan: &cfv3resource.ToOneRelationship{
					Data: &cfv3resource.Relationship{
						GUID: plan.ServicePlan.ValueString(),
					},
				},
			}
		}
		if !plan.Parameters.IsNull() {
			var params json.RawMessage
			err := json.Unmarshal([]byte(plan.Parameters.ValueString()), &params)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error in unmarshalling parameters",
					"Unable to unmarshal json parameters during update of service instance"+plan.Name.ValueString()+": "+err.Error(),
				)
				return
			}
			updateServiceInstance.Parameters = &params
		}
		// check if tag is not null and update the tags
		if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
			tags, diags := toTagsList(ctx, plan.Tags)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			updateServiceInstance.Tags = tags
		}

		updateServiceInstance.Metadata, diags = setClientMetadataForUpdate(ctx, previousState.Labels, previousState.Annotations, plan.Labels, plan.Annotations)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		jobID, _, err := r.cfClient.ServiceInstances.UpdateManaged(ctx, previousState.ID.ValueString(), &updateServiceInstance)
		if err != nil {
			resp.Diagnostics.AddError(
				"API Error in updating managed service instance",
				"Unable to update service instance "+plan.Name.ValueString()+": "+err.Error(),
			)
		}
		if jobID != "" {
			if pollJob(ctx, *r.cfClient, jobID) != nil {
				resp.Diagnostics.AddError(
					"Unable to verify service instance update",
					"Service Instance update verification failed for "+plan.Name.ValueString()+": "+err.Error(),
				)
			}
		}
		serviceInstance, err := r.cfClient.ServiceInstances.Get(ctx, plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error get service instance after update",
				"Unable to fetch updated service instance"+plan.Name.ValueString()+": "+err.Error(),
			)
		}
		plan, diags = mapResourceServiceInstanceValuesToType(ctx, serviceInstance, plan.Parameters)
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	case userProvidedServiceInstance:
		updateServiceInstance := cfv3resource.ServiceInstanceUserProvidedUpdate{
			Name: plan.Name.ValueStringPointer(),
		}
		if !plan.Credentials.IsNull() {
			var credentials json.RawMessage
			err := json.Unmarshal([]byte(plan.Credentials.ValueString()), &credentials)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error in unmarshalling credentials",
					"Unable to unmarshal json credentials during update of service instance"+plan.Name.ValueString()+": "+err.Error(),
				)
				return
			}
			updateServiceInstance.Credentials = &credentials
		}
		if !plan.SyslogDrainURL.IsNull() {
			updateServiceInstance.SyslogDrainURL = plan.SyslogDrainURL.ValueStringPointer()
		}
		if !plan.RouteServiceURL.IsNull() {
			updateServiceInstance.RouteServiceURL = plan.RouteServiceURL.ValueStringPointer()
		}
		if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
			tags, diags := toTagsList(ctx, plan.Tags)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			updateServiceInstance.Tags = tags
		}
		updateServiceInstance.Metadata, diags = setClientMetadataForUpdate(ctx, previousState.Labels, previousState.Annotations, plan.Labels, plan.Annotations)
		resp.Diagnostics.Append(diags...)
		_, err := r.cfClient.ServiceInstances.UpdateUserProvided(ctx, previousState.ID.ValueString(), &updateServiceInstance)
		if err != nil {
			resp.Diagnostics.AddError(
				"API Error in updating user-provided service instance",
				"Unable to update service instance "+plan.Name.ValueString()+": "+err.Error(),
			)
		}
		serviceInstance, err := r.cfClient.ServiceInstances.Get(ctx, plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error get service instance after update",
				"Unable to fetch updated service instance"+plan.Name.ValueString()+": "+err.Error(),
			)
		}
		plan, diags = mapResourceServiceInstanceValuesToType(ctx, serviceInstance, plan.Credentials)
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	}

}

func (r *serviceInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state serviceInstanceType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	jobID, err := r.cfClient.ServiceInstances.Delete(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"API Error in deleting service instance",
			"Unable to delete service instance "+state.Name.ValueString()+": "+err.Error(),
		)

	}
	if jobID != "" {
		if pollJob(ctx, *r.cfClient, jobID) != nil {
			resp.Diagnostics.AddError(
				"Unable to verify service instance deletion",
				"Service Instance deletion verification failed for "+state.ID.ValueString()+": "+err.Error(),
			)
		}
	}

}

func (rs *serviceInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
