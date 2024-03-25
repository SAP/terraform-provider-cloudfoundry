package provider

import (
	"context"
	"time"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type serviceCredentialBindingType struct {
	Name            types.String `tfsdk:"name"`
	ID              types.String `tfsdk:"id"`
	Type            types.String `tfsdk:"type"`
	ServiceInstance types.String `tfsdk:"service_instance"`
	App             types.String `tfsdk:"app"`
	LastOperation   types.Object `tfsdk:"last_operation"` //LastOperationType
	Labels          types.Map    `tfsdk:"labels"`
	Annotations     types.Map    `tfsdk:"annotations"`
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
}

func mapServiceCredentialBindingValuesToType(ctx context.Context, value *resource.ServiceCredentialBinding) (serviceCredentialBindingType, diag.Diagnostics) {
	var diags, diagnostics diag.Diagnostics
	serviceCredentialBindingType := serviceCredentialBindingType{
		ID:        types.StringValue(value.GUID),
		Type:      types.StringValue(value.Type),
		CreatedAt: types.StringValue(value.CreatedAt.Format(time.RFC3339)),
		UpdatedAt: types.StringValue(value.UpdatedAt.Format(time.RFC3339)),
	}

	Name:      types.StringValue(value.Name),

	switch value.Type {
	case managedSerivceInstance:
		dsServiceInstanceType.ServicePlan = types.StringValue(value.Relationships.ServicePlan.Data.GUID)
		if value.DashboardURL != nil {
			dsServiceInstanceType.DashboardURL = types.StringValue(*value.DashboardURL)
		}
	case userProvidedServiceInstance:
		if value.SyslogDrainURL != nil {
			dsServiceInstanceType.SyslogDrainURL = types.StringValue(*value.SyslogDrainURL)
		}
		if value.RouteServiceURL != nil {
			dsServiceInstanceType.RouteServiceURL = types.StringValue(*value.RouteServiceURL)
		}
	}
	dsServiceInstanceType.Labels, diags = mapMetadataValueToType(ctx, value.Metadata.Labels)
	diagnostics.Append(diags...)
	dsServiceInstanceType.Annotations, diags = mapMetadataValueToType(ctx, value.Metadata.Annotations)
	diagnostics.Append(diags...)
	if value.MaintenanceInfo != nil {
		dsServiceInstanceType.MaintenanceInfo, diags = types.ObjectValueFrom(ctx, maintenanceInfoAttrTypes, mapMaintenanceInfo(*value.MaintenanceInfo))
		diagnostics.Append(diags...)
	} else {
		dsServiceInstanceType.MaintenanceInfo = types.ObjectNull(maintenanceInfoAttrTypes)

	}
	dsServiceInstanceType.LastOperation, diags = types.ObjectValueFrom(ctx, lastOperationAttrTypes, mapLastOperation(value.LastOperation))
	diagnostics.Append(diags...)
	//tags mapping
	if len(value.Tags) > 0 {
		tags := make([]types.String, 0, len(value.Tags))
		for _, t := range value.Tags {
			tags = append(tags, types.StringValue(t))
		}
		dsServiceInstanceType.Tags, diags = types.SetValueFrom(ctx, types.StringType, tags)
		diagnostics.Append(diags...)
	} else {
		dsServiceInstanceType.Tags = types.SetNull(types.StringType)

	}

	return dsServiceInstanceType, diagnostics
}
