resource "cloudfoundry_isolation_segment" "isosegment" {
  name   = "hifi"
  labels = { "purpose" : "testing" }
} 