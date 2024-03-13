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
  output_path = "/Users/I524895/go/src/github.com/SAP/terraform-provider-cloudfoundry/examples/resources/cloudfoundry_app/cf-sample-app-nodejs.zip"
}

resource "cloudfoundry_app" "gobis-server" {
  name             = "cf-nodejs"
  space            = "tf-space-1"
  org              = "PerformanceTeamBLR"
  path             = zipper_file.fixture.output_path
  source_code_hash = zipper_file.fixture.output_sha
  instances        = 1
  environment = {
    MY_ENV = "red",
  }
  strategy = "rolling"
  routes = [
    {
      route = "cf-sample.cfapps.sap.hana.ondemand.com"
    }
  ]
}

# resource "cloudfoundry_app" "gobis-server-1" {
#   name = "go-server-lite-1"
#   space = "tf-space-1"
#   org = "PerformanceTeamBLR"
#   path = "/Users/I524895/go/src/github.com/SAP/terraform-provider-cloudfoundry/examples/resources/cloudfoundry_app/good-app.zip"
#   buildpacks = ["https://github.com/cloudfoundry/go-buildpack.git"]
#   environment = {
#     MY_ENV = "red",
#     GOPACKAGENAME = "github.com/my-cf-sample",
#     GOVERSION = "go1.22.1"
#   }
#   processes = [
#     {
#       type = "web",
#       instances = 2
#       memory = "256M"
#       health_check_type = "http"
#       health_check_http_endpoint = "/health"
#       readiness_health_check_type = "http"
#       readiness_health_check_http_endpoint = "/health"
#     }
#   ] 
#   routes = [
#     {
#       route = "go-server-lite-1.cfapps.sap.hana.ondemand.com"
#     }
#   ]
# }
