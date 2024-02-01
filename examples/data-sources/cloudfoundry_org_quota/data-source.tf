data "cloudfoundry_org_quota" "org_quota" {
  name = "tf-test-org-quota"
}

output "guid" {
  value = data.cloudfoundry_org_quota.org_quota.id
}
output "instance_memory" {
  value = data.cloudfoundry_org_quota.org_quota.instance_memory
}