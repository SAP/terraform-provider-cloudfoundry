terraform {
  required_providers {
    cloudfoundry = {
      source = "sap/cloudfoundry"
    }
  }
}
provider "cloudfoundry" {}


resource "cloudfoundry_space" "space" {
  name = "test"
  org_name = "PerformanceTeamBLR"
  labels = {"name":"test","type" : "test"}
  quota = "79878f64-e1d1-4a71-9b69-7db01ac0119f"
  allow_ssh = true
  asgs = ["2e298b15-abeb-43ff-898d-f0e2c05e8858","56eedab7-cb97-469b-a3e9-89521827c039"]
}

output "out" {
  value = resource.cloudfoundry_space.space
}