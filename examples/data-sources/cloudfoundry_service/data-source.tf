terraform {
  required_providers {
    cloudfoundry = {
      source = "sap/cloudfoundry"
    }
  }
}
provider "cloudfoundry" {
}

data "cloudfoundry_service" "xsuaa-offering" {
    name = "xsuaa"
}

output "serviceplans" {
    value = data.cloudfoundry_service.xsuaa-offering.service_plans
}