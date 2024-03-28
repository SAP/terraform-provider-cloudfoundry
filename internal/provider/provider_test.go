package provider

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"regexp"
	"testing"
	"text/template"

	cfconfig "github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	testingResource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"

	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

type CloudFoundryProviderConfigPtr struct {
	Endpoint          *string
	User              *string
	Password          *string
	CFClientID        *string
	CFClientSecret    *string
	SkipSslValidation *bool
}

var redactedTestUser = CloudFoundryProviderConfigPtr{
	Endpoint:       strtostrptr("https://api.x.x.x.x.com"),
	User:           strtostrptr("xx"),
	Password:       strtostrptr("xxxx"),
	CFClientID:     strtostrptr("xx"),
	CFClientSecret: strtostrptr("xxxx"),
}

func hclProvider(cfConfig *CloudFoundryProviderConfigPtr) string {
	if cfConfig != nil {
		s := `
			provider "cloudfoundry" {
			{{- if .Endpoint}}
				api_url  = "{{.Endpoint}}"
			{{- end -}}
			{{if .User}}
				user = "{{.User}}"
			{{- end -}}
			{{if .Password}}
				password = "{{.Password}}"
			{{- end -}}
			{{if .CFClientID}}
				cf_client_id = "{{.CFClientID}}"
			{{- end -}}
			{{if .CFClientSecret}}
				cf_client_secret = "{{.CFClientSecret}}"
			{{- end -}}
			{{if .SkipSslValidation}}
				skip_ssl_validation = "{{.SkipSslValidation}}"
			{{- end }}
			}`
		tmpl, err := template.New("provider").Parse(s)
		if err != nil {
			panic(err)
		}
		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, cfConfig)
		if err != nil {
			panic(err)
		}
		return buf.String()
	}
	return `provider "cloudfoundry" {}`
}
func hclProviderWithDataSource(cfConfig *CloudFoundryProviderConfigPtr) string {
	s := `
	data "cloudfoundry_org" "org" {
		name = "tf-test-do-not-delete"
	}`
	return hclProvider(cfConfig) + s
}

func TestCloudFoundryProvider_Configure(t *testing.T) {
	t.Parallel()
	t.Run("error path - user login with missing user/password data", func(t *testing.T) {

		testingResource.Test(t, testingResource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(http.DefaultClient),
			Steps: []testingResource.TestStep{
				{
					Config: hclProviderWithDataSource(&CloudFoundryProviderConfigPtr{
						Endpoint: redactedTestUser.Endpoint,
						Password: redactedTestUser.Password,
					}),
					ExpectError: regexp.MustCompile(`Error: Missing field user`),
				},
				{
					Config: hclProviderWithDataSource(&CloudFoundryProviderConfigPtr{
						Endpoint: redactedTestUser.Endpoint,
						User:     redactedTestUser.User,
					}),
					ExpectError: regexp.MustCompile(`Error: Missing field password`),
				},
			},
		})
	})
	t.Run("error path - user login with missing clientid/clientsecret data", func(t *testing.T) {

		testingResource.Test(t, testingResource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(http.DefaultClient),
			Steps: []testingResource.TestStep{
				{
					Config: hclProviderWithDataSource(&CloudFoundryProviderConfigPtr{
						Endpoint:   redactedTestUser.Endpoint,
						CFClientID: redactedTestUser.CFClientID,
					}),
					ExpectError: regexp.MustCompile(`Error: Missing field cf_client_secret`),
				},
				{
					Config: hclProviderWithDataSource(&CloudFoundryProviderConfigPtr{
						Endpoint:       redactedTestUser.Endpoint,
						CFClientSecret: redactedTestUser.CFClientSecret,
					}),
					ExpectError: regexp.MustCompile(`Error: Missing field cf_client_id`),
				},
				{
					Config: hclProviderWithDataSource(&CloudFoundryProviderConfigPtr{
						Endpoint: redactedTestUser.Endpoint,
					}),
					ExpectError: regexp.MustCompile(`Error: Unable to create CF Client due to missing values`),
				},
			},
		})
	})
	t.Run("user login with valid user/pass data", func(t *testing.T) {
		endpoint := strtostrptr(os.Getenv("TEST_CF_API_URL"))
		user := strtostrptr(os.Getenv("TEST_CF_USER"))
		password := strtostrptr(os.Getenv("TEST_CF_PASSWORD"))
		if *endpoint == "" || *user == "" || *password == "" {
			t.Logf("\nATTENTION: Using redacted user credentions since endpoint, username & password not set as env \n Make sure you are not triggering a recording else test will fail")
			endpoint = redactedTestUser.Endpoint
			user = redactedTestUser.User
			password = redactedTestUser.Password
		}
		cfg := CloudFoundryProviderConfigPtr{
			Endpoint: endpoint,
			User:     user,
			Password: password,
		}
		recUserPass := cfg.SetupVCR(t, "fixtures/provider.user_pwd")
		defer stopQuietly(recUserPass)

		testingResource.Test(t, testingResource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(recUserPass.GetDefaultClient()),
			Steps: []testingResource.TestStep{
				{
					Config: hclProvider(&cfg) + `
					data "cloudfoundry_org" "org" {
						name = "PerformanceTeamBLR"
					}`,
				},
			},
		})
	})
	t.Run("user login with valid home directory", func(t *testing.T) {
		cfg := getCFHomeConf()
		recHomeDir := cfg.SetupVCR(t, "fixtures/provider.home_dir")
		defer stopQuietly(recHomeDir)

		testingResource.Test(t, testingResource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(recHomeDir.GetDefaultClient()),
			Steps: []testingResource.TestStep{
				{
					Config: hclProviderWithDataSource(nil),
				},
			},
		})
	})
}

func getProviders(httpClient *http.Client) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"cloudfoundry": providerserver.NewProtocol6WithError(New("test", httpClient)()),
	}
}
func getCFHomeConf() *CloudFoundryProviderConfigPtr {
	cfConf, err := cfconfig.NewFromCFHome()
	if err != nil {
		return &CloudFoundryProviderConfigPtr{
			Endpoint: strtostrptr("https://api.x.x.x.x.com"),
		}
	}
	apiEndpointURL := cfConf.ApiURL("")
	cfg := CloudFoundryProviderConfigPtr{
		Endpoint: &apiEndpointURL,
	}
	return &cfg
}
func stopQuietly(rec *recorder.Recorder) {
	if err := rec.Stop(); err != nil {
		panic(err)
	}
}
func TestCloudFoundryProvider_HasResources(t *testing.T) {
	expectedResources := []string{
		"cloudfoundry_org",
		"cloudfoundry_org_quota",
		"cloudfoundry_space",
		"cloudfoundry_user",
		"cloudfoundry_space_quota",
		"cloudfoundry_role",
		"cloudfoundry_security_group",
		"cloudfoundry_service_instance",
		"cloudfoundry_route",
		"cloudfoundry_domain",
		"cloudfoundry_app",
	}

	ctx := context.Background()
	registeredResources := []string{}

	for _, resourceFunc := range New("test", &http.Client{})().Resources(ctx) {
		var resp resource.MetadataResponse

		resourceFunc().Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "cloudfoundry"}, &resp)

		registeredResources = append(registeredResources, resp.TypeName)
	}

	assert.ElementsMatch(t, expectedResources, registeredResources)
}

func TestProvider_HasDataSources(t *testing.T) {
	expectedDataSources := []string{
		"cloudfoundry_org",
		"cloudfoundry_space",
		"cloudfoundry_org_quota",
		"cloudfoundry_user",
		"cloudfoundry_space_quota",
		"cloudfoundry_role",
		"cloudfoundry_users",
		"cloudfoundry_security_group",
		"cloudfoundry_service_instance",
		"cloudfoundry_service",
		"cloudfoundry_route",
		"cloudfoundry_domain",
		"cloudfoundry_app",
	}

	ctx := context.Background()
	registeredDataSources := []string{}

	for _, resourceFunc := range New("test", &http.Client{})().DataSources(ctx) {
		var resp datasource.MetadataResponse

		resourceFunc().Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "cloudfoundry"}, &resp)

		registeredDataSources = append(registeredDataSources, resp.TypeName)
	}

	assert.ElementsMatch(t, expectedDataSources, registeredDataSources)
}
