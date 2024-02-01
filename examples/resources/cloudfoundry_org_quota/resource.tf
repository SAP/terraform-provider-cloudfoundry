resource "cloudfoundry_org_quota" "org_quota" {
  name                     = "tf-test-org-quota"
  allow_paid_service_plans = true
  instance_memory          = 2048
  total_memory             = 51200
  total_app_instances      = 100
  total_routes             = 50
  total_services           = 200
  total_route_ports        = 5
  total_app_log_rate_limit = 1000
  organizations = [
    "ff4ab280-90b7-46ab-9877-12e3820d992e",
    "21cdfd82-b828-4802-bb77-55df34aea6da",
  ]
}

output "guid" {
  value = resource.cloudfoundry_org_quota.org_quota.id
}