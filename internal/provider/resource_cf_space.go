package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/SAP/terraform-provider-cloudfoundry/internal/provider/managers"
	"github.com/SAP/terraform-provider-cloudfoundry/internal/validation"
	cfv3client "github.com/cloudfoundry-community/go-cfclient/v3/client"
	cfv3resource "github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &SpaceResource{}
	_ resource.ResourceWithConfigure   = &SpaceResource{}
	_ resource.ResourceWithImportState = &SpaceResource{}
)

func NewSpaceResource() resource.Resource {
	return &SpaceResource{}
}

type SpaceResource struct {
	cfClient *cfv3client.Client
}

func (r *SpaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space"
}

func (r *SpaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a Cloud Foundry resource for managing Cloud Foundry spaces within organizations.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Space in Cloud Foundry",
				Required:            true,
			},
			"org": schema.StringAttribute{
				MarkdownDescription: "The ID of the Org within which to create the space",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.String{
					validation.ValidUUID(),
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("org_name"),
					}...),
				},
			},
			"org_name": schema.StringAttribute{
				MarkdownDescription: "The ID of the Org within which to create the space",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("org"),
					}...),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The GUID of the space",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"quota": schema.StringAttribute{
				MarkdownDescription: "The ID of the Space quota or plan defined for the owning Org. Specifying an empty string request unassigns any space quota from the space. Defaults to empty string.",
				Optional:            true,
				Validators: []validator.String{
					validation.ValidUUID(),
				},
			},
			"allow_ssh": schema.BoolAttribute{
				MarkdownDescription: "Allows SSH to application containers via the CF CLI.",
				Optional:            true,
			},
			"isolation_segment": schema.StringAttribute{
				MarkdownDescription: "The ID of the isolation segment to assign to the space. The segment must be entitled to the space's parent organization. If the isolation segment id is unspecified, then Cloud Foundry assigns the space to the org’s default isolation segment if any. Note that existing apps in the space will not run in a newly assigned isolation segment until they are restarted.",
				Optional:            true,
				Validators: []validator.String{
					validation.ValidUUID(),
				},
			},
			"asgs": schema.SetAttribute{
				MarkdownDescription: "List of running application security groups to apply to applications running within this space. Defaults to empty list.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"staging_asgs": schema.SetAttribute{
				MarkdownDescription: "List of staging application security groups to apply to applications being staged for this space. Defaults to empty list.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"labels": schema.MapAttribute{
				MarkdownDescription: `The labels associated with Cloud Foundry resources. Add as described [here](https://docs.cloudfoundry.org/adminguide/metadata.html#-view-metadata-for-an-object).`,
				ElementType:         types.StringType,
				Optional:            true,
			},
			"annotations": schema.MapAttribute{
				MarkdownDescription: "The annotations associated with Cloud Foundry resources. Add as described [here](https://docs.cloudfoundry.org/adminguide/metadata.html#-view-metadata-for-an-object).",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was updated in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
		},
	}
}

func (r *SpaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.cfClient = session.CFClient
}

// Create creates the resource and sets the initial Terraform state.
func (r *SpaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan spaceType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = plan.populateOrgValues(ctx, r.cfClient)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createSpace, diags := plan.setCreateSpaceValuesFromPlan(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	space, err := r.cfClient.Spaces.Create(ctx, &createSpace)
	if err != nil {
		resp.Diagnostics.AddError(
			"API Error Creating Space",
			"Could not create Space "+plan.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	if !plan.Quota.IsNull() {
		_, err := r.cfClient.SpaceQuotas.Get(ctx, plan.Quota.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"API Error Getting Space Quota",
				"Could not get the Quota with ID "+plan.Quota.ValueString()+": "+err.Error(),
			)
			return
		}
		_, err = r.cfClient.SpaceQuotas.Apply(ctx, plan.Quota.ValueString(), []string{
			space.GUID,
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"API Error Applying Space Quota",
				"Could not apply the Quota with ID "+plan.Quota.ValueString()+" to the space "+plan.Name.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	if !plan.AllowSSH.IsNull() {
		allowSSH, err := r.cfClient.SpaceFeatures.IsSSHEnabled(ctx, space.GUID)
		if err != nil {
			resp.Diagnostics.AddError(
				"API Error Fetching Space SSH",
				"Could not get the SSH feature value on space with ID "+space.GUID+" to the space "+space.Name+": "+err.Error(),
			)
			return
		}
		if allowSSH != plan.AllowSSH.ValueBool() {
			err = r.cfClient.SpaceFeatures.EnableSSH(ctx, space.GUID, plan.AllowSSH.ValueBool())
			if err != nil {
				resp.Diagnostics.AddError(
					"API Error Setting Space SSH",
					"Could not set the SSH feature value on space with ID "+space.GUID+" to the space "+space.Name+": "+err.Error(),
				)
				return
			}
		}
	}

	if !plan.IsolationSegment.IsNull() {
		_, err := r.cfClient.IsolationSegments.Get(ctx, plan.IsolationSegment.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"API Error Getting Isolation Segment",
				"Could not get the Isolation Segment with ID "+plan.IsolationSegment.ValueString()+": "+err.Error(),
			)
			return
		}
		err = r.cfClient.Spaces.AssignIsolationSegment(ctx, space.GUID, plan.IsolationSegment.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"API Error Assigning Isolation Segment",
				"Could not assign the Isolation Segment with ID "+plan.IsolationSegment.ValueString()+" on space "+space.Name+": "+err.Error(),
			)
			return
		}
	}

	if !plan.RunningSecurityGroups.IsNull() {
		var runningGroupsInput []string
		runningSecurityGroupsDiagnostics := plan.RunningSecurityGroups.ElementsAs(ctx, &runningGroupsInput, false)
		resp.Diagnostics.Append(runningSecurityGroupsDiagnostics...)
		for _, securityGroup := range runningGroupsInput {
			_, err := r.cfClient.SecurityGroups.Get(ctx, securityGroup)
			if err != nil {
				resp.Diagnostics.AddError(
					"API Error Getting Security Group",
					"Could not get the Security Group with ID "+securityGroup+": "+err.Error(),
				)
				return
			}
			_, err = r.cfClient.SecurityGroups.BindRunningSecurityGroup(ctx, securityGroup, []string{space.GUID})
			if err != nil {
				resp.Diagnostics.AddError(
					"API Error Assigning Running Security Group",
					"Could not assign the Security Group with ID "+securityGroup+" to the space "+space.Name+": "+err.Error(),
				)
				return
			}
		}
	}

	if !plan.StagingSecurityGroups.IsNull() {
		var stagingGroupsInput []string
		stagingSecurityGroupsDiagnostics := plan.StagingSecurityGroups.ElementsAs(ctx, &stagingGroupsInput, false)
		resp.Diagnostics.Append(stagingSecurityGroupsDiagnostics...)
		for _, securityGroup := range stagingGroupsInput {
			_, err := r.cfClient.SecurityGroups.Get(ctx, securityGroup)
			if err != nil {
				resp.Diagnostics.AddError(
					"API Error Getting Security Group",
					"Could not get the Security Group with ID "+securityGroup+": "+err.Error(),
				)
				return
			}
			_, err = r.cfClient.SecurityGroups.BindStagingSecurityGroup(ctx, securityGroup, []string{space.GUID})
			if err != nil {
				resp.Diagnostics.AddError(
					"API Error Assigning Staging Security Group",
					"Could not assign the Security Group with ID "+securityGroup+" to the space "+space.Name+": "+err.Error(),
				)
				return
			}
		}
	}

	plan.setComputedTypeValuesFromSpace(ctx, space)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rs *SpaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data spaceType

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	space, err := rs.cfClient.Spaces.Get(ctx, data.Id.ValueString())
	if err != nil {
		cfError, isCfError := err.(cfv3resource.IsResourceNotFound)
		if isCfError {
			if cfError.Code == 10010 {
				resp.State.RemoveResource(ctx)
				return
			}
		}
		cfv3resource.IsResourceNotFoundError()
		resp.Diagnostics.AddError(
			"Unable to fetch space data.",
			fmt.Sprintf("Request failed with %s.", err.Error()),
		)
		return
	}

	diags = data.setTypeValuesFromSpace(ctx, space)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = data.populateOrgValues(ctx, rs.cfClient)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sshEnabled, err := rs.cfClient.SpaceFeatures.IsSSHEnabled(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch space features.",
			fmt.Sprintf("Request failed with %s.", err.Error()),
		)
		return
	}

	data.setTypeValueFromBool(sshEnabled)

	isolationSegment, err := rs.cfClient.Spaces.GetAssignedIsolationSegment(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch assigned isolation segment.",
			fmt.Sprintf("Request failed with %s.", err.Error()),
		)
		return
	}

	data.setTypeValueFromString(isolationSegment)

	runningSecurityGroups, err := rs.cfClient.SecurityGroups.ListRunningForSpaceAll(ctx, data.Id.ValueString(), nil)
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

	stagingSecurityGroups, err := rs.cfClient.SecurityGroups.ListStagingForSpaceAll(ctx, data.Id.ValueString(), nil)
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

	tflog.Trace(ctx, "read a space resource")
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (rs *SpaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, previousState spaceType
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &previousState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !previousState.Name.Equal(plan.Name) || !previousState.Labels.Equal(plan.Labels) || !previousState.Annotations.Equal(plan.Annotations) {

		updateSpace, diags := plan.setUpdateSpaceValuesFromPlan(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		_, err := rs.cfClient.Spaces.Update(ctx, plan.Id.ValueString(), &updateSpace)
		if err != nil {
			resp.Diagnostics.AddError(
				"API Error Updating Space",
				"Could not update Space with ID "+plan.Id.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	if !previousState.Quota.Equal(plan.Quota) {
		if plan.Quota.IsNull() {
			_, err := rs.cfClient.SpaceQuotas.(ctx, plan.Quota.ValueString(), []string{
				plan.Id.ValueString(),
			})
		}
		_, err := rs.cfClient.SpaceQuotas.Apply(ctx, plan.Quota.ValueString(), []string{
			plan.Id.ValueString(),
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"API Error Updating Org Quota",
				"Could not apply the Quota with ID "+plan.Quota.ValueString()+" to the space with ID"+plan.Id.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	if !previousState.AllowSSH.Equal(plan.AllowSSH) {
		err := rs.cfClient.SpaceFeatures.EnableSSH(ctx, plan.Id.ValueString(), plan.AllowSSH.ValueBool())
		if err != nil {
			resp.Diagnostics.AddError(
				"API Error Updating Space SSH",
				"Could not set the SSH feature value on space with ID "+plan.Id.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	if !previousState.IsolationSegment.Equal(plan.IsolationSegment) {
		err := rs.cfClient.Spaces.AssignIsolationSegment(ctx, plan.Id.ValueString(), plan.IsolationSegment.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"API Error Updating Isolation Segment",
				"Could not assign the Isolation Segment with ID "+plan.IsolationSegment.ValueString()+" on space with ID"+plan.Id.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	if !previousState.RunningSecurityGroups.Equal(plan.RunningSecurityGroups) {
		err := rs.cfClient.Spaces.AssignIsolationSegment(ctx, plan.Id.ValueString(), plan.IsolationSegment.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"API Error Updating Isolation Segment",
				"Could not assign the Isolation Segment with ID "+plan.IsolationSegment.ValueString()+" on space with ID"+plan.Id.ValueString()+": "+err.Error(),
			)
			return
		}
	}

}

func (rs *SpaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state spaceType

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	jobID, err := rs.cfClient.Spaces.Delete(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"API Error Deleting Space",
			"Could not delete the space with ID "+state.Id.ValueString()+" and name "+state.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	err = rs.cfClient.Jobs.PollComplete(ctx, jobID, &cfv3client.PollingOptions{
		FailedState:   "FAILED",
		Timeout:       time.Second * 30,
		CheckInterval: time.Second,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"API Error Deleting Space",
			"Could not delete the space with ID "+state.Id.ValueString()+" and name "+state.Name.ValueString()+": "+err.Error(),
		)
		return
	}

}

func (rs *SpaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
