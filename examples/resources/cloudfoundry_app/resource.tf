terraform {
  required_providers {
    cloudfoundry = {
      source = "sap/cloudfoundry"
    }
    zipper = {
      source = "ArthurHlt/zipper"
    }
  }
}
provider "cloudfoundry" {}

provider "zipper" {}

resource "zipper_file" "fixture" {
  source      = "https://github.com/cloudfoundry-samples/cf-sample-app-nodejs.git"
  output_path = "cf-sample-app-nodejs.zip"
}

resource "cloudfoundry_app" "gobis-server" {
  name             = "tf-test-do-not-delete-nodejs"
  space            = "tf-space-1"
  org              = "PerformanceTeamBLR"
  path             = zipper_file.fixture.output_path
  source_code_hash = zipper_file.fixture.output_sha
  instances        = 1
  environment = {
    MY_ENV = "red",
  }
  strategy = "rolling"
  service_bindings = [
    {
      service_instance: "xsuaa-tf"
      params = {
        role = "Viewer"
      }
    }
  ]
  routes = [
    {
      route = "tf-test-do-not-delete-nodejs.cfapps.sap.hana.ondemand.com"
    }
  ]
}

resource "cloudfoundry_app" "http-bin-server" {
  name         = "tf-test-do-not-delete-http-bin"
  space        = "tf-space-1"
  org          = "PerformanceTeamBLR"
  docker_image = "kennethreitz/httpbin"
  strategy     = "blue-green"
  labels = {
    "app" = "backend",
    "env" = "production"
  }
  annotations = {
    "created-by" = "me",
  }
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

resource "cloudfoundry_app" "http-bin-sidecar" {
  name         = "tf-test-do-not-delete-http-bin-sidecar"
  space        = "tf-space-1"
  org          = "PerformanceTeamBLR"
  docker_image = "kennethreitz/httpbin"
  sidecars = [
    {
      name = "sidecar-1"
      process_types = [
        "worker"
      ]
      command = "sleep 5200"
      memory  = "256M"
    }
  ]
  no_route = true
}