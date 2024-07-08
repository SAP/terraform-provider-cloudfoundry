terraform {
  required_providers {
    cloudfoundry = {
      source  = "SAP/cloudfoundry"
    }
  }
}
provider "cloudfoundry" {
}

resource "cloudfoundry_user" "my_user" {
  username = "hi"
  email = "hi@gmail.com"
  given_name = "Das"
  family_name = "Hi"
  annotations = { purpose : "testing" }
}
