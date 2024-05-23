resource "cloudfoundry_org_role" "my_role" {
  username = "test123@example.com"
  type     = "organization_user"
  org      = "784b4cd0-4771-4e4d-9052-a07e178bae56"
}