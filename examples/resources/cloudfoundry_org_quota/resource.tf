resource "cloudfoundry_org_quota" "org_quota" {
  name = "tf-test-org-quota"
  allow_paid_service_plans = true
  instance_memory = 2048
  total_memory = 51200
  total_app_instances = 100
  total_routes = 50
  total_services = 200
  total_route_ports = 5
  total_app_log_rate_limit = 1000
}

output "guid" {
  value = resource.cloudfoundry_org_quota.org_quota.id
}