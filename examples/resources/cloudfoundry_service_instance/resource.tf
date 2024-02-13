terraform {
  required_providers {
    cloudfoundry = {
      source = "sap/cloudfoundry"
    }
  }
}
provider "cloudfoundry" {}

locals {
    params = jsondecode(file("/Users/i342464/myworks/terraform/security.json"))
}

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
resource "cloudfoundry_service_instance" "xsuaa_svc" {
  name         = "tf-xsuaa-com1"
  type         = "managed"
  tags =       ["team","perf"]
  space        = data.cloudfoundry_space.team_space.id
  service_plan = data.cloudfoundry_service.xsuaa_svc.service_plans["application"]
  parameters = <<EOT
  {
  "xsappname": "tf-test1",
  "tenant-mode": "dedicated",
  "description": "tf test1",
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