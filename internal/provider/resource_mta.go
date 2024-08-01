package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/SAP/terraform-provider-cloudfoundry/internal/mta"
	"github.com/SAP/terraform-provider-cloudfoundry/internal/provider/managers"
	"github.com/SAP/terraform-provider-cloudfoundry/internal/validation"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &mtaResource{}
	_ resource.ResourceWithConfigure = &mtaResource{}
)

func NewMtaResource() resource.Resource {
	return &mtaResource{}
}

type mtaResource struct {
	mtaClient *mta.APIClient
}

func (r *mtaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mta"
}

func (r *mtaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Allows deploying applications and services via an MTAR archive or URL.
		
__Further documentation:__ 
 [Multitarget Applications in the Cloud Foundry Environment](https://help.sap.com/docs/btp/sap-business-technology-platform/multitarget-applications-in-cloud-foundry-environment)
`,
		Attributes: map[string]schema.Attribute{
			"mtar_path": schema.StringAttribute{
				MarkdownDescription: "The local path where the MTA archive is present. Either this attribute or mtar_url need to be set.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("mtar_path"),
						path.MatchRoot("mtar_url"),
					}...),
				},
			},
			"mtar_url": schema.StringAttribute{
				MarkdownDescription: "The remote URL where the MTA archive is present",
				Optional:            true,
			},
			"extension_descriptors": schema.SetAttribute{
				MarkdownDescription: "The paths for the MTA deployment extension files.",
				Optional:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"space": schema.StringAttribute{
				MarkdownDescription: "The GUID of the space where the MTA will be deployed",
				Required:            true,
				Validators: []validator.String{
					validation.ValidUUID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"deploy_url": schema.StringAttribute{
				MarkdownDescription: "The URL of the deploy service, if a custom one has been used(should be present in the same landscape). By default 'deploy-service.<system-domain>'",
				Optional:            true,
			},
			"namespace": schema.StringAttribute{
				MarkdownDescription: "The namespace of the MTA. Should be of valid host format",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The MTA ID of the deployment",
				Computed:            true,
			},
			"mta": schema.SingleNestedAttribute{
				MarkdownDescription: "contains the details of the MTA object",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"metadata": schema.SingleNestedAttribute{
						MarkdownDescription: "an identifier, version and namespace that uniquely identify the MTA",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed: true,
							},
							"version": schema.StringAttribute{
								Computed: true,
							},
							"namespace": schema.StringAttribute{
								Computed: true,
							},
						},
					},
					"modules": schema.ListNestedAttribute{
						MarkdownDescription: "the deployable parts contained in the MTA deployment archive, most commonly Cloud Foundry applications or content",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"module_name": schema.StringAttribute{
									Computed: true,
								},
								"app_name": schema.StringAttribute{
									Computed: true,
								},
								"created_on": schema.StringAttribute{
									Computed: true,
								},
								"updated_on": schema.StringAttribute{
									Computed: true,
								},
								"provided_dendency_names": schema.ListAttribute{
									ElementType: types.StringType,
									Computed:    true,
								},
								"services": schema.ListAttribute{
									ElementType: types.StringType,
									Computed:    true,
								},
								"uris": schema.ListAttribute{
									ElementType: types.StringType,
									Computed:    true,
								},
							},
						},
					},
					"services": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
		},
	}
}

func (r *mtaResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	session, ok := req.ProviderData.(*managers.Session)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *managers.Session, got: %T. Please report this issue to the provider developers", req.ProviderData),
		)
		return
	}

	apiEndpointURL := session.CFClient.ApiURL("")
	conf := mta.NewConfiguration(apiEndpointURL, session.CFClient.UserAgent(), session.CFClient.HTTPAuthClient())
	r.mtaClient = mta.NewAPIClient(conf)

	subDomainWithProtocol := strings.Split(apiEndpointURL, ".")[0]
	subDomain := strings.Split(subDomainWithProtocol, "//")[1]
	deploySubdomainWithProtocol := strings.Replace(subDomainWithProtocol, subDomain, "deploy-service", 1)
	deployURL := strings.Replace(apiEndpointURL, subDomainWithProtocol, deploySubdomainWithProtocol, 1)

	r.mtaClient.ChangeBasePath(deployURL)
}

func (r *mtaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.upsert(ctx, &req.Plan, &resp.State, &resp.Diagnostics)
}

func (r *mtaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.upsert(ctx, &req.Plan, &resp.State, &resp.Diagnostics)
}

func (r *mtaResource) upsert(ctx context.Context, reqPlan *tfsdk.Plan, respState *tfsdk.State, respDiags *diag.Diagnostics) {
	var (
		mtarType             MtarType
		uploadedFile         mta.FileMetadata
		err                  error
		mtaId                string
		extensionDescriptors string
	)
	diags := reqPlan.Get(ctx, &mtarType)
	respDiags.Append(diags...)
	if respDiags.HasError() {
		return
	}

	spaceGuid := mtarType.Space.ValueString()
	namespace := mtarType.Namespace.ValueString()

	if !mtarType.DeployUrl.IsNull() {
		r.mtaClient.ChangeBasePath(mtarType.DeployUrl.ValueString())
	}

	if !mtarType.MtarPath.IsNull() {
		fileLocation := mtarType.MtarPath.ValueString()

		uploadedFile, _, err = r.mtaClient.DefaultApi.UploadMtaFile(ctx, spaceGuid, namespace, fileLocation)
		if err != nil {
			respDiags.AddError(
				"Unable to upload mtar file",
				fmt.Sprintf("Request failed with %s ", err.Error()),
			)
			return
		}

		// Extract mta id from archive file
		descriptor, err := mta.GetMtaDescriptorFromArchive(fileLocation)
		if err != nil {
			respDiags.AddError(
				"MTA ID missing",
				fmt.Sprintf("Could not get MTA ID from deployment descriptor %s ", err),
			)
			return
		}
		mtaId = descriptor.ID
	}

	if !mtarType.MtarUrl.IsNull() {
		fileLocation := mtarType.MtarUrl.ValueString()
		uploadJobID, uploadResp, err := r.mtaClient.DefaultApi.AsyncUploadFileFromURL(ctx, spaceGuid, namespace, fileLocation)
		if err != nil {
			respDiags.AddError(
				"Unable to upload remote mtar file",
				fmt.Sprintf("Request failed with %s ", err.Error()),
			)
			return
		}

		jobResponse, err := mta.PollMtaJob(ctx, r.mtaClient, spaceGuid, uploadJobID, mta.FinishedState, uploadResp.Header.Get("x-cf-app-instance"), namespace)
		if err != nil {
			respDiags.AddError(
				"Unable to poll MTAR upload job",
				fmt.Sprintf("Request failed with %s ", err.Error()),
			)
			return
		}
		mtaId = jobResponse.MtaId
		uploadedFile = jobResponse.File
	}

	if !mtarType.ExtensionDescriptors.IsNull() {
		var (
			extensionDescriptorsList []string
			extensionFileID          []string
		)
		diags = mtarType.ExtensionDescriptors.ElementsAs(ctx, &extensionDescriptorsList, false)
		respDiags.Append(diags...)

		for _, descriptorLocation := range extensionDescriptorsList {
			uploadedExtensionDescriptor, _, err := r.mtaClient.DefaultApi.UploadMtaFile(ctx, spaceGuid, namespace, descriptorLocation)
			if err != nil {
				respDiags.AddError(
					"Unable to upload mta extension descriptor",
					fmt.Sprintf("Request failed with %s ", err.Error()),
				)
			}
			extensionFileID = append(extensionFileID, uploadedExtensionDescriptor.Id)
		}
		extensionDescriptors = strings.Join(extensionFileID, ",")
	}

	// Check for an ongoing operation for this MTA ID and abort it
	_, err = mta.CheckOngoingOperation(ctx, r.mtaClient, mtaId, uploadedFile.Namespace, spaceGuid)
	if err != nil {
		respDiags.AddError(
			"Unable to check for and abort ongoing MTA operation",
			fmt.Sprintf("Request failed with %s ", err.Error()),
		)
		return
	}

	operationParams := mta.Operation{
		ProcessType: "DEPLOY",
		Namespace:   namespace,
		Parameters: map[string]interface{}{
			"appArchiveId": uploadedFile.Id,
			"mtaId":        mtaId,
		},
	}

	if extensionDescriptors != "" {
		operationParams.Parameters["mtaExtDescriptorId"] = extensionDescriptors
	}

	//Starting deploy operation
	operationId, _, _, err := r.mtaClient.DefaultApi.StartMtaOperation(ctx, spaceGuid, operationParams)
	if err != nil {
		respDiags.AddError(
			"Unable to start MTA DEPLOY operation",
			fmt.Sprintf("Request failed with %s ", err.Error()),
		)
		return
	}

	err = mta.PollMtaOperation(ctx, r.mtaClient, spaceGuid, operationId, mta.FinishedState)
	if err != nil {
		respDiags.AddError(
			"Failure in polling MTA operation",
			fmt.Sprintf("Request failed with %s ", err.Error()),
		)
		return
	}

	//get details of MTA
	mtaObject, _, err := r.mtaClient.DefaultApi.GetMta(ctx, spaceGuid, mtaId, namespace)
	if err != nil {
		respDiags.AddError(
			"Unable to fetch MTA details",
			fmt.Sprintf("Request failed with %s ", err.Error()),
		)
		return
	}
	mtaTfType, diags := mapMtaValuesToType(ctx, mtaObject)
	respDiags.Append(diags...)
	mtarType.Mta, diags = types.ObjectValueFrom(ctx, mtaObjAttributes, mtaTfType)
	respDiags.Append(diags...)
	mtarType.Id = types.StringValue(mtaObject.Metadata.Id)
	respDiags.Append(respState.Set(ctx, mtarType)...)
}

func (r *mtaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MtarType
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.DeployUrl.IsNull() {
		r.mtaClient.ChangeBasePath(data.DeployUrl.ValueString())
	}

	//get details of MTA
	mtaObject, _, err := r.mtaClient.DefaultApi.GetMta(ctx, data.Space.ValueString(), data.Id.ValueString(), data.Namespace.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), mta.MTA_NOT_FOUND) {
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError(
				"Unable to fetch MTA details",
				fmt.Sprintf("Request failed with %s ", err.Error()),
			)
			return
		}

	}

	mtaTfType, diags := mapMtaValuesToType(ctx, mtaObject)
	resp.Diagnostics.Append(diags...)
	data.Mta, diags = types.ObjectValueFrom(ctx, mtaObjAttributes, mtaTfType)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "read an mtar resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mtaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var mtarType MtarType
	diags := req.State.Get(ctx, &mtarType)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mtaId := mtarType.Id.ValueString()
	spaceGuid := mtarType.Space.ValueString()

	if !mtarType.DeployUrl.IsNull() {
		r.mtaClient.ChangeBasePath(mtarType.DeployUrl.ValueString())
	}

	// Check for an ongoing operation for this MTA ID and abort it
	_, err := mta.CheckOngoingOperation(ctx, r.mtaClient, mtaId, mtarType.Namespace.ValueString(), spaceGuid)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to check for and abort ongoing MTA operation",
			fmt.Sprintf("Request failed with %s ", err.Error()),
		)
		return
	}

	operationParams := mta.Operation{
		ProcessType: "UNDEPLOY",
		Namespace:   mtarType.Namespace.ValueString(),
		Parameters: map[string]interface{}{
			"mtaId":          mtaId,
			"deleteServices": true,
		},
	}

	operationId, _, _, err := r.mtaClient.DefaultApi.StartMtaOperation(ctx, spaceGuid, operationParams)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to start MTA UNDEPLOY operation",
			fmt.Sprintf("Request failed with %s ", err.Error()),
		)
		return
	}

	err = mta.PollMtaOperation(ctx, r.mtaClient, spaceGuid, operationId, mta.FinishedState)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failure in polling MTA operation",
			fmt.Sprintf("Request failed with %s ", err.Error()),
		)
		return
	}
}

func (r *mtaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) > 3 || len(parts) < 2 {
		resp.Diagnostics.AddError(
			"Resource Import ID of Invalid format",
			"The format for import ID should be of [space_guid]/[mta_id] OR [space_guid]/[mta_id]/[namespace]",
		)
		return
	}
	spaceGuid := parts[0]
	mtaId := parts[1]

	if len(parts) == 3 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("namespace"), parts[2])...)
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("space"), spaceGuid)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), mtaId)...)
}
