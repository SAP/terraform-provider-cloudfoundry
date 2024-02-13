package provider

import (
	"context"
	"time"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type serviceInstanceType struct {
	Name             types.String         `tfsdk:"name"`
	ID               types.String         `tfsdk:"id"`
	Type             types.String         `tfsdk:"type"`
	Space            types.String         `tfsdk:"space"`
	ServicePlan      types.String         `tfsdk:"service_plan"`
	Parameters       jsontypes.Normalized `tfsdk:"parameters"`
	LastOperation    types.List           `tfsdk:"last_operation"` //LastOperationType
	Tags             types.Set            `tfsdk:"tags"`
	DashboardURL     types.String         `tfsdk:"dashboard_url"`
	Credentials      jsontypes.Normalized `tfsdk:"credentials"`
	SyslogDrainURL   types.String         `tfsdk:"syslog_drain_url"`
	RouteServiceURL  types.String         `tfsdk:"route_service_url"`
	MaintenanceInfo  types.List           `tfsdk:"maintenance_info"` //maintenanceInfoType
	UpgradeAvailable types.Bool           `tfsdk:"upgrade_available"`
	Labels           types.Map            `tfsdk:"labels"`
	Annotations      types.Map            `tfsdk:"annotations"`
	CreatedAt        types.String         `tfsdk:"created_at"`
	UpdatedAt        types.String         `tfsdk:"updated_at"`
}

type lastOperationType struct {
	Type        types.String `tfsdk:"type"`
	State       types.String `tfsdk:"state"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

type maintenanceInfoType struct {
	Version     types.String `tfsdk:"version"`
	Description types.String `tfsdk:"description"`
}

func mapServiceInstanceValuesToType(ctx context.Context, value *resource.ServiceInstance, parameters jsontypes.Normalized) (serviceInstanceType, diag.Diagnostics) {

	serviceInstanceType := serviceInstanceType{
		Name:        types.StringValue(value.Name),
		ID:          types.StringValue(value.GUID),
		Type:        types.StringValue(value.Type),
		Space:       types.StringValue(value.Relationships.Space.Data.GUID),
		ServicePlan: types.StringValue(value.Relationships.ServicePlan.Data.GUID),
		CreatedAt:   types.StringValue(value.CreatedAt.Format(time.RFC3339)),
		UpdatedAt:   types.StringValue(value.UpdatedAt.Format(time.RFC3339)),
	}
	if !parameters.IsNull() {
		serviceInstanceType.Parameters = jsontypes.NewNormalizedValue(parameters.ValueString())
	} else {
		serviceInstanceType.Parameters = jsontypes.NewNormalizedNull()
	}
	if value.DashboardURL != nil {
		serviceInstanceType.DashboardURL = types.StringValue(*value.DashboardURL)
	}
	if value.SyslogDrainURL != nil {
		serviceInstanceType.SyslogDrainURL = types.StringValue(*value.SyslogDrainURL)
	}
	if value.RouteServiceURL != nil {
		serviceInstanceType.RouteServiceURL = types.StringValue(*value.RouteServiceURL)
	}
	if value.UpgradeAvailable != nil {
		serviceInstanceType.UpgradeAvailable = types.BoolValue(*value.UpgradeAvailable)
	}

	var diags, diagnostics diag.Diagnostics
	serviceInstanceType.Labels, diags = mapMetadataValueToType(ctx, value.Metadata.Labels)
	diagnostics.Append(diags...)
	serviceInstanceType.Annotations, diags = mapMetadataValueToType(ctx, value.Metadata.Annotations)
	diagnostics.Append(diags...)

	// tags mapping
	if len(value.Tags) > 0 {
		serviceInstanceType.Tags, diags = types.SetValueFrom(ctx, types.StringType, value.Tags)
		diagnostics.Append(diags...)
	} else {
		var empty []string
		serviceInstanceType.Tags, diags = types.SetValueFrom(ctx, types.StringType, empty)
		diagnostics.Append(diags...)
	}
	maintenanceInfo, diags := flattenMaintenanceInfo(ctx, value.MaintenanceInfo)
	diagnostics.Append(diags...)
	serviceInstanceType.MaintenanceInfo = *maintenanceInfo

	lastOperation, diags := flattenLastOperation(ctx, &value.LastOperation)
	diagnostics.Append(diags...)
	serviceInstanceType.LastOperation = *lastOperation

	return serviceInstanceType, diagnostics
}

func flattenMaintenanceInfo(ctx context.Context, maintenanceInfo *resource.ServiceInstanceMaintenanceInfo) (*basetypes.ListValue, diag.Diagnostics) {
	result := make(map[string]attr.Value)

	result["version"] = types.StringValue(maintenanceInfo.Version)
	result["description"] = types.StringValue(maintenanceInfo.Description)

	obj, diags := types.ObjectValue(maintenanceInfoAttrTypes.AttrTypes, result)
	if diags.HasError() {
		return nil, diags
	}
	objList := []attr.Value{obj}

	resultList, diag := basetypes.NewListValue(
		maintenanceInfoAttrTypes,
		objList,
	)
	if diag.HasError() {
		return nil, diag
	}

	return &resultList, nil
}

var maintenanceInfoAttrTypes = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"version":     types.StringType,
		"description": types.StringType,
	},
}

var lastOperationAttrTypes = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"type":        types.StringType,
		"state":       types.StringType,
		"description": types.StringType,
		"created_at":  types.StringType,
		"updated_at":  types.StringType,
	},
}

func flattenLastOperation(ctx context.Context, lastOperation *resource.LastOperation) (*basetypes.ListValue, diag.Diagnostics) {
	result := make(map[string]attr.Value)

	result["type"] = types.StringValue(lastOperation.Type)
	result["state"] = types.StringValue(lastOperation.State)
	result["description"] = types.StringValue(lastOperation.Description)
	result["created_at"] = types.StringValue(lastOperation.CreatedAt.Format(time.RFC3339))
	result["updated_at"] = types.StringValue(lastOperation.UpdatedAt.Format(time.RFC3339))

	obj, diags := types.ObjectValue(lastOperationAttrTypes.AttrTypes, result)
	if diags.HasError() {
		return nil, diags
	}
	objList := []attr.Value{obj}

	resultList, diag := basetypes.NewListValue(
		lastOperationAttrTypes,
		objList,
	)
	if diag.HasError() {
		return nil, diag
	}

	return &resultList, nil
}

// ToAliasList converts aliases of type types.Set into a slice of strings.
func toTagsList(ctx context.Context, tagsSet types.Set) ([]string, diag.Diagnostics) {
	if tagsSet.IsNull() || tagsSet.IsUnknown() {
		return []string{}, nil
	}

	tags := make([]string, 0, len(tagsSet.Elements()))
	diags := tagsSet.ElementsAs(ctx, &tags, false)
	return tags, diags
}
