terraform {
  required_providers {
    cloudfoundry = {
      source = "sap/cloudfoundry"
    }
  }
}
provider "cloudfoundry" {
}

data "cloudfoundry_service_instance" "svc" {
    name = "tf-test-do-not-delete"
    space = "02c0cc92-6ecc-44b1-b7b2-096ca19ee143"
}

output "guid" {
    value = data.cloudfoundry_service_instance.svc.id
}
