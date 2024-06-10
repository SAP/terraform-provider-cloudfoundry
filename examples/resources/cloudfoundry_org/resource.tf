resource "cloudfoundry_org" "org" {
  name      = "tf-test-iso"
  suspended = false
  labels = {
    env = "test"
  }
  annotations = {
    env-ann = "test-ann"
  }
}