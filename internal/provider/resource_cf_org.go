package provider

import (
	"context"
	"fmt"
	"maps"

	"github.com/SAP/terraform-provider-cloudfoundry/internal/provider/managers"
	"github.com/SAP/terraform-provider-cloudfoundry/internal/validation"
	cfv3client "github.com/cloudfoundry-community/go-cfclient/v3/client"
	cfv3resource "github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/samber/lo"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &orgResource{}
)

// NeworgResource is a helper function to simplify the provider implementation.
func NeworgResource() resource.Resource {
	return &orgResource{}
}

// orgResource is the resource implementation.
type orgResource struct {
	cfClient *cfv3client.Client
}

// Metadata returns the resource type name.
func (r *orgResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_org"
}

// Configure defines the configuration for the resource
func (r *orgResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

// Schema defines the schema for the resource.
func (r *orgResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Creates a Cloud Foundry Organization 
		
		__Further documentation:__
		https://docs.cloudfoundry.org/concepts/roles.html#orgs
		`,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Organization in Cloud Foundry",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the Organization",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
			"quota": schema.StringAttribute{
				MarkdownDescription: "The ID of quota to be applied to this Org. By default, no quota is assigned to the org.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					validation.ValidUUID(),
				},
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was updated in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"suspended": schema.BoolAttribute{
				MarkdownDescription: "Whether an organization is suspended or not",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}

}

// Create creates the resource and sets the initial Terraform state.
func (r *orgResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan orgType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var createOrg cfv3resource.OrganizationCreate

	if !plan.Suspended.IsUnknown() {
		suspend := plan.Suspended.ValueBool()
		createOrg.Suspended = &suspend
		//TODO @k7 check if below works
		//createOrg.Suspended = plan.Suspended.ValueBoolPointer()
	}
	labels := make(map[string]*string)
	labelsDiags := plan.Labels.ElementsAs(ctx, &labels, false)
	resp.Diagnostics.Append(labelsDiags...)
	maps.Copy(createOrg.Metadata.Labels, labels)

	annotations := make(map[string]*string)
	annotationsDiags := plan.Annotations.ElementsAs(ctx, &annotations, false)
	resp.Diagnostics.Append(annotationsDiags...)
	maps.Copy(createOrg.Metadata.Annotations, annotations)

	org, err := r.cfClient.Organizations.Create(ctx, &cfv3resource.OrganizationCreate{})

	if err != nil {
		resp.Diagnostics.AddError(
			"API Error Creating Organization",
			"Could not create Organization "+plan.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	// Assign Quota to the organization
	if !plan.Quota.IsUnknown() {
		quota, err := r.cfClient.OrganizationQuotas.Get(ctx, plan.Quota.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"API Error Getting Org Quota",
				"Could not get the Quota with ID "+plan.Quota.ValueString()+"and name "+quota.Name+": "+err.Error(),
			)
			return
		}
		_, err = r.cfClient.OrganizationQuotas.Apply(ctx, plan.Quota.ValueString(), []string{
			org.GUID,
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"API Error Applying Org Quota",
				"Could not apply the Quota with ID "+plan.Quota.ValueString()+"and name "+quota.Name+" to the org "+plan.Name.ValueString()+": "+err.Error(),
			)
			return
		}

	}

	plan, diags = mapOrgValuesToType(ctx, org)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *orgResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data orgType

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	orgs, err := r.cfClient.Organizations.ListAll(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch org data.",
			fmt.Sprintf("Request failed with %s.", err.Error()),
		)
		return
	}
	org, found := lo.Find(orgs, func(org *cfv3resource.Organization) bool {
		return org.Name == data.Name.ValueString()
	})
	if !found {
		resp.Diagnostics.AddError(
			"Unable to find org data in list",
			fmt.Sprintf("Given name %s not in the list of orgs.", data.Name.ValueString()),
		)
		return
	}
	data, diags = mapOrgValuesToType(ctx, org)
	resp.Diagnostics.Append(diags...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *orgResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan orgType

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateOrg := cfv3resource.OrganizationUpdate{
		Name:      "",
		Suspended: new(bool),
		Metadata:  &cfv3resource.Metadata{},
	}
	org, err := r.cfClient.Organizations.Update(ctx, plan.ID.ValueString(), &updateOrg)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update Org",
			"Could not update org "+plan.Name.ValueString()+": "+err.Error(),
		)
	}
	plan, diags = mapOrgValuesToType(ctx, org)
	resp.Diagnostics.Append(diags...)

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *orgResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state orgType

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaces, err := r.cfClient.Spaces.ListAll(ctx, &cfv3client.SpaceListOptions{
		OrganizationGUIDs: cfv3client.Filter{
			Values: []string{
				state.ID.ValueString(),
			},
		},
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list spaces for Org Deletion",
			"Could not list spaces to delete them before deleting organization"+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
	for _, space := range spaces {
		_, err := r.cfClient.Spaces.Delete(ctx, space.GUID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to list spaces for Org Deletion",
				"Could not list spaces to delete them before deleting organization"+state.ID.ValueString()+": "+err.Error(),
			)
			return
		}
	}

}
