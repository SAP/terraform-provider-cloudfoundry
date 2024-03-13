terraform {
  required_providers {
    cloudfoundry = {
      source = "sap/cloudfoundry"
    }
  }
}
provider "cloudfoundry" {}

data "cloudfoundry_org" "team_org" {
  name = "PerformanceTeamBLR"
}

data "cloudfoundry_space" "team_space" {
  name = "PerformanceTeamBLR"
  org = data.cloudfoundry_org.team_org.id
}

data "cloudfoundry_service" "xsuaa_svc" {
  name = "xsuaa"
}
data "cloudfoundry_service" "autoscaler_svc" {
  name = "autoscaler"
}
resource "cloudfoundry_service_instance" "xsuaa_svc" {
  name         = "xsuaa_svc"
  type         = "managed"
  tags =       ["terraform-test","test1"]
  space        = data.cloudfoundry_space.team_space.id
  service_plan = data.cloudfoundry_service.xsuaa_svc.service_plans["application"]
  parameters = <<EOT
  {
  "xsappname": "tf-test23",
  "tenant-mode": "dedicated",
  "description": "tf test123",
  "foreign-scope-references": ["user_attributes"],
  "scopes": [
    {
      "name": "uaa.user",
      "description": "UAA"
    }
  ],
  "role-templates": [
    {
      "name": "Token_Exchange",
      "description": "UAA",
      "scope-references": ["uaa.user"]
    }
  ]
}
EOT
}

# Managed service instance without parameters
resource "cloudfoundry_service_instance" "dev-autoscaler" {
  name = "tf-autoscaler-test"
  type = "managed"
  tags = ["terraform-test","autoscaler"]
  space = data.cloudfoundry_space.team_space.id
  service_plan = data.cloudfoundry_service.autoscaler_svc.service_plans["standard"]
}
# User provided service instance
resource "cloudfoundry_service_instance" "dev-usp" {
  name = "tf-usp-test"
  type = "user-provided"
  tags = ["terraform-test","usp"]
  space = data.cloudfoundry_space.team_space.id
  credentials = <<EOT
  {
    "user": "user1",
    "password": "demo122"
  }
  EOT
}
