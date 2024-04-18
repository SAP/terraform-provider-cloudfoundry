data "cloudfoundry_service_credential_binding" "scbd" {
  service_instance = "5e2976bb-332e-41e1-8be3-53baafea9296"
  app              = "ec6ac2b3-fb79-43c4-9734-000d4299bd59"
}

data "cloudfoundry_service_credential_binding" "scbdaf" {
  service_instance = "e9ec29ca-993d-42e2-9c5b-cb17b1972cce"
  name             = "hifi"
}