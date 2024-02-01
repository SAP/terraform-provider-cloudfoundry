package provider

import (
	"context"
	"fmt"
	"time"

	cfv3client "github.com/cloudfoundry-community/go-cfclient/v3/client"
	cfv3resource "github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const (
	labelsKey      = "labels"
	annotationsKey = "annotations"
)

const DefaultTimeout = 20 * time.Minute

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

func pollJob(ctx context.Context, client cfv3client.Client, jobID string, timeout time.Duration) error {

	return client.Jobs.PollComplete(ctx, jobID, &cfv3client.PollingOptions{
		Timeout:       timeout,
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
