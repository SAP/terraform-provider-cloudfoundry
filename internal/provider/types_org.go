package provider

import (
	"context"
	"time"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type orgType struct {
	Name        types.String `tfsdk:"name"`
	ID          types.String `tfsdk:"id"`
	Labels      types.Map    `tfsdk:"labels"`
	Annotations types.Map    `tfsdk:"annotations"`
	Quota       types.String `tfsdk:"quota"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	Suspended   types.Bool   `tfsdk:"suspended"`
}

func mapOrgValuesToType(ctx context.Context, value *resource.Organization) (orgType, diag.Diagnostics) {
	orgType := orgType{
		Name:      types.StringValue(value.Name),
		ID:        types.StringValue(value.GUID),
		Quota:     types.StringValue(value.Relationships.Quota.Data.GUID),
		CreatedAt: types.StringValue(value.CreatedAt.Format(time.RFC3339)),
		UpdatedAt: types.StringValue(value.UpdatedAt.Format(time.RFC3339)),
		Suspended: types.BoolValue(*value.Suspended),
	}
	var diags, diagnostics diag.Diagnostics
	orgType.Labels, diags = types.MapValueFrom(ctx, types.StringType, value.Metadata.Labels)
	diagnostics.Append(diags...)

	orgType.Annotations, diags = types.MapValueFrom(ctx, types.StringType, value.Metadata.Annotations)
	diagnostics.Append(diags...)

	return orgType, diagnostics
}
