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

	if value.Name != nil {
		serviceCredentialBindingType.Name = types.StringValue(value.Name)
	}

	if value.Relationships.App != nil {
		serviceCredentialBindingType.Name = types.StringValue(value.Name)
	}

	serviceCredentialBindingType.Labels, diags = mapMetadataValueToType(ctx, value.Metadata.Labels)
	diagnostics.Append(diags...)
	serviceCredentialBindingType.Annotations, diags = mapMetadataValueToType(ctx, value.Metadata.Annotations)
	diagnostics.Append(diags...)
	serviceCredentialBindingType.LastOperation, diags = types.ObjectValueFrom(ctx, lastOperationAttrTypes, mapLastOperation(value.LastOperation))
	diagnostics.Append(diags...)

	return serviceCredentialBindingType, diagnostics
}
