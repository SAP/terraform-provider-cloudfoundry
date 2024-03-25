package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-cloudfoundry/internal/provider/managers"
	"github.com/SAP/terraform-provider-cloudfoundry/internal/validation"
	cfv3client "github.com/cloudfoundry-community/go-cfclient/v3/client"
	cfv3resource "github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/samber/lo"
)

var (
	_ resource.Resource              = &orgQuotaResource{}
	_ resource.ResourceWithConfigure = &orgQuotaResource{}
)

func NewOrgQuotaResource() resource.Resource {
	return &orgQuotaResource{}
}

type orgQuotaResource struct {
	cfClient *cfv3client.Client
}

func (r *orgQuotaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_org_quota"
}

func (r *orgQuotaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a Cloud Foundry resource to manage org quota definitions.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name you use to identify the quota or plan in Cloud Foundry",
				Required:            true,
			},
			"allow_paid_service_plans": schema.BoolAttribute{
				MarkdownDescription: "Determines whether users can provision instances of non-free service plans. Does not control plan visibility. When false, non-free service plans may be visible in the marketplace but instances can not be provisioned.",
				Required:            true,
			},
			"total_services": schema.Int64Attribute{
				MarkdownDescription: "Maximum services allowed.",
				Optional:            true,
			},
			"total_service_keys": schema.Int64Attribute{
				MarkdownDescription: "Maximum service keys allowed.",
				Optional:            true,
			},
			"total_routes": schema.Int64Attribute{
				MarkdownDescription: "Maximum routes allowed.",
				Optional:            true,
			},
			"total_route_ports": schema.Int64Attribute{
				MarkdownDescription: "Maximum routes with reserved ports.",
				Optional:            true,
			},
			"total_private_domains": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of private domains allowed to be created within the Org.",
				Optional:            true,
			},
			"total_memory": schema.Int64Attribute{
				MarkdownDescription: "Maximum memory usage allowed.",
				Optional:            true,
			},
			"instance_memory": schema.Int64Attribute{
				MarkdownDescription: "Maximum memory per application instance.",
				Optional:            true,
			},
			"total_app_instances": schema.Int64Attribute{
				MarkdownDescription: "Maximum app instances allowed.",
				Optional:            true,
			},
			"total_app_tasks": schema.Int64Attribute{
				MarkdownDescription: "Maximum tasks allowed per app.",
				Optional:            true,
			},
			"total_app_log_rate_limit": schema.Int64Attribute{
				MarkdownDescription: "Maximum log rate allowed for all the started processes and running tasks in bytes/second.",
				Optional:            true,
			},
			"orgs": schema.SetAttribute{
				MarkdownDescription: "Set of Org GUIDs to which this org quota would be assigned.",
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(validation.ValidUUID()),
					setvalidator.SizeAtLeast(1),
				},
			},
			idKey:        guidSchema(),
			createdAtKey: createdAtSchema(),
			updatedAtKey: updatedAtSchema(),
		},
	}
}
func (r *orgQuotaResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *orgQuotaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var orgQuotaType OrgQuotaType
	diags := req.Plan.Get(ctx, &orgQuotaType)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgsQuotaValues, diags := orgQuotaType.mapOrgQuotaTypeToValues(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	orgsQuotaResp, err := r.cfClient.OrganizationQuotas.Create(ctx, orgsQuotaValues)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create org quota",
			fmt.Sprintf("Request failed with %s ", err.Error()),
		)
		return
	}
	orgsQuotaType, diags := mapOrgQuotaValuesToType(orgsQuotaResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, orgsQuotaType)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *orgQuotaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var orgQuotaTypeState OrgQuotaType
	diags := req.State.Get(ctx, &orgQuotaTypeState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	orgsQuotas, err := r.cfClient.OrganizationQuotas.ListAll(ctx, &cfv3client.OrganizationQuotaListOptions{
		GUIDs: cfv3client.Filter{
			Values: []string{orgQuotaTypeState.ID.ValueString()},
		},
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch org quota data",
			fmt.Sprintf("Request failed with %s", err.Error()),
		)
		return
	}
	orgsQuota, found := lo.Find(orgsQuotas, func(orgQuota *cfv3resource.OrganizationQuota) bool {
		return orgQuota.GUID == orgQuotaTypeState.ID.ValueString()
	})
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	orgsQuotaType, diags := mapOrgQuotaValuesToType(orgsQuota)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, orgsQuotaType)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *orgQuotaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var orgQuotaTypePlan OrgQuotaType
	var orgQuotaTypeState OrgQuotaType
	diags := req.Plan.Get(ctx, &orgQuotaTypePlan)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Get(ctx, &orgQuotaTypeState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	removed, added, diags := findChangedRelationsFromTFState(ctx, orgQuotaTypePlan.Organizations, orgQuotaTypeState.Organizations)
	resp.Diagnostics.Append(diags...)
	orgsQuotaValues, diags := orgQuotaTypePlan.mapOrgQuotaTypeToValues(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if len(removed) != 0 {
		resp.Diagnostics.AddError(
			"Unable to update org quota",
			fmt.Sprintf("Cannot unassign org quota from org %v", removed),
		)
		return
	}
	if len(added) != 0 {
		_, err := r.cfClient.OrganizationQuotas.Apply(ctx, orgQuotaTypePlan.ID.ValueString(), added)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to update org quota",
				fmt.Sprintf("Request failed with %s", err.Error()),
			)
			return
		}
	}
	orgsQuotaValues.Relationships = nil
	orgsQuotaResp, err := r.cfClient.OrganizationQuotas.Update(ctx, orgQuotaTypePlan.ID.ValueString(), orgsQuotaValues)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update org quota",
			fmt.Sprintf("Request failed with %s", err.Error()),
		)
	}
	orgsQuotaType, diags := mapOrgQuotaValuesToType(orgsQuotaResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, orgsQuotaType)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *orgQuotaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var orgQuotaType OrgQuotaType
	diags := req.State.Get(ctx, &orgQuotaType)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	jobID, err := r.cfClient.OrganizationQuotas.Delete(ctx, orgQuotaType.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete organization quota",
			fmt.Sprintf("Org quota deletion verification failed %s with %s", orgQuotaType.Name.ValueString(), err.Error()),
		)
		return
	}
	if err = pollJob(ctx, *r.cfClient, jobID); err != nil {
		resp.Diagnostics.AddError(
			"Unable to verify org quota deletion",
			"Org quota deletion verification failed for "+orgQuotaType.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}

func (r *orgQuotaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
