package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAppResource_Configure(t *testing.T) {
	t.Parallel()
	t.Run("happy path - create app with bits", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_app_bits")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + `
resource "cloudfoundry_app" "app" {
	name                                 = "cf-nodejs"
  space                                = "tf-space-1" 
  org                                  = "PerformanceTeamBLR"
  path                                 = "../../asset/cf-sample-app-nodejs.zip"
	memory                               = "256M"
	disk_quota                           = "1024M"
	health_check_type                    = "http"
	health_check_http_endpoint           = "/"
	readiness_health_check_type          = "http"
	readiness_health_check_http_endpoint = "/"
  instances                            = 2
  environment = {
    MY_ENV = "red",
  }
  strategy = "rolling"
  routes = [
    {
      route = "cf-sample-test.cfapps.sap.hana.ondemand.com"
    }
  ]
}
					`,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "name", "cf-nodejs"),
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "space", "tf-space-1"),
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "org", "PerformanceTeamBLR"),
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "instances", "2"),
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "memory", "256M"),
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "disk_quota", "1024M"),
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "health_check_type", "http"),
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "health_check_http_endpoint", "/"),
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "strategy", "rolling"),
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "environment.MY_ENV", "red"),
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "routes.0.route", "cf-sample-test.cfapps.sap.hana.ondemand.com"),
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "routes.0.protocol", "http1"),
					),
				},
			},
		})
	})
	t.Run("happy path - create app with docker and process", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_app_docker")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + `
resource "cloudfoundry_app" "app" {
	name         = "http-bin"
	space        = "tf-space-1"
	org          = "PerformanceTeamBLR"
	docker_image = "kennethreitz/httpbin"
	strategy		 = "blue-green"
	processes = [
		{
			type                                 = "web",
			instances                            = 1
			memory                               = "256M"
			disk_quota                           = "1024M"
			health_check_type                    = "http"
			health_check_http_endpoint           = "/get"
			readiness_health_check_type          = "http"
			readiness_health_check_http_endpoint = "/get"
		}
	]
	no_route = true
}
					`,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "docker_image", "kennethreitz/httpbin"),
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "strategy", "blue-green"),
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "no_route", "true"),
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "processes.0.instances", "1"),
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "processes.0.memory", "256M"),
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "processes.0.disk_quota", "1024M"),
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "processes.0.health_check_type", "http"),
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "processes.0.health_check_http_endpoint", "/get"),
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "processes.0.readiness_health_check_type", "http"),
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "processes.0.readiness_health_check_http_endpoint", "/get"),
						resource.TestCheckResourceAttr("cloudfoundry_app.app", "processes.0.type", "web"),
					),
				},
			},
		})
	})
	t.Run("happy path - create app with sidecar", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_app_sidecar")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + `
resource "cloudfoundry_app" "http-bin-sidecar" {
	name         = "http-bin-sidecar"
	space        = "tf-space-1"
	org          = "PerformanceTeamBLR"
	docker_image = "kennethreitz/httpbin"
	sidecars = [
		{
			name         = "sidecar-1"
			process_types = [
				"worker"
			]
			command      = "sleep 5200"
			memory       = "256M"
		}
	]
	no_route = true
}
					`,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("cloudfoundry_app.http-bin-sidecar", "sidecars.0.name", "sidecar-1"),
						resource.TestCheckResourceAttr("cloudfoundry_app.http-bin-sidecar", "sidecars.0.command", "sleep 5200"),
						resource.TestCheckResourceAttr("cloudfoundry_app.http-bin-sidecar", "sidecars.0.process_types.#", "1"),
					),
				},
			},
		})
	})
}
