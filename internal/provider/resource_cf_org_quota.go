package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/SAP/terraform-provider-cloudfoundry/internal/provider/managers"
	cfv3client "github.com/cloudfoundry-community/go-cfclient/v3/client"
	cfv3resource "github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
		MarkdownDescription: "Provides a Cloud Foundry resource to manage org quotas definitions.",
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
				MarkdownDescription: "Maximum services allowed",
				Optional:            true,
			},
			"total_service_keys": schema.Int64Attribute{
				MarkdownDescription: "Maximum service keys allowed",
				Optional:            true,
			},
			"total_routes": schema.Int64Attribute{
				MarkdownDescription: "Maximum routes allowed",
				Optional:            true,
			},
			"total_route_ports": schema.Int64Attribute{
				MarkdownDescription: "Maximum routes with reserved ports",
				Optional:            true,
			},
			"total_private_domains": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of private domains allowed to be created within the Org",
				Optional:            true,
			},
			"total_memory": schema.Int64Attribute{
				MarkdownDescription: "Maximum memory usage allowed",
				Optional:            true,
			},
			"instance_memory": schema.Int64Attribute{
				MarkdownDescription: "Maximum memory per application instance",
				Optional:            true,
			},
			"total_app_instances": schema.Int64Attribute{
				MarkdownDescription: "Maximum app instances allowed",
				Optional:            true,
			},
			"total_app_tasks": schema.Int64Attribute{
				MarkdownDescription: "Maximum tasks allowed per app",
				Optional:            true,
			},
			"total_app_log_rate_limit": schema.Int64Attribute{
				MarkdownDescription: "Maximum log rate allowed for all the started processes and running tasks in bytes/second.",
				Optional:            true,
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
			fmt.Sprintf("Expected *managers.Session, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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

	orgsQuotaValues := orgQuotaType.mapOrgQuotaTypeToValues()
	orgsQuotaResp, err := r.cfClient.OrganizationQuotas.Create(ctx, orgsQuotaValues)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create org quota.",
			fmt.Sprintf("Request failed with %s.", err.Error()),
		)
		return
	}
	orgsQuotaType := mapOrgQuotaValuesToType(orgsQuotaResp)
	diags = resp.State.Set(ctx, orgsQuotaType)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *orgQuotaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var orgQuotaType OrgQuotaType
	diags := req.State.Get(ctx, &orgQuotaType)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	orgsQuotas, err := r.cfClient.OrganizationQuotas.ListAll(ctx, &cfv3client.OrganizationQuotaListOptions{
		Names: cfv3client.Filter{
			Values: []string{orgQuotaType.Name.ValueString()},
		},
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch org quota data.",
			fmt.Sprintf("Request failed with %s.", err.Error()),
		)
		return
	}
	orgsQuota, found := lo.Find(orgsQuotas, func(orgQuota *cfv3resource.OrganizationQuota) bool {
		return orgQuota.Name == orgQuotaType.Name.ValueString()
	})
	if !found {
		resp.Diagnostics.AddError(
			"Unable to find org quota.",
			"Org quota does not exist.",
		)
		return
	}
	orgsQuotaType := mapOrgQuotaValuesToType(orgsQuota)
	diags = resp.State.Set(ctx, orgsQuotaType)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *orgQuotaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var orgQuotaType OrgQuotaType
	diags := req.Plan.Get(ctx, &orgQuotaType)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgsQuotaValues := orgQuotaType.mapOrgQuotaTypeToValues()
	orgsQuotaResp, err := r.cfClient.OrganizationQuotas.Update(ctx, orgQuotaType.Id.ValueString(), orgsQuotaValues)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update org quota.",
			fmt.Sprintf("Request failed with %s.", err.Error()),
		)
		return
	}
	orgsQuotaType := mapOrgQuotaValuesToType(orgsQuotaResp)
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
	err := r.cfClient.OrganizationQuotas.Delete(ctx, orgQuotaType.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list spaces for Org Deletion",
			"Could not list spaces to delete them before deleting organization",
		)
		return
	}
	_, _, err = lo.AttemptWithDelay(10, 3*time.Second, func(i int, duration time.Duration) error {
		_, err := r.cfClient.OrganizationQuotas.Get(ctx, orgQuotaType.Id.ValueString())
		if err != nil {
			if cfv3resource.IsResourceNotFoundError(err) {
				return nil
			} else {
				return err
			}
		} else {
			return fmt.Errorf("resource still exists")
		}
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to verify Org Quota Deletion",
			"Org quota deletion verification failed with %s."+err.Error(),
		)
		return
	}
}
