resource "cloudfoundry_user_groups" "my_user_groups" {
  user   = "ad6bb1e0-05f6-4440-9485-6fe20b38c500"
  origin = "uaa"
  groups = ["cloud_controller.admin", "scim.read", "scim.write"]
}