package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAppDataSource_Configure(t *testing.T) {
	t.Parallel()
	t.Run("happy path - read app with bits", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/datasource_app_bits")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + `
data "cloudfoundry_app" "app" {
	name = "tf-test-do-not-delete-http-bin"
	space = "tf-space-1"
	org = "PerformanceTeamBLR"
}`,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.cloudfoundry_app.app", "name", "tf-test-do-not-delete-http-bin"),
						resource.TestCheckResourceAttr("data.cloudfoundry_app.app", "space", "tf-space-1"),
						resource.TestCheckResourceAttr("data.cloudfoundry_app.app", "org", "PerformanceTeamBLR"),
						resource.TestCheckResourceAttr("data.cloudfoundry_app.app", "instances", "1"),
						resource.TestCheckResourceAttr("data.cloudfoundry_app.app", "memory", "256M"),
						resource.TestCheckResourceAttr("data.cloudfoundry_app.app", "disk_quota", "1024M"),
						resource.TestCheckResourceAttr("data.cloudfoundry_app.app", "health_check_type", "http"),
						resource.TestCheckResourceAttr("data.cloudfoundry_app.app", "health_check_http_endpoint", "/get"),
						resource.TestCheckResourceAttr("data.cloudfoundry_app.app", "readiness_health_check_type", "http"),
						resource.TestCheckResourceAttr("data.cloudfoundry_app.app", "readiness_health_check_http_endpoint", "/get"),
						resource.TestCheckResourceAttr("data.cloudfoundry_app.app", "type", "web"),
						resource.TestCheckResourceAttr("data.cloudfoundry_app.app", "docker_image", ""),
						resource.TestCheckResourceAttr("data.cloudfoundry_app.app", "strategy", ""),
						resource.TestCheckResourceAttr("data.cloudfoundry_app.app", "no_route", "false"),
						resource.TestCheckResourceAttr("data.cloudfoundry_app.app", "buildpacks.#", "1"),
						resource.TestCheckResourceAttr("data.cloudfoundry_app.app", "buildpacks.0", "go_buildpack"),
						resource.TestCheckResourceAttr("data.cloudfoundry_app.app", "services.#", "1"),
						resource.TestCheckResourceAttr("data.cloudfoundry_app.app", "services.0", "test-service"),
						resource.TestCheckResourceAttr("data.cloudfoundry_app.app", "routes.#", "1"),
						resource.TestCheckResourceAttr("data.cloudfoundry_app.app", "routes.0.route", "http-bin.cfapps.sap.hana.ondemand.com"),
						resource.TestCheckResourceAttr("data.cloudfoundry_app.app", "routes.0.protocol", "http1"),
					),
				},
			},
		})
	})
	t.Run("error path - get unavailable datasource app", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/datasource_app_invalid_app_name")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + `
data "cloudfoundry_app" "app" {
	name = "testunavailableapp"
	space = "tf-space-1"
	org = "PerformanceTeamBLR"
}`,
					ExpectError: regexp.MustCompile(`Error: Error finding given app`),
				},
			},
		})
	})
}
