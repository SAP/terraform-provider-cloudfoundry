terraform {
  required_providers {
    cloudfoundry = {
      source = "sap/cloudfoundry"
    }
  }
}
provider "cloudfoundry" {}

resource "cloudfoundry_app" "app" {
  name      = "tf-test" 
}

resource "cloudfoundry_app" "gobis-server" {
  name = "go-server-lite"
  space = data.cloudfoundry_space.space.id
  path = "./good-app.zip"
  buildpack = "https://github.com/cloudfoundry/go-buildpack.git"
  health_check_type = "http"
  instances = 3
  health_check_http_endpoint = "/health"
  environment = {
    "MY_ENV" = "red",
    "GOPACKAGENAME" = "github.com/my-cf-sample",
    "GOVERSION" = "go1.22.0"
  }
  strategy = "rolling"
  routes {
    route = cloudfoundry_route.default.id
  }
}