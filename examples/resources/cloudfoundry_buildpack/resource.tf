resource "cloudfoundry_buildpack" "mybuildpack" {
  name     = "hi"
  position = 2
  stack    = "cflinuxfs3"
  enabled  = false
  locked   = false
  labels   = { "hi" : "fi" }
  path     = "somethin.zip"
} 