package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-cloudfoundry/internal/provider/managers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/samber/lo"

	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SpaceDataSource{}
var _ datasource.DataSourceWithConfigure = &SpaceDataSource{}

func NewSpaceDataSource() datasource.DataSource {
	return &SpaceDataSource{}
}

type SpaceDataSource struct {
	cfClient *client.Client
}

type SpaceDataSourceModel struct {
	Name        types.String `tfsdk:"name"`
	Id          types.String `tfsdk:"id"`
	OrgName     types.String `tfsdk:"org_name"`
	Org         types.String `tfsdk:"org"`
	Quota       types.String `tfsdk:"quota"`
	Labels      types.Map    `tfsdk:"labels"`
	Annotations types.Map    `tfsdk:"annotations"`
}

func (d *SpaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space"
}

func (d *SpaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Gets information on a Cloud Foundry space.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the space to look up",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The GUID of the space",
				Computed:            true,
			},
			"org_name": schema.StringAttribute{
				MarkdownDescription: "The name of the organization under which the space exists",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("org"),
					}...),
				},
			},
			"org": schema.StringAttribute{
				MarkdownDescription: "The GUID of the organization under which the space exists",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("org_name"),
					}...),
				},
			},
			"quota": schema.StringAttribute{
				MarkdownDescription: "The space quota applied to the space",
				Computed:            true,
			},
			labelsKey:      labelsSchema(),
			annotationsKey: annotationsSchema(),
		},
	}
}

func (d *SpaceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SpaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SpaceDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Ensure org details is present in state file else populate org name or guid accordingly
	if data.Org.IsNull() && data.OrgName.IsNull() {
		resp.Diagnostics.AddError(
			"Neither Org GUID nor Org Name is present.",
			"Expected either 'org' or 'org_name' attribute.",
		)
		return

	} else if data.Org.IsNull() {
		orgs, err := d.cfClient.Organizations.ListAll(ctx, &client.OrganizationListOptions{
			Names: client.Filter{
				Values: []string{
					data.OrgName.ValueString(),
				},
			},
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to fetch org data.",
				fmt.Sprintf("Request failed with %s.", err.Error()),
			)
			return
		}

		org, found := lo.Find(orgs, func(org *resource.Organization) bool {
			return org.Name == data.OrgName.ValueString()
		})

		if !found {
			resp.Diagnostics.AddError(
				"Unable to find org data in list",
				fmt.Sprintf("Given name %s not in the list of orgs.", data.OrgName.ValueString()),
			)
			return
		}
		data.Org = types.StringValue(org.GUID)
	} else {

		//Fetching organization with GUID
		org, err := d.cfClient.Organizations.Get(ctx, data.Org.ValueString())
		if err != nil {
			switch err.(type) {
			case resource.CloudFoundryError:
				resp.Diagnostics.AddError(
					"Unable to find org data in list.",
					fmt.Sprintf("Given org %s not in the list of orgs.", data.Org.ValueString()),
				)
			default:
				resp.Diagnostics.AddError(
					"Unable to fetch org data.",
					fmt.Sprintf("Request failed with %s.", err.Error()),
				)
			}

			return
		}

		data.OrgName = types.StringValue(org.Name)
	}

	//Filtering for spaces under the org with GUID
	spaces, err := d.cfClient.Spaces.ListAll(ctx, &client.SpaceListOptions{
		OrganizationGUIDs: client.Filter{
			Values: []string{
				data.Org.ValueString(),
			},
		},
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch space data.",
			fmt.Sprintf("Request failed with %s.", err.Error()),
		)
		return
	}

	//No spaces with name might be present under org
	space, found := lo.Find(spaces, func(space *resource.Space) bool {
		return space.Name == data.Name.ValueString()
	})
	if !found {
		resp.Diagnostics.AddError(
			"Unable to find space data in list",
			fmt.Sprintf("Given name %s not in the list of spaces under %s.", data.Name.ValueString(), data.OrgName.ValueString()),
		)
		return
	}

	data.Id = types.StringValue(space.GUID)
	data.Labels = *setMapToBaseMap(ctx, resp, space.Metadata.Labels)
	data.Annotations = *setMapToBaseMap(ctx, resp, space.Metadata.Annotations)

	//Checking if quota exists, then taking the guid value
	if space.Relationships.Quota.Data != nil {
		data.Quota = types.StringValue(space.Relationships.Quota.Data.GUID)
	} else {
		data.Quota = types.StringValue("")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "read a space data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}
