data "cloudfoundry_user" "my_user" {
  name = "debaditya.ray@sap.com"
}

output "labels" {
  value = data.cloudfoundry_user.my_user.users.0.labels.enviroment
}