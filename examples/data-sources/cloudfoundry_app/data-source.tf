data "cloudfoundry_app" "http-bin-server" {
  name       = "tf-test-do-not-delete-http-bin"
  space_name = "tf-space-1"
  org_name   = "PerformanceTeamBLR"
}

output "id" {
  value = data.cloudfoundry_app.http-bin-server.id
}

output "space" {
  value = data.cloudfoundry_app.http-bin-server.space_name
}

output "name" {
  value = data.cloudfoundry_app.http-bin-server.name
}
output "environment" {
  value = data.cloudfoundry_app.http-bin-server.environment
}
output "instances" {
  value = data.cloudfoundry_app.http-bin-server.instances
}
output "memory" {
  value = data.cloudfoundry_app.http-bin-server.memory
}
output "disk_quota" {
  value = data.cloudfoundry_app.http-bin-server.disk_quota
}
output "routes" {
  value = data.cloudfoundry_app.http-bin-server.routes
}
output "buildpacks" {
  value = data.cloudfoundry_app.http-bin-server.buildpacks
}
output "service_bindings" {
  value = data.cloudfoundry_app.http-bin-server.service_bindings
}