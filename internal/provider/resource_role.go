package provider

import (
	"context"
	"fmt"

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
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &RoleResource{}
	_ resource.ResourceWithConfigure   = &RoleResource{}
	_ resource.ResourceWithImportState = &RoleResource{}
)

// Instantiates a space resource
func NewRoleResource() resource.Resource {
	return &RoleResource{}
}

// Contains reference to the v3 client to be used for making the API calls
type RoleResource struct {
	cfClient *cfv3client.Client
}

func (r *RoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (r *RoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a Cloud Foundry resource for assigning roles. Space roles cannot be assigned until the user has the relevant role in the organization. (Updating a role is not supported according to the docs)",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				MarkdownDescription: "Role type; see [Valid role types](https://v3-apidocs.cloudfoundry.org/version/3.154.0/index.html#valid-role-types)",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("space_auditor", "space_developer", "space_manager", "space_supporter",
						"organization_user", "organization_auditor", "organization_manager", "organization_billing_manager",
					),
				},
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "The guid of the cloudfoundry user to assign the role with",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"org": schema.StringAttribute{
				MarkdownDescription: "The guid of the organization to assign the role to",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validation.ValidUUID(),
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("space"),
						path.MatchRoot("org"),
					}...),
				},
			},
			"space": schema.StringAttribute{
				MarkdownDescription: "The guid of the space to assign the role to",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validation.ValidUUID(),
				},
			},
			idKey:        guidSchema(),
			createdAtKey: createdAtSchema(),
			updatedAtKey: updatedAtSchema(),
		},
	}
}

func (r *RoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan roleType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		role *cfv3resource.Role
		err  error
	)
	if !plan.Organization.IsNull() {
		orgRoleType := plan.getOrgRoleType()
		if orgRoleType == cfv3resource.OrganizationRoleNone {
			resp.Diagnostics.AddError(
				"Invalid Role Type",
				"Could not register Space Role "+plan.Type.ValueString()+" for the organization with ID "+plan.Organization.ValueString()+". Please assign an organization role instead.",
			)
			return
		}
		role, err = r.cfClient.Roles.CreateOrganizationRole(ctx, plan.Organization.ValueString(), plan.User.ValueString(), orgRoleType)
	} else {
		spaceRoleType := plan.getSpaceRoleType()
		if spaceRoleType == cfv3resource.SpaceRoleNone {
			resp.Diagnostics.AddError(
				"Invalid Role Type",
				"Could not register Organization Role "+plan.Type.ValueString()+" for the space with ID "+plan.Space.ValueString()+". Please assign a space role instead.",
			)
			return
		}
		role, err = r.cfClient.Roles.CreateSpaceRole(ctx, plan.Space.ValueString(), plan.User.ValueString(), spaceRoleType)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"API Error Registering Role",
			"Could not register Role with user ID "+plan.Id.ValueString()+" : "+err.Error(),
		)
		return
	}

	plan = mapRoleValuesToType(role)
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, "created a role resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (rs *RoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data roleType

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	role, err := rs.cfClient.Roles.Get(ctx, data.Id.ValueString())
	if err != nil {
		handleReadErrors(ctx, resp, err, "role", data.Id.ValueString())
		return
	}

	data = mapRoleValuesToType(role)
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, "read a role resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update for role is not possible
func (rs *RoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (rs *RoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state roleType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	jobID, err := rs.cfClient.Roles.Delete(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"API Error Deleting Role",
			"Could not delete the role with ID "+state.Id.ValueString()+" and user ID "+state.User.ValueString()+" : "+err.Error(),
		)
		return
	}

	if pollJob(ctx, *rs.cfClient, jobID) != nil {
		resp.Diagnostics.AddError(
			"API Error Deleting Role",
			"Failed in deleting the role with ID "+state.Id.ValueString()+" and user ID "+state.User.ValueString()+" : "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "deleted a role resource")
}

func (rs *RoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
