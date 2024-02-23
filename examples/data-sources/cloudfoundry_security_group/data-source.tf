terraform {
  required_providers {
    cloudfoundry = {
      source = "sap/cloudfoundry"
    }
  }
}
provider "cloudfoundry" {}

data "cloudfoundry_security_group" "sgroup" {
  name = "riemann"
}

output "sgroup" {
  value = data.cloudfoundry_security_group.sgroup
}