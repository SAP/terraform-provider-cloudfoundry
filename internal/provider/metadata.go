package provider

import (
	"context"
	"fmt"
	"time"

	cfv3client "github.com/cloudfoundry-community/go-cfclient/v3/client"
	cfv3resource "github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/samber/lo"
)

const (
	idKey          = "id"
	labelsKey      = "labels"
	annotationsKey = "annotations"
	createdAtKey   = "created_at"
	updatedAtKey   = "updated_at"
)

const defaultTimeout = 20 * time.Minute

func datasourceLabelsSchema() *schema.MapAttribute {
	return &schema.MapAttribute{
		MarkdownDescription: "The labels associated with Cloud Foundry resources. Add as described [here](https://docs.cloudfoundry.org/adminguide/metadata.html#-view-metadata-for-an-object).",
		ElementType:         types.StringType,
		Computed:            true,
	}
}

func datasourceAnnotationsSchema() *schema.MapAttribute {
	return &schema.MapAttribute{
		MarkdownDescription: "The annotations associated with Cloud Foundry resources.Add as described [here](https://docs.cloudfoundry.org/adminguide/metadata.html#-view-metadata-for-an-object).",
		ElementType:         types.StringType,
		Computed:            true,
	}
}

func resourceLabelsSchema() *schema.MapAttribute {
	return &schema.MapAttribute{
		MarkdownDescription: `The labels associated with Cloud Foundry resources. Add as described [here](https://docs.cloudfoundry.org/adminguide/metadata.html#-view-metadata-for-an-object).`,
		ElementType:         types.StringType,
		Optional:            true,
	}
}

func resourceAnnotationsSchema() *schema.MapAttribute {
	return &schema.MapAttribute{
		MarkdownDescription: "The annotations associated with Cloud Foundry resources. Add as described [here](https://docs.cloudfoundry.org/adminguide/metadata.html#-view-metadata-for-an-object).",
		ElementType:         types.StringType,
		Optional:            true,
	}
}

func createdAtSchema() *schema.StringAttribute {
	return &schema.StringAttribute{
		MarkdownDescription: "The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}
func updatedAtSchema() *schema.StringAttribute {
	return &schema.StringAttribute{
		MarkdownDescription: "The date and time when the resource was updated in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
		Computed:            true,
	}
}
func guidSchema() *schema.StringAttribute {
	return &schema.StringAttribute{
		MarkdownDescription: "The GUID of the object.",
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

// Take relationship from cfclient and return set type of terraform
func setRelationshipToTFSet(r []cfv3resource.Relationship) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	var bt basetypes.SetValue
	if len(r) != 0 {
		tfVal := []attr.Value{}
		for _, val := range r {
			tfVal = append(tfVal, types.StringValue(val.GUID))
		}
		bt, diags = types.SetValue(types.StringType, tfVal)
	} else {
		bt = types.SetNull(types.StringType)
	}
	return bt, diags
}

// Returns removed and added element in the new plan which existed in state
func findChangedRelationsFromTFState(ctx context.Context, planSet basetypes.SetValue, stateSet basetypes.SetValue) ([]string, []string, diag.Diagnostics) {
	var diags diag.Diagnostics
	var planSetStr, stateSetStr []string
	diags = append(diags, planSet.ElementsAs(ctx, &planSetStr, false)...)
	diags = append(diags, stateSet.ElementsAs(ctx, &stateSetStr, false)...)
	removed, added := lo.Difference(stateSetStr, planSetStr)
	return removed, added, diags
}
func setMapToBaseMap(ctx context.Context, resp *datasource.ReadResponse, mt map[string]*string) *basetypes.MapValue {
	labels, diag := types.MapValueFrom(ctx, types.StringType, mt)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return nil
	}
	return &labels
}

func handleReadErrors(ctx context.Context, resp *resource.ReadResponse, err error, res string, resName string) {
	if cfv3resource.IsResourceNotFoundError(err) {
		resp.State.RemoveResource(ctx)
	} else {
		resp.Diagnostics.AddError(fmt.Sprintf("API Error Reading %s %s", res, resName), fmt.Sprintf("%s", err.Error()))
	}

}

func pollJob(ctx context.Context, client cfv3client.Client, jobID string) error {

	return client.Jobs.PollComplete(ctx, jobID, &cfv3client.PollingOptions{
		Timeout:       defaultTimeout,
		CheckInterval: time.Second * 10,
		FailedState:   string(cfv3resource.JobStateFailed),
	})
}

func mapMetadataValueToType(ctx context.Context, generic map[string]*string) (basetypes.MapValue, diag.Diagnostics) {

	var out basetypes.MapValue
	var diagnostics diag.Diagnostics
	if len(generic) == 0 {
		out = types.MapNull(types.StringType)
	} else {
		out, diagnostics = types.MapValueFrom(ctx, types.StringType, generic)
	}

	return out, diagnostics
}
