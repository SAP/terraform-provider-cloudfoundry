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
}

output "guid" {
    value = data.cloudfoundry_service_instance.svc.id
}
