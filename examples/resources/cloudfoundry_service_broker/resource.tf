resource "cloudfoundry_service_broker" "mysql" {
  name     = "broker"
  url      = "example.broker.com"
  username = "test"
  password = "test"
}