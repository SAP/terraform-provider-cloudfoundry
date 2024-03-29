package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAppResource_Configure(t *testing.T) {
	t.Parallel()
	resourceName := "cloudfoundry_app.app"
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
  path                                 = "../../assets/cf-sample-app-nodejs.zip"
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
						resource.TestCheckResourceAttr(resourceName, "name", "cf-nodejs"),
						resource.TestCheckResourceAttr(resourceName, "space", "tf-space-1"),
						resource.TestCheckResourceAttr(resourceName, "org", "PerformanceTeamBLR"),
						resource.TestCheckResourceAttr(resourceName, "instances", "2"),
						resource.TestCheckResourceAttr(resourceName, "memory", "256M"),
						resource.TestCheckResourceAttr(resourceName, "disk_quota", "1024M"),
						resource.TestCheckResourceAttr(resourceName, "health_check_type", "http"),
						resource.TestCheckResourceAttr(resourceName, "health_check_http_endpoint", "/"),
						resource.TestCheckResourceAttr(resourceName, "strategy", "rolling"),
						resource.TestCheckResourceAttr(resourceName, "environment.MY_ENV", "red"),
						resource.TestCheckResourceAttr(resourceName, "routes.0.route", "cf-sample-test.cfapps.sap.hana.ondemand.com"),
						resource.TestCheckResourceAttr(resourceName, "routes.0.protocol", "http1"),
					),
				},
			},
		})
	})
	t.Run("happy path - update app with bits", func(t *testing.T) {
		cfg := getCFHomeConf()
		rec := cfg.SetupVCR(t, "fixtures/resource_app_bits_update")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider(nil) + `
resource "cloudfoundry_app" "app" {
	name                                 = "cf-nodejs-update"
  space                                = "tf-space-1" 
  org                                  = "PerformanceTeamBLR"
  path                                 = "../../assets/cf-sample-app-nodejs.zip"
	source_code_hash                     = "1234567890"
	memory                               = "256M"
	disk_quota                           = "1024M"
  instances                            = 1
  environment = {
    MY_ENV = "red",
  }
	labels = {
		MY_LABEL = "red",
	}
  strategy = "blue-green"
}
					`,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", "cf-nodejs-update"),
						resource.TestCheckResourceAttr(resourceName, "space", "tf-space-1"),
						resource.TestCheckResourceAttr(resourceName, "org", "PerformanceTeamBLR"),
						resource.TestCheckResourceAttr(resourceName, "instances", "1"),
						resource.TestCheckResourceAttr(resourceName, "memory", "256M"),
						resource.TestCheckResourceAttr(resourceName, "disk_quota", "1024M"),
						resource.TestCheckResourceAttr(resourceName, "strategy", "blue-green"),
						resource.TestCheckResourceAttr(resourceName, "environment.MY_ENV", "red"),
						resource.TestCheckResourceAttr(resourceName, "labels.MY_LABEL", "red"),
					),
				},
				{
					Config: hclProvider(nil) + `
resource "cloudfoundry_app" "app" {
	name                                 = "cf-nodejs-update"
  space                                = "tf-space-1" 
  org                                  = "PerformanceTeamBLR"
  path                                 = "../../assets/cf-sample-app-nodejs.zip"
	source_code_hash                     = "999999"
	memory                               = "256M"
	disk_quota                           = "1024M"
  instances                            = 2
  labels = {
		MY_LABEL = "blue",
	}
  strategy = "blue-green"
}
					`,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", "cf-nodejs-update"),
						resource.TestCheckResourceAttr(resourceName, "space", "tf-space-1"),
						resource.TestCheckResourceAttr(resourceName, "org", "PerformanceTeamBLR"),
						resource.TestCheckResourceAttr(resourceName, "instances", "2"),
						resource.TestCheckResourceAttr(resourceName, "memory", "256M"),
						resource.TestCheckResourceAttr(resourceName, "disk_quota", "1024M"),
						resource.TestCheckResourceAttr(resourceName, "strategy", "blue-green"),
						resource.TestCheckResourceAttr(resourceName, "labels.MY_LABEL", "blue"),
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
						resource.TestCheckResourceAttr(resourceName, "docker_image", "kennethreitz/httpbin"),
						resource.TestCheckResourceAttr(resourceName, "strategy", "blue-green"),
						resource.TestCheckResourceAttr(resourceName, "no_route", "true"),
						resource.TestCheckResourceAttr(resourceName, "processes.0.instances", "1"),
						resource.TestCheckResourceAttr(resourceName, "processes.0.memory", "256M"),
						resource.TestCheckResourceAttr(resourceName, "processes.0.disk_quota", "1024M"),
						resource.TestCheckResourceAttr(resourceName, "processes.0.health_check_type", "http"),
						resource.TestCheckResourceAttr(resourceName, "processes.0.health_check_http_endpoint", "/get"),
						resource.TestCheckResourceAttr(resourceName, "processes.0.readiness_health_check_type", "http"),
						resource.TestCheckResourceAttr(resourceName, "processes.0.readiness_health_check_http_endpoint", "/get"),
						resource.TestCheckResourceAttr(resourceName, "processes.0.type", "web"),
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
