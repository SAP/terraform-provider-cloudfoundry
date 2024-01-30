package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/samber/lo"
)

type spaceType struct {
	Name                  types.String `tfsdk:"name"`
	Id                    types.String `tfsdk:"id"`
	OrgId                 types.String `tfsdk:"org"`
	OrgName               types.String `tfsdk:"org_name"`
	Quota                 types.String `tfsdk:"quota"`
	AllowSSH              types.Bool   `tfsdk:"allow_ssh"`
	IsolationSegment      types.String `tfsdk:"isolation_segment"`
	RunningSecurityGroups types.Set    `tfsdk:"asgs"`
	StagingSecurityGroups types.Set    `tfsdk:"staging_asgs"`
	Labels                types.Map    `tfsdk:"labels"`
	Annotations           types.Map    `tfsdk:"annotations"`
	CreatedAt             types.String `tfsdk:"created_at"`
	UpdatedAt             types.String `tfsdk:"updated_at"`
}

// Sets the terraform struct values from the space resource returned by the cf-client
func (data *spaceType) setTypeValuesFromSpace(ctx context.Context, space *resource.Space) diag.Diagnostics {

	data.Name = types.StringValue(space.Name)
	data.Id = types.StringValue(space.GUID)
	data.OrgId = types.StringValue(space.Relationships.Organization.Data.GUID)

	if space.Relationships.Quota.Data != nil {
		data.Quota = types.StringValue(space.Relationships.Quota.Data.GUID)
	} else {
		data.Quota = types.StringValue("")
	}

	data.CreatedAt = types.StringValue(space.CreatedAt.Format(time.RFC3339))
	data.UpdatedAt = types.StringValue(space.UpdatedAt.Format(time.RFC3339))

	var diags, diagnostics diag.Diagnostics
	data.Labels, diags = types.MapValueFrom(ctx, types.StringType, space.Metadata.Labels)
	diagnostics.Append(diags...)
	data.Annotations, diags = types.MapValueFrom(ctx, types.StringType, space.Metadata.Annotations)
	diagnostics.Append(diags...)

	return diagnostics
}

// Sets the terraform struct allow_ssh value from a bool
func (data *spaceType) setTypeValueFromBool(sshEnabled bool) {
	data.AllowSSH = types.BoolValue(sshEnabled)
}

// Sets the terraform struct isolation_segment value from a string
func (data *spaceType) setTypeValueFromString(isolationSegment string) {
	data.IsolationSegment = types.StringValue(isolationSegment)
}

// Sets the terraform struct asgs or staging_asgs value from the security group resource returned by the cf-client depending on the "running" or "staging" passed as SecurityGroupType
func (data *spaceType) setTypeValueFromSecurityGroups(ctx context.Context, groups []*resource.SecurityGroup, SecurityGroupType string) diag.Diagnostics {

	//Not sure of security group logics
	var spaceSecurityGroups []string
	var diags, diagnostics diag.Diagnostics

	if SecurityGroupType == "running" {
		spaceSecurityGroups = lo.FilterMap(groups, func(group *resource.SecurityGroup, _ int) (string, bool) {
			if !group.GloballyEnabled.Running {
				return group.GUID, true
			}
			return "", false
		})
		data.RunningSecurityGroups, diags = types.SetValueFrom(ctx, types.StringType, spaceSecurityGroups)

	} else if SecurityGroupType == "staging" {
		spaceSecurityGroups = lo.FilterMap(groups, func(group *resource.SecurityGroup, _ int) (string, bool) {
			if !group.GloballyEnabled.Staging {
				return group.GUID, true
			}
			return "", false
		})
		data.StagingSecurityGroups, diags = types.SetValueFrom(ctx, types.StringType, spaceSecurityGroups)
	}

	diagnostics.Append(diags...)
	return diagnostics
}

// Sets the space resource values for creation with cf-client from the terraform struct values
func (data *spaceType) setCreateSpaceValuesFromPlan(ctx context.Context) (resource.SpaceCreate, diag.Diagnostics) {

	createSpace := resource.NewSpaceCreate(data.Name.ValueString(), data.OrgId.ValueString())

	var diagnostics diag.Diagnostics
	createSpace.Metadata = resource.NewMetadata()

	labelsDiags := data.Labels.ElementsAs(ctx, &createSpace.Metadata.Labels, false)
	diagnostics.Append(labelsDiags...)

	annotationsDiags := data.Annotations.ElementsAs(ctx, &createSpace.Metadata.Annotations, false)
	diagnostics.Append(annotationsDiags...)

	return *createSpace, diagnostics
}

// Sets the computed terraform struct values from the created space resource obtained from cf-client
func (data *spaceType) setComputedTypeValuesFromSpace(ctx context.Context, space *resource.Space) {

	data.Id = types.StringValue(space.GUID)
	data.CreatedAt = types.StringValue(space.CreatedAt.Format(time.RFC3339))
	data.UpdatedAt = types.StringValue(space.UpdatedAt.Format(time.RFC3339))

}

// Sets org name and org guid details in terraform struct if values missing
func (data *spaceType) populateOrgValues(ctx context.Context, cfClient *client.Client) diag.Diagnostics {

	var Diagnostics diag.Diagnostics

	if data.OrgId.IsUnknown() && data.OrgName.IsUnknown() || data.OrgId.IsNull() && data.OrgName.IsNull() {
		Diagnostics.AddError(
			"Neither Org GUID nor Org Name is present.",
			"Expected either 'org' or 'org_name' attribute.",
		)
		return Diagnostics

	} else if data.OrgId.IsUnknown() || data.OrgId.IsNull() {
		orgs, err := cfClient.Organizations.ListAll(ctx, &client.OrganizationListOptions{
			Names: client.Filter{
				Values: []string{
					data.OrgName.ValueString(),
				},
			},
		})
		if err != nil {
			Diagnostics.AddError(
				"Unable to fetch org data.",
				fmt.Sprintf("Request failed with %s.", err.Error()),
			)
			return Diagnostics
		}

		org, found := lo.Find(orgs, func(org *resource.Organization) bool {
			return org.Name == data.OrgName.ValueString()
		})

		if !found {
			Diagnostics.AddError(
				"Unable to find org data in list",
				fmt.Sprintf("Given name %s not in the list of orgs.", data.OrgName.ValueString()),
			)
			return Diagnostics
		}
		data.OrgId = types.StringValue(org.GUID)
	} else {
		//Fetching organization with GUID
		orgs, err := cfClient.Organizations.ListAll(ctx, &client.OrganizationListOptions{
			GUIDs: client.Filter{
				Values: []string{
					data.OrgId.ValueString(),
				},
			},
		})

		if err != nil {
			Diagnostics.AddError(
				"Unable to fetch org data.",
				fmt.Sprintf("Request failed with %s.", err.Error()),
			)
			return Diagnostics
		}

		org, found := lo.Find(orgs, func(org *resource.Organization) bool {
			return org.GUID == data.OrgId.ValueString()
		})

		if !found {
			Diagnostics.AddError(
				"Unable to find org data in list",
				fmt.Sprintf("Given org %s not in the list of orgs.", data.OrgId.ValueString()),
			)
			return Diagnostics
		}
		data.OrgName = types.StringValue(org.Name)
	}

	return Diagnostics
}

// Sets the space resource values for updating with cf-client from the terraform struct values
func (data *spaceType) setUpdateSpaceValuesFromPlan(ctx context.Context) (resource.SpaceUpdate, diag.Diagnostics) {

	updateSpace := &resource.SpaceUpdate{
		Name: data.Name.ValueString(),
	}
	var diagnostics diag.Diagnostics
	updateSpace.Metadata = resource.NewMetadata()

	labelsDiags := data.Labels.ElementsAs(ctx, &updateSpace.Metadata.Labels, false)
	diagnostics.Append(labelsDiags...)

	annotationsDiags := data.Annotations.ElementsAs(ctx, &updateSpace.Metadata.Annotations, false)
	diagnostics.Append(annotationsDiags...)

	return *updateSpace, diagnostics
}
