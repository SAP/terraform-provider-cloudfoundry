package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/SAP/terraform-provider-cloudfoundry/internal/provider/managers"
	cfv3client "github.com/cloudfoundry-community/go-cfclient/v3/client"
	cfv3operation "github.com/cloudfoundry-community/go-cfclient/v3/operation"
	cfv3resource "github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gopkg.in/yaml.v2"
)

var (
	_ resource.Resource              = &appResource{}
	_ resource.ResourceWithConfigure = &appResource{}
)

func NewAppResource() resource.Resource {
	return &appResource{}
}

type appResource struct {
	cfClient *cfv3client.Client
}

func (r *appResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app"
}

func (r *appResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a Cloud Foundry resource to manage applications.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the application.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"space": schema.StringAttribute{
				MarkdownDescription: "The name of the associated Cloud Foundry space.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"org": schema.StringAttribute{
				MarkdownDescription: "The name of the associated Cloud Foundry organization.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"stack": schema.StringAttribute{
				MarkdownDescription: "The name of the stack the application will be deployed to.",
				Optional:            true,
				Computed:            true,
			},
			"buildpacks": schema.SetAttribute{
				MarkdownDescription: "Multiple buildpacks used to stage the application.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "An uri or path to target a zip file. this can be in the form of unix path (/my/path.zip) or url path (http://zip.com/my.zip).",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_code_hash": schema.StringAttribute{
				MarkdownDescription: "Used to trigger updates. Must be set to a base64-encoded SHA256 hash of the path specified.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"docker_image": schema.StringAttribute{
				MarkdownDescription: "The URL to the docker image with tag e.g registry.example.com:5000/user/repository/tag or docker image name from the public repo e.g. redis:4.0",
				Optional:            true,
			},
			"docker_credentials": schema.MapNestedAttribute{
				MarkdownDescription: "Defines login credentials for private docker repositories",
				Optional:            true,
				Validators: []validator.Map{
					mapvalidator.AlsoRequires(path.MatchRoot("docker_image")),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"username": schema.StringAttribute{
							MarkdownDescription: "The username for the private docker repository. Password must be provided in the environment variable CF_DOCKER_PASSWORD.",
							Required:            true,
							Sensitive:           true,
						},
					},
				},
			},
			"strategy": schema.StringAttribute{
				MarkdownDescription: "The deployment strategy to use when deploying the application. Valid values are 'none', 'rolling', and 'blue-green', defaults to 'none'.",
				Optional:            true,
				Default:             stringdefault.StaticString("none"),
				Validators: []validator.String{
					stringvalidator.OneOf("none", "rolling", "blue-green"),
				},
			},
			"service_bindings": schema.SetNestedAttribute{
				MarkdownDescription: "Service instances to bind to the application.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"service_instance": schema.StringAttribute{
							MarkdownDescription: "The service instance GUID.",
							Required:            true,
						},
						"params": schema.MapAttribute{
							MarkdownDescription: "A list of key/value parameters used by the service broker to create the binding.",
							Optional:            true,
						},
					},
				},
			},
			"routes": schema.SetNestedAttribute{
				MarkdownDescription: "The routes to map to the application to control its ingress traffic.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"route": schema.StringAttribute{
							MarkdownDescription: "The route GUID.",
							Required:            true,
						},
						"protocol": schema.StringAttribute{
							MarkdownDescription: "The protocol to use for the route. Valid values are http2, http1, and tcp.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("http2", "http1", "tcp"),
							},
						},
					},
				},
			},
			"environment": schema.MapAttribute{
				MarkdownDescription: "Key/value pairs of custom environment variables to set in your app. Does not include any system or service variables.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"health_check_interval": schema.Int64Attribute{
				MarkdownDescription: "The interval in seconds between health checks.",
			},
			"readiness_health_check_type": schema.StringAttribute{
				MarkdownDescription: "The readiness health check type which can be one of 'port', 'process', 'http'.",
				Optional:            true,
			},
			"readiness_health_check_http_endpoint": schema.StringAttribute{
				MarkdownDescription: "The endpoint for the http readiness health check type.",
				Optional:            true,
			},
			"readiness_health_check_invocation_timeout": schema.Int64Attribute{
				MarkdownDescription: "The timeout in seconds for the readiness health check requests for http and port health checks.",
				Optional:            true,
			},
			"readiness_health_check_interval": schema.Int64Attribute{
				MarkdownDescription: "The interval in seconds between readiness health checks.",
				Optional:            true,
			},
			"log_rate_limit_per_second": schema.StringAttribute{
				MarkdownDescription: "The attribute specifies the log rate limit for all instances of an app.",
				Optional:            true,
			},
			"no_route": schema.BoolAttribute{
				MarkdownDescription: "The attribute with a value of true to prevent a route from being created for your app.",
				Optional:            true,
			},
			"random_route": schema.BoolAttribute{
				MarkdownDescription: "The random-route attribute to generate a unique route and avoid name collisions.",
				Optional:            true,
			},
			"processes": schema.SetNestedAttribute{
				MarkdownDescription: "List of configurations for individual process types.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: r.ProcessSchemaAttributes(),
				},
			},
			"sidecars": schema.SetNestedAttribute{
				MarkdownDescription: "The attribute specifies additional processes to run in the same container as your app",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Sidecar name. The identifier for the sidecars to be configured.",
							Required:            true,
						},
						"command": schema.StringAttribute{
							MarkdownDescription: "The command used to start the sidecar.",
							Required:            true,
						},
						"process_types": schema.SetAttribute{
							MarkdownDescription: "List of processes to associate sidecar with.",
							ElementType:         types.StringType,
							Required:            true,
							Validators: []validator.Set{
								setvalidator.ValueStringsAre(stringvalidator.OneOf("web", "worker")),
							},
						},
						"memory": schema.StringAttribute{
							MarkdownDescription: "The memory limit for the sidecar.",
							Optional:            true,
						},
					},
				},
			},
			idKey:        guidSchema(),
			createdAtKey: createdAtSchema(),
			updatedAtKey: updatedAtSchema(),
		},
	}
	for k, v := range r.ProcessAppCommonSchema() {
		if _, ok := resp.Schema.Attributes[k]; !ok {
			resp.Schema.Attributes[k] = v
		}
	}
}

func (r *appResource) ProcessSchemaAttributes() map[string]schema.Attribute {
	pSchema := map[string]schema.Attribute{
		"type": schema.StringAttribute{
			MarkdownDescription: "The process type. Can be web or worker.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("web", "worker"),
			},
		},
	}
	for k, v := range r.ProcessAppCommonSchema() {
		if _, ok := pSchema[k]; !ok {
			pSchema[k] = v
		}
	}
	return pSchema
}
func (r *appResource) ProcessAppCommonSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"command": schema.StringAttribute{
			MarkdownDescription: "A custom start command for the application. This overrides the start command provided by the buildpack.",
			Optional:            true,
		},
		"disk_quota": schema.StringAttribute{
			MarkdownDescription: "The disk space to be allocated for each application instance.",
			Optional:            true,
			Computed:            true,
		},
		"health_check_http_endpoint": schema.StringAttribute{
			MarkdownDescription: "The endpoint for the http health check type.",
			Optional:            true,
		},
		"health_check_invocation_timeout": schema.Int64Attribute{
			MarkdownDescription: "The timeout in seconds for the health check requests for http and port health checks.",
			Optional:            true,
		},
		"health_check_type": schema.StringAttribute{
			MarkdownDescription: "The health check type which can be one of 'port', 'process', 'http'.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("port", "process", "http"),
			},
		},
		"instances": schema.Int64Attribute{
			MarkdownDescription: "The number of app instances that you want to start. Defaults to 1.",
			Optional:            true,
			Default:             int64default.StaticInt64(int64(1)),
		},
		"memory": schema.Int64Attribute{
			MarkdownDescription: "The memory limit for each application instance in megabytes. If not provided, value is computed and retreived from Cloud Foundry.",
			Optional:            true,
			Computed:            true,
		},
		"timeout": schema.Int64Attribute{
			MarkdownDescription: "Defines the number of seconds that Cloud Foundry allocates for starting your app.",
			Optional:            true,
		},
	}
}
func (r *appResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	session, ok := req.ProviderData.(*managers.Session)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *managers.Session, got: %T. Please report this issue to the provider developers", req.ProviderData),
		)
		return
	}
	r.cfClient = session.CFClient
}

func (r *appResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.upsert(ctx, &req.Plan, &resp.State, &resp.Diagnostics)
}

func (r *appResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var appType AppType
	diags := req.State.Get(ctx, &appType)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	appRaw, err := r.cfClient.Manifests.Generate(ctx, appType.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading app", err.Error())
		return
	}
	var appManifest cfv3operation.Manifest
	err = yaml.Unmarshal([]byte(appRaw), &appManifest)
	if err != nil {
		resp.Diagnostics.AddError("Error unmarshalling app", err.Error())
		return
	}
	appResp, err := r.cfClient.Applications.Get(ctx, appType.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error fetching app", err.Error())
		return
	}
	plan, diags := mapAppValuesToType(ctx, appManifest.Applications[0], appResp)
	resp.Diagnostics.Append(diags...)
	resp.State.Set(ctx, &plan)
}

func (r *appResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.upsert(ctx, &req.Plan, &resp.State, &resp.Diagnostics)
}
func (r *appResource) upsert(ctx context.Context, reqPlan *tfsdk.Plan, respState *tfsdk.State, respDiags *diag.Diagnostics) {
	var appType AppType
	diags := reqPlan.Get(ctx, &appType)
	respDiags.Append(diags...)
	if respDiags.HasError() {
		return
	}

	appManifestValue, diags := appType.mapAppTypeToValues(ctx)
	respDiags.Append(diags...)
	if respDiags.HasError() {
		return
	}
	appResp, err := r.push(appType, appManifestValue, ctx)
	if err != nil {
		respDiags.AddError("Error pushing app", err.Error())
		return
	}
	manifestRespRaw, err := r.cfClient.Manifests.Generate(ctx, appResp.GUID)
	if err != nil {
		respDiags.AddError("Error generating manifest", err.Error())
	}
	var manifest *cfv3operation.Manifest
	err = yaml.Unmarshal([]byte(manifestRespRaw), &manifest)
	if err != nil {
		respDiags.AddError("Error unmarshalling manifest", err.Error())
	}
	plan, diags := mapAppValuesToType(ctx, manifest.Applications[0], appResp)
	respDiags.Append(diags...)
	plan.Space = appType.Space
	plan.Org = appType.Org
	plan.Path = appType.Path
	plan.Strategy = appType.Strategy
	plan.SourceCodeHash = appType.SourceCodeHash
	respDiags.Append(respState.Set(ctx, &plan)...)
}
func (r *appResource) push(appType AppType, appManifestValue *cfv3operation.AppManifest, ctx context.Context) (*cfv3resource.App, error) {
	file, err := os.Open(appType.Path.ValueString())
	if err != nil {
		return nil, err
	}
	manifestOp := cfv3operation.NewAppPushOperation(r.cfClient, appType.Org.ValueString(), appType.Space.ValueString())
	if !appType.Strategy.IsNull() {
		var sm cfv3operation.StrategyMode
		switch appType.Strategy.ValueString() {
		case "rolling":
			sm = cfv3operation.StrategyRolling
		case "blue-green":
			sm = cfv3operation.StrategyBlueGreen
		default:
			sm = cfv3operation.StrategyNone
		}
		manifestOp.WithStrategy(sm)
	}
	appResp, err := manifestOp.Push(ctx, appManifestValue, file)
	if err != nil {
		return nil, err
	}
	return appResp, nil
}

func (r *appResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var appType AppType
	diags := req.State.Get(ctx, &appType)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	jobID, err := r.cfClient.Applications.Delete(ctx, appType.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete application",
			fmt.Sprintf("Application deletion verification failed %s with %s", appType.Name.ValueString(), err.Error()),
		)
		return
	}
	err = pollJob(ctx, *r.cfClient, jobID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to verify org quota deletion",
			"Application deletion verification failed for "+appType.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}

func (r *appResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
