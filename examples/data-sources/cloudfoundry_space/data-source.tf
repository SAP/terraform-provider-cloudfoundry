terraform {
  required_providers {
    cloudfoundry = {
      source = "sap/cloudfoundry"
    }
  }
}
provider "cloudfoundry" {}


data "cloudfoundry_space" "space" {
  name = "PerformanceTeamBLR"
  org = "784b4cd0-4771-4e4d-9052-a07e178bae56"
}

output "out" {
  value = data.cloudfoundry_space.space
}
