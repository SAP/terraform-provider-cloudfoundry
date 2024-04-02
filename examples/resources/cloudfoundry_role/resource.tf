resource "cloudfoundry_role" "my_role" {
  username = "debaditya.ray@sap.com"
  type     = "organization_user"
  org      = "784b4cd0-4771-4e4d-9052-a07e178bae56"
}