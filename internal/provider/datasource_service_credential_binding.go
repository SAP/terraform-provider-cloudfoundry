package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-cloudfoundry/internal/provider/managers"
	cfv3client "github.com/cloudfoundry-community/go-cfclient/v3/client"
	cfv3resource "github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/samber/lo"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource                   = &ServiceCredentialBindingDataSource{}
	_ datasource.DataSourceWithConfigure      = &ServiceCredentialBindingDataSource{}
	_ datasource.DataSourceWithValidateConfig = &ServiceCredentialBindingDataSource{}
)

func NewServiceCredentialBindingDataSource() datasource.DataSource {
	return &ServiceCredentialBindingDataSource{}
}

type ServiceCredentialBindingDataSource struct {
	cfClient *cfv3client.Client
}

func (d *ServiceCredentialBindingDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_credential_binding"
}

func (d *ServiceCredentialBindingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Gets information of a Service Credential Binding.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the service credential binding when the type is app",
				Optional:            true,
			},
			"service_instance": schema.StringAttribute{
				MarkdownDescription: "The GUID of the service instance that this binding is originated from",
				Required:            true,
			},
			"app": schema.StringAttribute{
				MarkdownDescription: "The GUID of the app using this binding for key bindings",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the service credential binding. Either app or key.",
				Required:            true,
			},
			"last_operation": schema.SingleNestedAttribute{
				MarkdownDescription: "The last operation performed on the service credential binding",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the last operation",
						Computed:            true,
					},
					"state": schema.StringAttribute{
						MarkdownDescription: "The state of the last operation",
						Computed:            true,
					},
					"description": schema.StringAttribute{
						MarkdownDescription: "A description of the last operation",
						Computed:            true,
					},
					"updated_at": schema.StringAttribute{
						MarkdownDescription: "The time at which the last operation was updated",
						Computed:            true,
					},
					"created_at": schema.StringAttribute{
						MarkdownDescription: "The time at which the last operation was created",
						Computed:            true,
					},
				},
			},
			idKey:          guidSchema(),
			labelsKey:      datasourceLabelsSchema(),
			annotationsKey: datasourceAnnotationsSchema(),
			createdAtKey:   createdAtSchema(),
			updatedAtKey:   updatedAtSchema(),
		},
	}
}

func (d *ServiceCredentialBindingDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	session, ok := req.ProviderData.(*managers.Session)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *managers.Session, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.cfClient = session.CFClient
}

func (d *ServiceCredentialBindingDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
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

func (d *ServiceCredentialBindingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data serviceCredentialBindingType

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	svcCredentialBindings, err := d.cfClient.ServiceCredentialBindings.ListAll(ctx, &cfv3client.ServiceCredentialBindingListOptions{
		Names: cfv3client.Filter{
			Values: []string{data.Name.ValueString()},
		},
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"API Error Fetching Service Credential Binding.",
			fmt.Sprintf("Request failed with %s.", err.Error()),
		)
		return
	}
	svcCredentialBinding, found := lo.Find(svcCredentialBindings, func(svcCredentialBinding *cfv3resource.ServiceCredentialBinding) bool {
		return svcCredentialBinding.Name == data.Name.ValueString()
	})
	if !found {
		resp.Diagnostics.AddError(
			"Unable to find service credential binding in list",
			fmt.Sprintf("Given name %s not in the list of service instances.", data.Name.ValueString()),
		)
		return
	}

	data, diags = mapServiceCredentialBindingValuesToType(ctx, svcCredentialBinding)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
