package provider

import (
	"context"
	"time"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Terraform struct for storing values for user data source and resource
type userType struct {
	UserName         types.String `tfsdk:"username"`
	PresentationName types.String `tfsdk:"presentation_name"`
	Origin           types.String `tfsdk:"origin"`
	Id               types.String `tfsdk:"id"`
	Labels           types.Map    `tfsdk:"labels"`
	Annotations      types.Map    `tfsdk:"annotations"`
	CreatedAt        types.String `tfsdk:"created_at"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
}

type datasourceUserType struct {
	Users []userType   `tfsdk:"users"`
	Name  types.String `tfsdk:"name"`
}

type spaceOrgUsersType struct {
	Users        []userType   `tfsdk:"users"`
	Space        types.String `tfsdk:"space"`
	Organization types.String `tfsdk:"org"`
}

// Sets the user resource values for creation with cf-client from the terraform struct values
func (data *userType) mapCreateUserTypeToValues(ctx context.Context) (resource.UserCreate, diag.Diagnostics) {

	createUser := &resource.UserCreate{GUID: data.Id.ValueString()}
	var diagnostics diag.Diagnostics
	createUser.Metadata = resource.NewMetadata()

	labelsDiags := data.Labels.ElementsAs(ctx, &createUser.Metadata.Labels, false)
	diagnostics.Append(labelsDiags...)

	annotationsDiags := data.Annotations.ElementsAs(ctx, &createUser.Metadata.Annotations, false)
	diagnostics.Append(annotationsDiags...)

	return *createUser, diagnostics
}

// Sets the terraform struct values from the user resource returned by the cf-client
func mapUserValuesToType(ctx context.Context, user *resource.User) (userType, diag.Diagnostics) {

	userType := userType{
		Id:               types.StringValue(user.GUID),
		CreatedAt:        types.StringValue(user.CreatedAt.Format(time.RFC3339)),
		UpdatedAt:        types.StringValue(user.UpdatedAt.Format(time.RFC3339)),
		PresentationName: types.StringValue(user.PresentationName),
		UserName:         types.StringValue(user.Username),
		Origin:           types.StringValue(user.Origin),
	}

	var diags, diagnostics diag.Diagnostics
	userType.Labels, diags = mapMetadataValueToType(ctx, user.Metadata.Labels)
	diagnostics.Append(diags...)
	userType.Annotations, diags = mapMetadataValueToType(ctx, user.Metadata.Annotations)
	diagnostics.Append(diags...)

	return userType, diagnostics
}

// Sets the user resource values for updation with cf-client from the terraform struct values
func (plan *userType) mapUpdateUserTypeToValues(ctx context.Context, state userType) (resource.UserUpdate, diag.Diagnostics) {

	updateUser := &resource.UserUpdate{}

	var diagnostics diag.Diagnostics
	updateUser.Metadata, diagnostics = setClientMetadataForUpdate(ctx, state.Labels, state.Annotations, plan.Labels, plan.Annotations)

	return *updateUser, diagnostics
}

// Prepares a terraform list from the user resources returned by the cf-client
func mapUsersValuesToListType(ctx context.Context, users []*resource.User) ([]userType, diag.Diagnostics) {

	var diagnostics diag.Diagnostics
	userValues := []userType{}
	for _, user := range users {
		userValue, diags := mapUserValuesToType(ctx, user)
		diagnostics.Append(diags...)
		userValues = append(userValues, userValue)
	}

	return userValues, diagnostics
}

// Sets the terraform struct values from the user resources returned by the cf-client
func mapUsersValuesToType(ctx context.Context, data datasourceUserType, users []*resource.User) (datasourceUserType, diag.Diagnostics) {

	var diags, diagnostics diag.Diagnostics
	usersType := datasourceUserType{
		Name: data.Name,
	}

	usersType.Users, diags = mapUsersValuesToListType(ctx, users)
	diagnostics.Append(diags...)

	return usersType, diagnostics
}

// Sets the terraform struct values from the user resources returned by the cf-client
func mapSpaceOrgUsersValuesToType(ctx context.Context, data spaceOrgUsersType, users []*resource.User) (spaceOrgUsersType, diag.Diagnostics) {

	var diags, diagnostics diag.Diagnostics
	spaceOrgUsersType := spaceOrgUsersType{}

	if !data.Organization.IsNull() {
		spaceOrgUsersType.Organization = data.Organization
	}

	if !data.Space.IsNull() {
		spaceOrgUsersType.Space = data.Space
	}

	spaceOrgUsersType.Users, diags = mapUsersValuesToListType(ctx, users)
	diagnostics.Append(diags...)

	return spaceOrgUsersType, diagnostics
}
