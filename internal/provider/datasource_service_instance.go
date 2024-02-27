package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-cloudfoundry/internal/provider/managers"
	cfv3client "github.com/cloudfoundry-community/go-cfclient/v3/client"
	cfv3resource "github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/samber/lo"
)

func NewServiceInstanceDataSource() datasource.DataSource {
	return &ServiceInstanceDataSource{}
}

type ServiceInstanceDataSource struct {
	cfClient *cfv3client.Client
}

func (d *ServiceInstanceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_instance"
}

func (d *ServiceInstanceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Gets information of a Service Instance.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the service instance to look up",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the service instnace. Either managed or user-provided.",
				Computed:            true,
			},
			"space": schema.StringAttribute{
				MarkdownDescription: "The ID of the space in which to create the service instance",
				Computed:            true,
			},
			"service_plan": schema.StringAttribute{
				MarkdownDescription: "The ID of the service plan from which to create the service instance",
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "Set of tags used by apps to identify service instances. They are shown in the app VCAP_SERVICES env.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"syslog_drain_url": schema.StringAttribute{
				MarkdownDescription: "URL to which logs for bound applications will be streamed; only shown when type is user-provided.",
				Computed:            true,
			},
			"route_service_url": schema.StringAttribute{
				MarkdownDescription: "URL to which requests for bound routes will be forwarded; only shown when type is user-provided.",
				Computed:            true,
			},
			"maintenance_info": schema.MapAttribute{
				MarkdownDescription: "Information about the version of this service instance; only shown when type is managed",
				ElementType:         maintenanceInfoAttrTypes,
				Computed:            true,
			},
			"upgrade_available": schema.BoolAttribute{
				MarkdownDescription: "Whether or not an upgrade of this service instance is available on the current Service Plan; details are available in the maintenance_info object; Only shown when type is managed",
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
			labelsKey:      datasourceLabelsSchema(),
			annotationsKey: datasourceAnnotationsSchema(),
			createdAtKey:   createdAtSchema(),
			updatedAtKey:   updatedAtSchema(),
		},
	}
}

func (d *ServiceInstanceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServiceInstanceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data datasourceServiceInstanceType

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	svcInstances, err := d.cfClient.ServiceInstances.ListAll(ctx, &cfv3client.ServiceInstanceListOptions{
		Names: cfv3client.Filter{
			Values: []string{
				data.Name.ValueString(),
			},
		},
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"API Error to fetch service instance data.",
			fmt.Sprintf("Request failed with %s.", err.Error()),
		)
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

	switch svcInstance.Type {
	case "managed":
		data, diags = mapDataSourceServiceInstanceValuesToType(ctx, svcInstance)
	case "user-provided":
		data, diags = mapDataSourceServiceInstanceValuesToType(ctx, svcInstance)
	}
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
