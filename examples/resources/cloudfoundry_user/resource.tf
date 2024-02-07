resource "cloudfoundry_user" "my_user" {
  id          = "test-user567"
  annotations = { purpose : "testing" }
}
