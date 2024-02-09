data "cloudfoundry_users" "my_users" {
  org = "784b4cd0-4771-4e4d-9052-a07e178bae56"
}

output "users" {
  value = data.cloudfoundry_users.my_users
}