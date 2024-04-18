resource "cloudfoundry_service_credential_binding" "scb7" {
  type             = "app"
  name             = "hifi"
  service_instance = "e9ec29ca-993d-42e2-9c5b-cb17b1972cce"
  app              = "ec6ac2b3-fb79-43c4-9734-000d4299bd59"
}

resource "cloudfoundry_service_credential_binding" "scb1" {
  type             = "key"
  name             = "hifi"
  service_instance = "e9ec29ca-993d-42e2-9c5b-cb17b1972cce"
}