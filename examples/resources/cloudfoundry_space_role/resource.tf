resource "cloudfoundry_space_role" "my_role" {
  username = "debaditya.ray@sap.com"
  type     = "space_manager"
  space    = "dd457c79-f7c9-4828-862b-35843d3b646d"
}