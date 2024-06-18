package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-cloudfoundry/internal/provider/managers"
	cfv3client "github.com/cloudfoundry/go-cfclient/v3/client"
	cfv3resource "github.com/cloudfoundry/go-cfclient/v3/resource"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

type serviceRouteBindingResource struct {
	cfClient *cfv3client.Client
}

var (
	_ resource.ResourceWithConfigure   = &serviceRouteBindingResource{}
	_ resource.ResourceWithImportState = &serviceRouteBindingResource{}
)

func NewServiceRouteBindingResource() resource.Resource {
	return &serviceRouteBindingResource{}
}

func (r *serviceRouteBindingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_route_binding"
}

func (r *serviceRouteBindingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Service route bindings are relations between a service instance and a route.

Not all service instances support route binding. In order to bind to a managed service instance, the service instance should be created from a service offering that has requires route forwarding (requires=[route_forwarding]). In order to bind to a user-provided service instance, the service instance must have route_service_url set.`,

		Attributes: map[string]schema.Attribute{
			"service_instance": schema.StringAttribute{
				MarkdownDescription: "The GUID of the service instance to be bound",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"route": schema.StringAttribute{
				MarkdownDescription: "The GUID of the route to be bound",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"parameters": schema.StringAttribute{
				MarkdownDescription: "A JSON object that is passed to the service broker for managed service instance.",
				Optional:            true,
				CustomType:          jsontypes.NormalizedType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"route_service_url": schema.StringAttribute{
				MarkdownDescription: "URL to which requests for bound routes will be forwarded.",
				Computed:            true,
			},
			"last_operation": lastOperationSchema(),
			idKey:            guidSchema(),
			labelsKey:        resourceLabelsSchema(),
			annotationsKey:   resourceAnnotationsSchema(),
			createdAtKey:     createdAtSchema(),
			updatedAtKey:     updatedAtSchema(),
		},
	}
}

func (r *serviceRouteBindingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	session, ok := req.ProviderData.(*managers.Session)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *managers.Session, got: %T. Please report this issue to the provider developers", req.ProviderData),
		)
		return
	}
	r.cfClient = session.CFClient
}

func (r *serviceRouteBindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var (
		plan                serviceRouteBindingType
		serviceRouteBinding *cfv3resource.ServiceRouteBinding
		err                 error
	)
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createServiceRouteBinding, diags := plan.mapCreateServiceRouteBindingTypeToValues(ctx)
	resp.Diagnostics.Append(diags...)

	jobID, serviceRouteBinding, err := r.cfClient.ServiceRouteBindings.Create(ctx, &createServiceRouteBinding)
	if err != nil {
		resp.Diagnostics.AddError(
			"API Error in creating service Route Binding",
			"Unable to create Route Binding with Route "+plan.Route.ValueString()+": "+err.Error(),
		)
		return

	} else if jobID != "" {
		err = pollJob(ctx, *r.cfClient, jobID, defaultTimeout)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to verify service Route binding creation",
				"Service Route Binding verification failed for Route "+plan.Route.ValueString()+": "+err.Error(),
			)
			return
		}

		getOptions := cfv3client.ServiceRouteBindingListOptions{
			ServiceInstanceGUIDs: cfv3client.Filter{
				Values: []string{
					plan.ServiceInstance.ValueString(),
				},
			},
			RouteGUIDs: cfv3client.Filter{
				Values: []string{
					plan.Route.ValueString(),
				},
			},
		}

		serviceRouteBinding, err = r.cfClient.ServiceRouteBindings.Single(ctx, &getOptions)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error fetching service route binding after creation",
				"Unable to fetch created service route binding with route "+plan.Route.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	data, diags := mapServiceRouteBindingValuesToType(ctx, serviceRouteBinding)
	data.Parameters = plan.Parameters
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serviceRouteBindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data serviceRouteBindingType

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	serviceRouteBinding, err := r.cfClient.ServiceRouteBindings.Get(ctx, data.ID.ValueString())
	if err != nil {
		handleReadErrors(ctx, resp, err, "service_route_binding", data.ID.ValueString())
		return
	}

	state, diags := mapServiceRouteBindingValuesToType(ctx, serviceRouteBinding)
	state.Parameters = data.Parameters
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}

func (r *serviceRouteBindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, previousState serviceRouteBindingType
	var diags diag.Diagnostics
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &previousState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateServiceRouteBinding, diags := plan.mapUpdateServiceRouteBindingTypeToValues(ctx, previousState)
	resp.Diagnostics.Append(diags...)

	serviceRouteBinding, err := r.cfClient.ServiceRouteBindings.Update(ctx, plan.ID.ValueString(), &updateServiceRouteBinding)
	if err != nil {
		resp.Diagnostics.AddError(
			"API Error Updating Service Route Binding",
			"Could not update Service Route Binding with ID "+plan.ID.ValueString()+" : "+err.Error(),
		)
		return
	}

	data, diags := mapServiceRouteBindingValuesToType(ctx, serviceRouteBinding)
	data.Parameters = plan.Parameters
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serviceRouteBindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state serviceRouteBindingType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	jobID, err := r.cfClient.ServiceRouteBindings.Delete(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"API Error in deleting service route binding",
			"Unable to delete route binding "+state.ID.ValueString()+": "+err.Error(),
		)

	}
	if jobID != "" {
		if err := pollJob(ctx, *r.cfClient, jobID, defaultTimeout); err != nil {
			resp.Diagnostics.AddError(
				"Unable to verify service route binding deletion",
				"service route binding deletion verification failed for "+state.ID.ValueString()+": "+err.Error(),
			)
		}
	}

}

func (rs *serviceRouteBindingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
