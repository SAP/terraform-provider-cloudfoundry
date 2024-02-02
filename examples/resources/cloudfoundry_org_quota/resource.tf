resource "cloudfoundry_org_quota" "org_quota" {
  name                     = "tf-test-do-not-delete"
  allow_paid_service_plans = true
  instance_memory          = 2048
  total_memory             = 51200
  total_app_instances      = 100
  total_routes             = 50
  total_services           = 200
  total_service_keys       = 120
  total_private_domains    = 40
  total_app_tasks          = 10
  total_route_ports        = 5
  total_app_log_rate_limit = 1000
  organizations = [
    "30edf44a-2d4c-432c-9680-9a61123edcf1",
  ]
}

output "guid" {
  value = resource.cloudfoundry_org_quota.org_quota.id
}