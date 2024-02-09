package provider

import (
	"time"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Terraform struct for storing values for user data source and resource
type roleType struct {
	Type         types.String `tfsdk:"type"`
	User         types.String `tfsdk:"user"`
	Space        types.String `tfsdk:"space"`
	Id           types.String `tfsdk:"id"`
	Organization types.String `tfsdk:"org"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
}

// Returns the OrganizationRoleType value needed for org role creation
func (data *roleType) getOrgRoleType() resource.OrganizationRoleType {

	switch data.Type.ValueString() {
	case "organization_user":
		return resource.OrganizationRoleUser
	case "organization_auditor":
		return resource.OrganizationRoleAuditor
	case "organization_manager":
		return resource.OrganizationRoleManager
	case "organization_billing_manager":
		return resource.OrganizationRoleBillingManager
	default:
		return resource.OrganizationRoleNone

	}
}

// Returns the OrganizationRoleType value needed for org role creation
func (data *roleType) getSpaceRoleType() resource.SpaceRoleType {

	switch data.Type.ValueString() {
	case "space_auditor":
		return resource.SpaceRoleAuditor
	case "space_developer":
		return resource.SpaceRoleDeveloper
	case "space_manager":
		return resource.SpaceRoleManager
	case "space_supporter":
		return resource.SpaceRoleSupporter
	default:
		return resource.SpaceRoleNone

	}
}

// Sets the terraform struct values from the user resource returned by the cf-client
func mapRoleValuesToType(role *resource.Role) roleType {

	roleType := roleType{
		Id:        types.StringValue(role.GUID),
		CreatedAt: types.StringValue(role.CreatedAt.Format(time.RFC3339)),
		UpdatedAt: types.StringValue(role.UpdatedAt.Format(time.RFC3339)),
		Type:      types.StringValue(role.Type),
		User:      types.StringValue(role.Relationships.User.Data.GUID),
	}

	if role.Relationships.Org.Data != nil {
		roleType.Organization = types.StringValue(role.Relationships.Org.Data.GUID)
	}

	if role.Relationships.Space.Data != nil {
		roleType.Space = types.StringValue(role.Relationships.Space.Data.GUID)
	}

	return roleType
}
