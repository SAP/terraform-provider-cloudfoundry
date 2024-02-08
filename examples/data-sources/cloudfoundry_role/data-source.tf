data "cloudfoundry_role" "my_role" {
  id = "e17839d9-cd4f-4e4b-baf0-18786f12fede"
}

output "role_object" {
  value = data.cloudfoundry_role.my_role
}