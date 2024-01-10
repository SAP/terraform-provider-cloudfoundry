package provider

import (
	"bytes"
	"net/http"
	"os"
	"regexp"
	"testing"
	"text/template"

	cfconfig "github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	testingResource "github.com/hashicorp/terraform-plugin-testing/helper/resource"

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

func hclProvider(cfConfig CloudFoundryProviderConfigPtr) string {
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
	}
	data "cloudfoundry_org" "org" {
		name = "PerformanceTeamBLR"
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
func strtostrptr(s string) *string {
	return &s
}

func TestCloudFoundryProvider_Configure(t *testing.T) {
	t.Parallel()
	t.Run("user login with missing user/password data", func(t *testing.T) {

		testingResource.Test(t, testingResource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(http.DefaultClient),
			Steps: []testingResource.TestStep{
				{
					Config: hclProvider(CloudFoundryProviderConfigPtr{
						Endpoint: redactedTestUser.Endpoint,
						Password: redactedTestUser.Password,
					}),
					ExpectError: regexp.MustCompile(`Error: Missing field user`),
				},
				{
					Config: hclProvider(CloudFoundryProviderConfigPtr{
						Endpoint: redactedTestUser.Endpoint,
						User:     redactedTestUser.User,
					}),
					ExpectError: regexp.MustCompile(`Error: Missing field password`),
				},
			},
		})
	})
	t.Run("user login with missing clientid/clientsecret data", func(t *testing.T) {

		testingResource.Test(t, testingResource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(http.DefaultClient),
			Steps: []testingResource.TestStep{
				{
					Config: hclProvider(CloudFoundryProviderConfigPtr{
						Endpoint:   redactedTestUser.Endpoint,
						CFClientID: redactedTestUser.CFClientID,
					}),
					ExpectError: regexp.MustCompile(`Error: Missing field cf_client_secret`),
				},
				{
					Config: hclProvider(CloudFoundryProviderConfigPtr{
						Endpoint:       redactedTestUser.Endpoint,
						CFClientSecret: redactedTestUser.CFClientSecret,
					}),
					ExpectError: regexp.MustCompile(`Error: Missing field cf_client_id`),
				},
				{
					Config: hclProvider(CloudFoundryProviderConfigPtr{
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
			t.Logf("\nATTENTION: Using redacted user credentions since endpoint, username & password not set as env")
			t.Logf("Make sure you are not triggering a recording else test will fail")
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
					Config: hclProvider(cfg),
				},
			},
		})
	})
	t.Run("user login with valid home directory", func(t *testing.T) {
		cfConf, err := cfconfig.NewFromCFHome()
		if err != nil {
			panic(err)
		}
		cfg := CloudFoundryProviderConfigPtr{
			Endpoint: &cfConf.APIEndpointURL,
		}
		recHomeDir := cfg.SetupVCR(t, "fixtures/provider.home_dir")
		defer stopQuietly(recHomeDir)

		testingResource.Test(t, testingResource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(recHomeDir.GetDefaultClient()),
			Steps: []testingResource.TestStep{
				{
					Config: hclProvider(CloudFoundryProviderConfigPtr{}),
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
func stopQuietly(rec *recorder.Recorder) {
	if err := rec.Stop(); err != nil {
		panic(err)
	}
}
