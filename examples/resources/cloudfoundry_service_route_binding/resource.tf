resource "cloudfoundry_service_route_binding" "srb7" {
  service_instance = "3a8588f9-f846-444f-ab9e-48282f06449b"
  route            = "3966c2fb-d84d-462d-82a5-a81cf7cdab20"
  labels           = { "hi" : "fi" }
}