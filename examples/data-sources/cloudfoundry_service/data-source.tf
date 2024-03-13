data "cloudfoundry_service" "xsuaa-offering" {
  name = "xsuaa"
}

output "serviceplans" {
  value = data.cloudfoundry_service.xsuaa-offering.service_plans
}