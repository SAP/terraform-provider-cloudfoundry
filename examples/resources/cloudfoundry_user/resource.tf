resource "cloudfoundry_user" "my_user" {
  username    = "test"
  email       = "test@gmail.com"
  password    = "test123"
  given_name  = "test"
  family_name = "test"
  annotations = { "purpose" : "testing", hi : "hello" }
}