resource "cloudfoundry_org" "org" {
  name      = "tf-test"
  suspended = false
  labels = {
    env = "test"
  }
  annotations = {
    env-ann = "test-ann"
  }
}