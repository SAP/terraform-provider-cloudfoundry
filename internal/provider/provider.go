package provider

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/SAP/terraform-provider-cloudfoundry/internal/provider/managers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &CloudFoundryProvider{}

type CloudFoundryProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type CloudFoundryProviderModel struct {
	Endpoint                  types.String `tfsdk:"api_url"`
	User                      types.String `tfsdk:"user"`
	Password                  types.String `tfsdk:"password"`
	SSOPasscode               types.String `tfsdk:"sso_passcode"`
	CFClientID                types.String `tfsdk:"cf_client_id"`
	CFClientSecret            types.String `tfsdk:"cf_client_secret"`
	UaaClientID               types.String `tfsdk:"uaa_client_id"`
	UaaClientSecret           types.String `tfsdk:"uaa_client_secret"`
	SkipSslValidation         types.Bool   `tfsdk:"skip_ssl_validation"`
	AppLogsMax                types.Int64  `tfsdk:"app_logs_max"`
	DefaultQuotaName          types.String `tfsdk:"default_quota_name"`
	PurgeWhenDelete           types.Bool   `tfsdk:"purge_when_delete"`
	StoreTokensPath           types.String `tfsdk:"store_tokens_path"`
	ForceNotFailBrokerCatalog types.Bool   `tfsdk:"force_broker_not_fail_when_catalog_not_accessible"`
}

func (p *CloudFoundryProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cf"
	resp.Version = p.version
}

func (p *CloudFoundryProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Cloud Foundry Terraform plugin is an integration that allows users to leverage Terraform, an infrastructure as code tool, to define and provision infrastructure resources within the Cloud Foundry platform.",
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				MarkdownDescription: "Specific URL representing the entry point for communication between the client and a Cloud Foundry instance.",
				Optional:            false,
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "A unique identifier associated with an individual or entity for authentication & authorization purposes.",
				Optional:            true,
				Sensitive:           true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "A confidential alphanumeric code associated with a user account on the Cloud Foundry platform",
				Optional:            true,
				Sensitive:           true,
			},
			"sso_passcode": schema.StringAttribute{
				MarkdownDescription: "A temporary or one-time code used as part of the authentication process.",
				Optional:            true,
				Sensitive:           true,
			},
			"cf_client_id": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"cf_client_secret": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"uaa_client_id": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"uaa_client_secret": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"skip_ssl_validation": schema.BoolAttribute{
				Optional: true,
			},
			"default_quota_name": schema.StringAttribute{
				MarkdownDescription: "Name of the default quota",
				Optional:            true,
			},
			"app_logs_max": schema.Int64Attribute{
				MarkdownDescription: "Number of logs message which can be see when app creation is errored (-1 means all messages stored)",
				Optional:            true,
			},
			"purge_when_delete": schema.BoolAttribute{
				MarkdownDescription: "Set to true to purge when deleting a resource (e.g.: service instance, service broker)",
				Optional:            true,
			},
			"store_tokens_path": schema.StringAttribute{
				MarkdownDescription: "Path to a file to store tokens used for login. (this is useful for sso, this avoid requiring each time sso passcode)",
				Optional:            true,
			},
			"force_broker_not_fail_when_catalog_not_accessible": schema.BoolAttribute{
				MarkdownDescription: "Set to true to not trigger fail on catalog on service broker",
				Optional:            true,
			},
		},
	}
}
func addGenericAttributeError(resp *provider.ConfigureResponse, status string, pathRoot string, commonName string, envName string) {
	resp.Diagnostics.AddAttributeError(
		path.Root(pathRoot),
		`{{status}} field {{pathRoot}}`,
		`The provider cannot create the Cloud Foundry API client as there is an unknown configuration value for the Cloud Foundry {{commonName}}.`+
			`Either target apply the source of the value first, set the value statically in the configuration, or use the {{envName}} environment variable, ensure value is not empty.`,
	)
}
func addTypeCastAttributeError(resp *provider.ConfigureResponse, expectedType string, pathRoot string, commonName string, envName string) {
	resp.Diagnostics.AddAttributeError(
		path.Root(pathRoot),
		`Expected {{expectedType}} in field {{pathRoot}}`,
		`The provider cannot create the Cloud Foundry API client as there is an invalid configuration value for the Cloud Foundry {{commonName}}.`+
			`Ensure {{envName}} is of type {{expectedType}}`,
	)
}
func checkConfigUnknown(config *CloudFoundryProviderModel, resp *provider.ConfigureResponse) {
	if config.Endpoint.IsUnknown() {
		addGenericAttributeError(resp, "Unknown", "api_url", "API Endpoint", "CF_API_URL")
	}
	if config.User.IsUnknown() && !config.Password.IsUnknown() {
		addGenericAttributeError(resp, "Unknown", "user", "Username", "CF_USER")
	} else if !config.User.IsUnknown() && config.Password.IsUnknown() {
		addGenericAttributeError(resp, "Unknown", "password", "Password", "CF_PASSWORD")
	} else if config.CFClientID.IsUnknown() && !config.CFClientSecret.IsUnknown() {
		addGenericAttributeError(resp, "Unknown", "cf_client_id", "CF Client ID", "CF_CF_CLIENT_ID")
	} else if !config.CFClientID.IsUnknown() && config.CFClientSecret.IsUnknown() {
		addGenericAttributeError(resp, "Unknown", "cf_client_secret", "CF Client Secret", "CF_CF_CLIENT_SECRET")
	} else if config.UaaClientID.IsUnknown() && !config.UaaClientSecret.IsUnknown() {
		addGenericAttributeError(resp, "Unknown", "uaa_client_id", "UAA Client ID", "CF_UAA_CLIENT_ID")
	} else if !config.UaaClientID.IsUnknown() && config.UaaClientSecret.IsUnknown() {
		addGenericAttributeError(resp, "Unknown", "uaa_client_secret", "UAA Client Secret", "CF_UAA_CLIENT_SECRET")
	} else if config.User.IsUnknown() && config.Password.IsUnknown() && config.UaaClientID.IsUnknown() && config.UaaClientSecret.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unable to create CF Client due to unknown values",
			"Couple of user/password or uaa_client_id/uaa_client_secret must be set",
		)
	}
}

func checkConfigNull(resp *provider.ConfigureResponse, endpoint string, user string, password string, cfclientid string, cfclientsecret string, uaaclientid string, uaaclientsecret string, storetokenpath string) {
	if endpoint == "" {
		addGenericAttributeError(resp, "Missing", "api_url", "API Endpoint", "CF_API_URL")
	}
	switch {
	case user == "" && password != "":
		addGenericAttributeError(resp, "Missing", "user", "Username", "CF_USER")
	case user != "" && password == "":
		addGenericAttributeError(resp, "Missing", "password", "Password", "CF_PASSWORD")
	case cfclientid == "" && cfclientsecret != "":
		addGenericAttributeError(resp, "Missing", "cf_client_id", "Client ID", "CF_CF_CLIENT_ID")
	case cfclientid != "" && cfclientsecret == "":
		addGenericAttributeError(resp, "Missing", "cf_client_secret", " Client Secret", "CF_CF_CLIENT_SECRET")
	case uaaclientid == "" && uaaclientsecret != "":
		addGenericAttributeError(resp, "Missing", "uaa_client_id", "UAA Client ID", "CF_UAA_CLIENT_ID")
	case uaaclientid != "" && uaaclientsecret == "":
		addGenericAttributeError(resp, "Missing", "uaa_client_secret", "UAA Client Secret", "CF_UAA_CLIENT_SECRET")
	case user == "" && password == "" && uaaclientid == "" && uaaclientsecret == "" && storetokenpath == "":
		resp.Diagnostics.AddError(
			"Unable to create CF Client due to missing values",
			"Couple of user/password or uaa_client_id/uaa_client_secret or store token path must be set",
		)
	}
}

func getAndSetProviderValues(config *CloudFoundryProviderModel, resp *provider.ConfigureResponse) *managers.CloudFoundryProviderConfig {
	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	endpoint := os.Getenv("CF_API_URL")
	user := os.Getenv("CF_USER")
	password := os.Getenv("CF_PASSWORD")
	ssopasscode := os.Getenv("CF_SSO_PASSCODE")
	cfclientid := os.Getenv("CF_CF_CLIENT_ID")
	cfclientsecret := os.Getenv("CF_CF_CLIENT_SECRET")
	uaaclientid := os.Getenv("CF_UAA_CLIENT_ID")
	uaaclientsecret := os.Getenv("CF_UAA_CLIENT_SECRET")

	var skipsslvalidation, purgewhendelete, forcenotfailbrokercatalog bool
	var defaultquotaname, storetokenpath string
	var applogsmax int64
	var err error
	if os.Getenv("CF_SKIP_SSL_VALIDATION") != "" {
		skipsslvalidation, err = strconv.ParseBool(os.Getenv("CF_SKIP_SSL_VALIDATION"))
		if err != nil {
			addTypeCastAttributeError(resp, "Boolean", "skip_ssl_validation", "Skip SSL Validation", "CF_SKIP_SSL_VALIDATION")
			return nil
		}
	}
	defaultquotaname = os.Getenv("CF_DEFAULT_QUOTA_NAME")
	if os.Getenv("CF_APP_LOGS_MAX") != "" {
		applogsmax, err = strconv.ParseInt(os.Getenv("CF_APP_LOGS_MAX"), 10, 64)
		if err != nil {
			addTypeCastAttributeError(resp, "Int64", "app_logs_max", "Maximum App Logs", "CF_APP_LOGS_MAX")
			return nil
		}
	}
	if os.Getenv("CF_PURGE_WHEN_DELETE") != "" {
		purgewhendelete, err = strconv.ParseBool(os.Getenv("CF_PURGE_WHEN_DELETE"))
		if err != nil {
			addTypeCastAttributeError(resp, "Boolean", "purge_when_delete", "Purge On Deletion", "CF_PURGE_WHEN_DELETE")
			return nil
		}
	}
	storetokenpath = os.Getenv("CF_STORE_TOKEN_PATH")
	if os.Getenv("CF_FORCE_BROKER_NOT_FAIL_WHEN_CATALOG_NOT_ACCESSIBLE") != "" {
		forcenotfailbrokercatalog, err = strconv.ParseBool(os.Getenv("CF_FORCE_BROKER_NOT_FAIL_WHEN_CATALOG_NOT_ACCESSIBLE"))
		if err != nil {
			addTypeCastAttributeError(resp, "Boolean", "force_broker_not_fail_when_catalog_not_accessible", "Force Broker Not Fail When Catalog Not Accessible", "CF_FORCE_BROKER_NOT_FAIL_WHEN_CATALOG_NOT_ACCESSIBLE")
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
	if !config.SSOPasscode.IsNull() {
		ssopasscode = config.SSOPasscode.ValueString()
	}
	if !config.UaaClientID.IsNull() {
		uaaclientid = config.UaaClientID.ValueString()
	}
	if !config.UaaClientSecret.IsNull() {
		uaaclientsecret = config.UaaClientSecret.ValueString()
	}
	checkConfigNull(resp, endpoint, user, password, cfclientid, cfclientsecret, uaaclientid, uaaclientsecret, storetokenpath)
	if resp.Diagnostics.HasError() {
		return nil
	}
	if !config.SkipSslValidation.IsNull() {
		skipsslvalidation = config.SkipSslValidation.ValueBool()
	}
	if !config.DefaultQuotaName.IsNull() {
		defaultquotaname = config.DefaultQuotaName.ValueString()
	}
	if !config.AppLogsMax.IsNull() {
		applogsmax = config.AppLogsMax.ValueInt64()
	}
	if !config.PurgeWhenDelete.IsNull() {
		purgewhendelete = config.PurgeWhenDelete.ValueBool()
	}
	if !config.StoreTokensPath.IsNull() {
		storetokenpath = config.StoreTokensPath.ValueString()
	}
	if !config.ForceNotFailBrokerCatalog.IsNull() {
		forcenotfailbrokercatalog = config.ForceNotFailBrokerCatalog.ValueBool()
	}

	if defaultquotaname == "" {
		defaultquotaname = "default"
	}
	if applogsmax == 0 {
		applogsmax = 30
	}
	c := managers.CloudFoundryProviderConfig{
		Endpoint:                  strings.TrimSuffix(endpoint, "/"),
		User:                      user,
		Password:                  password,
		SSOPasscode:               ssopasscode,
		CFClientID:                cfclientid,
		CFClientSecret:            cfclientsecret,
		UaaClientID:               uaaclientid,
		UaaClientSecret:           uaaclientsecret,
		SkipSslValidation:         skipsslvalidation,
		AppLogsMax:                applogsmax,
		DefaultQuotaName:          defaultquotaname,
		PurgeWhenDelete:           purgewhendelete,
		StoreTokensPath:           storetokenpath,
		ForceNotFailBrokerCatalog: forcenotfailbrokercatalog,
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
	cloudFoundryProviderConfig.NewSession()

	// // Make the Cloud Foundry client available during DataSource and Resource
	// // type Configure methods.
	// resp.DataSourceData = client
	// resp.ResourceData = client

}

func (p *CloudFoundryProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *CloudFoundryProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CloudFoundryProvider{
			version: version,
		}
	}
}
