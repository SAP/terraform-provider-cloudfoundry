package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-cloudfoundry/internal/provider/managers"
	"github.com/SAP/terraform-provider-cloudfoundry/internal/validation"
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
var (
	_ datasource.DataSource              = &SpaceDataSource{}
	_ datasource.DataSourceWithConfigure = &SpaceDataSource{}
)

func NewSpaceDataSource() datasource.DataSource {
	return &SpaceDataSource{}
}

type SpaceDataSource struct {
	cfClient *client.Client
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
					validation.ValidUUID(),
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("org_name"),
					}...),
				},
			},
			"quota": schema.StringAttribute{
				MarkdownDescription: "The GUID of the space's quota",
				Computed:            true,
			},
			"allow_ssh": schema.BoolAttribute{
				MarkdownDescription: "Allows SSH to application containers via the CF CLI.",
				Computed:            true,
			},
			"isolation_segment": schema.StringAttribute{
				MarkdownDescription: "The ID of the isolation segment assigned to the space.",
				Computed:            true,
			},
			"asgs": schema.SetAttribute{
				MarkdownDescription: "List of running application security groups applied to applications running within this space.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"staging_asgs": schema.SetAttribute{
				MarkdownDescription: "List of staging application security groups applied to applications being staged for this space.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was updated in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
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
	var data spaceType

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = data.populateOrgValues(ctx, d.cfClient)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Filtering for spaces under the org with GUID
	spaces, err := d.cfClient.Spaces.ListAll(ctx, &client.SpaceListOptions{
		OrganizationGUIDs: client.Filter{
			Values: []string{
				data.OrgId.ValueString(),
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

	diags = data.setTypeValuesFromSpace(ctx, space)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sshEnabled, err := d.cfClient.SpaceFeatures.IsSSHEnabled(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch space features.",
			fmt.Sprintf("Request failed with %s.", err.Error()),
		)
		return
	}

	data.setTypeValueFromBool(sshEnabled)

	isolationSegment, err := d.cfClient.Spaces.GetAssignedIsolationSegment(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch assigned isolation segment.",
			fmt.Sprintf("Request failed with %s.", err.Error()),
		)
		return
	}

	data.setTypeValueFromString(isolationSegment)

	runningSecurityGroups, err := d.cfClient.SecurityGroups.ListRunningForSpaceAll(ctx, data.Id.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch running security groups.",
			fmt.Sprintf("Request failed with %s.", err.Error()),
		)
		return
	}

	diags = data.setTypeValueFromSecurityGroups(ctx, runningSecurityGroups, "running")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	stagingSecurityGroups, err := d.cfClient.SecurityGroups.ListStagingForSpaceAll(ctx, data.Id.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch staging security groups.",
			fmt.Sprintf("Request failed with %s.", err.Error()),
		)
		return
	}

	diags = data.setTypeValueFromSecurityGroups(ctx, stagingSecurityGroups, "staging")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "read a space data source")
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}
