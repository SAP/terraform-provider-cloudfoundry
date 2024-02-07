package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/SAP/terraform-provider-cloudfoundry/internal/provider/managers"
	cfconfig "github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &CloudFoundryProvider{}

type CloudFoundryProvider struct {
	version    string
	httpClient *http.Client
}

type CloudFoundryProviderModel struct {
	Endpoint          types.String `tfsdk:"api_url"`
	User              types.String `tfsdk:"user"`
	Password          types.String `tfsdk:"password"`
	CFClientID        types.String `tfsdk:"cf_client_id"`
	CFClientSecret    types.String `tfsdk:"cf_client_secret"`
	SkipSslValidation types.Bool   `tfsdk:"skip_ssl_validation"`
}

func (p *CloudFoundryProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cloudfoundry"
	resp.Version = p.version
}

func (p *CloudFoundryProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Cloud Foundry Terraform plugin is an integration that allows users to leverage Terraform, an infrastructure as code tool, to define and provision infrastructure resources within the Cloud Foundry platform.",
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				MarkdownDescription: "Specific URL representing the entry point for communication between the client and a Cloud Foundry instance.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "A unique identifier associated with an individual or entity for authentication & authorization purposes.",
				Optional:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "A confidential alphanumeric code associated with a user account on the Cloud Foundry platform",
				Optional:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"cf_client_id": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"cf_client_secret": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"skip_ssl_validation": schema.BoolAttribute{
				Optional: true,
			},
		},
	}
}
func addGenericAttributeError(resp *provider.ConfigureResponse, status string, pathRoot string, commonName string, envName string) {
	resp.Diagnostics.AddAttributeError(
		path.Root(pathRoot),
		fmt.Sprintf("%s field %s", status, pathRoot),
		fmt.Sprintf("The provider cannot create the Cloud Foundry API client as there is an unknown configuration value for the Cloud Foundry %s. "+
			"Either target apply the source of the value first, set the value statically in the configuration, or use the %s environment variable, ensure value is not empty. ", commonName, envName),
	)
}
func addTypeCastAttributeError(resp *provider.ConfigureResponse, expectedType string, pathRoot string, commonName string, envName string) {
	resp.Diagnostics.AddAttributeError(
		path.Root(pathRoot),
		fmt.Sprintf("Expected %s in field %s", expectedType, pathRoot),
		fmt.Sprintf("The provider cannot create the Cloud Foundry API client as there is an invalid configuration value for the Cloud Foundry %s. "+
			"Ensure %s is of type %s ", commonName, envName, expectedType),
	)
}
func checkConfigUnknown(config *CloudFoundryProviderModel, resp *provider.ConfigureResponse) {
	_, cfconfigerr := cfconfig.NewFromCFHome()

	anyParamExists := !config.User.IsUnknown() || !config.Password.IsUnknown() || !config.CFClientID.IsUnknown() || !config.CFClientSecret.IsUnknown()

	// If endpoint is unknown check if any of the auth param exists if yes throw error as api_url is manadatory if not
	// check is cf home directory config exists
	if config.Endpoint.IsUnknown() && (anyParamExists || cfconfigerr != nil) {
		addGenericAttributeError(resp, "Unknown", "api_url", "API Endpoint", "CF_API_URL")
	}
	switch {
	case config.User.IsUnknown() && !config.Password.IsUnknown():
		addGenericAttributeError(resp, "Unknown", "user", "Username", "CF_USER")
	case !config.User.IsUnknown() && config.Password.IsUnknown():
		addGenericAttributeError(resp, "Unknown", "password", "Password", "CF_PASSWORD")
	case config.User.IsUnknown() && config.Password.IsUnknown():
		switch {
		case config.CFClientID.IsUnknown() && !config.CFClientSecret.IsUnknown():
			addGenericAttributeError(resp, "Unknown", "cf_client_id", "CF Client ID", "CF_CF_CLIENT_ID")
		case !config.CFClientID.IsUnknown() && config.CFClientSecret.IsUnknown():
			addGenericAttributeError(resp, "Unknown", "cf_client_secret", "CF Client Secret", "CF_CF_CLIENT_SECRET")
		case config.CFClientID.IsUnknown() && config.CFClientSecret.IsUnknown():
			if !config.Endpoint.IsUnknown() || cfconfigerr != nil {
				resp.Diagnostics.AddError(
					"Unable to create CF Client due to unknown values",
					"Either user/password or client_id/client_secret must be set with api_url or CF config must exist in path (default ~/.cf/config.json)",
				)
			}
		}
	}
}

func checkConfig(resp *provider.ConfigureResponse, endpoint string, user string, password string, cfclientid string, cfclientsecret string) {
	_, cfconfigerr := cfconfig.NewFromCFHome()

	anyParamExists := user != "" || password != "" || cfclientid != "" || cfclientsecret != ""

	if endpoint == "" && (anyParamExists || cfconfigerr != nil) {
		addGenericAttributeError(resp, "Missing", "api_url", "API Endpoint", "CF_API_URL")
	}
	switch {
	case user == "" && password != "":
		addGenericAttributeError(resp, "Missing", "user", "Username", "CF_USER")
	case user != "" && password == "":
		addGenericAttributeError(resp, "Missing", "password", "Password", "CF_PASSWORD")
	case user == "" && password == "":
		switch {
		case cfclientid == "" && cfclientsecret != "":
			addGenericAttributeError(resp, "Missing", "cf_client_id", "Client ID", "CF_CF_CLIENT_ID")
		case cfclientid != "" && cfclientsecret == "":
			addGenericAttributeError(resp, "Missing", "cf_client_secret", " Client Secret", "CF_CF_CLIENT_SECRET")
		case cfclientid == "" && cfclientsecret == "":
			if endpoint != "" || cfconfigerr != nil {
				resp.Diagnostics.AddError(
					"Unable to create CF Client due to missing values",
					"Either user/password or client_id/client_secret must be set with api_url or CF config must exist in path (default ~/.cf/config.json)",
				)
			}
		}
	}
}

func getAndSetProviderValues(config *CloudFoundryProviderModel, resp *provider.ConfigureResponse) *managers.CloudFoundryProviderConfig {
	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	endpoint := os.Getenv("CF_API_URL")
	user := os.Getenv("CF_USER")
	password := os.Getenv("CF_PASSWORD")
	// Attribute name is cf_client_id (for backward compatability) and convention is to add `CF_` prefix to all environment variables, hence `CF_CF_CLIENT_ID`
	cfclientid := os.Getenv("CF_CF_CLIENT_ID")
	cfclientsecret := os.Getenv("CF_CF_CLIENT_SECRET")

	var skipsslvalidation bool
	var err error
	if os.Getenv("CF_SKIP_SSL_VALIDATION") != "" {
		skipsslvalidation, err = strconv.ParseBool(os.Getenv("CF_SKIP_SSL_VALIDATION"))
		if err != nil {
			addTypeCastAttributeError(resp, "Boolean", "skip_ssl_validation", "Skip SSL Validation", "CF_SKIP_SSL_VALIDATION")
			return nil
		}
	}
	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}
	if !config.User.IsNull() {
		user = config.User.ValueString()
	}
	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}
	if !config.CFClientID.IsNull() {
		cfclientid = config.CFClientID.ValueString()
	}
	if !config.CFClientSecret.IsNull() {
		cfclientsecret = config.CFClientSecret.ValueString()
	}
	checkConfig(resp, endpoint, user, password, cfclientid, cfclientsecret)
	if resp.Diagnostics.HasError() {
		return nil
	}
	if !config.SkipSslValidation.IsNull() {
		skipsslvalidation = config.SkipSslValidation.ValueBool()
	}

	c := managers.CloudFoundryProviderConfig{
		Endpoint:          strings.TrimSuffix(endpoint, "/"),
		User:              user,
		Password:          password,
		CFClientID:        cfclientid,
		CFClientSecret:    cfclientsecret,
		SkipSslValidation: skipsslvalidation,
	}
	return &c
}
func (p *CloudFoundryProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config CloudFoundryProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	checkConfigUnknown(&config, resp)

	if resp.Diagnostics.HasError() {
		return
	}

	cloudFoundryProviderConfig := getAndSetProviderValues(&config, resp)
	if resp.Diagnostics.HasError() {
		return
	}
	session, err := cloudFoundryProviderConfig.NewSession(p.httpClient)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create CF Client",
			"Client creation failed with error "+err.Error(),
		)
	}

	// Make the Cloud Foundry session available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = session
	resp.ResourceData = session
}

func (p *CloudFoundryProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NeworgResource,
		NewOrgQuotaResource,
		NewSpaceResource,
		NewUserResource,
	}
}

func (p *CloudFoundryProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewOrgDataSource,
		NewOrgQuotaDataSource,
		NewSpaceDataSource,
		NewUserDataSource,
	}
}

func New(version string, httpClient *http.Client) func() provider.Provider {
	return func() provider.Provider {
		return &CloudFoundryProvider{
			version:    version,
			httpClient: httpClient,
		}
	}
}
