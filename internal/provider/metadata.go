package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const (
	labelsKey      = "labels"
	annotationsKey = "annotations"
)

func labelsSchema() *schema.MapAttribute {
	return &schema.MapAttribute{
		MarkdownDescription: "The labels associated with Cloud Foundry resources",
		ElementType:         types.StringType,
		Computed:            true,
	}
}

func annotationsSchema() *schema.MapAttribute {
	return &schema.MapAttribute{
		MarkdownDescription: "The annotations associated with Cloud Foundry resources",
		ElementType:         types.StringType,
		Computed:            true,
	}
}

func setMapToBaseMap(ctx context.Context, resp *datasource.ReadResponse, mt map[string]*string) *basetypes.MapValue {
	labels, diag := types.MapValueFrom(ctx, types.StringType, mt)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return nil
	}
	return &labels
}
