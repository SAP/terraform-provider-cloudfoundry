resource "cloudfoundry_isolation_segment_entitlement" "isosegment" {
  segment = "63ae51b9-9073-4409-81b0-3704b8de85dd"
  orgs    = ["c8e454cc-7a24-4d71-b146-51d69538acfb"]
  default = true
} 