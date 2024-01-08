data "cloudfoundry_org" "org" {
  name = "PerformanceTeamBLR"
}

output "id" {
  value = data.cloudfoundry_org.org.id
}

output "labels" {
  value = data.cloudfoundry_org.org.labels
}