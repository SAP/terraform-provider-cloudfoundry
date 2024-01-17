data "cloudfoundry_space" "space" {
  name = "PerformanceTeamBLR"
  org  = "784b4cd0-4771-4e4d-9052-a07e178bae56"
}

output "id" {
  value = data.cloudfoundry_space.space.id
}