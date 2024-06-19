resource "cloudfoundry_space" "space" {
  name      = "space"
  org       = "c8e454cc-7a24-4d71-b146-51d69538acfb"
  allow_ssh = "true"
  labels    = { test : "pass", purpose : "prod" }
}
