terraform {
  required_providers {
    cloudfoundry = {
      source  = "SAP/cloudfoundry"
      version = "0.2.0-beta"
    }
  }
}
provider "cloudfoundry" {
  api_url  = "https://api.cf.example.com"
  username = "admin"
  password = "admin"
}