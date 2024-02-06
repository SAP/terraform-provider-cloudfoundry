data "cloudfoundry_space_quota" "my_space_quota" {
  name = "tf-test-do-not-delete"
  org  = "ca721b24-e24d-4171-83e1-1ef6bd836b38"
}

output "id" {
  value = data.cloudfoundry_space_quota.my_space_quota.id
}