terraform {
  required_providers {
    cloudfoundry = {
      source  = "SAP/cloudfoundry"
      version = "1.0.0-rc1"
    }
  }
}
provider "cloudfoundry" {
  api_url  = "https://api.cf.example.com"
  user     = "admin"
  password = "admin"
}