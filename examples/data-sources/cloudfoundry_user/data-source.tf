data "cloudfoundry_user" "my_user" {
  name = "test123@example.com"
}

output "labels" {
  value = data.cloudfoundry_user.my_user.users.0.labels.enviroment
}