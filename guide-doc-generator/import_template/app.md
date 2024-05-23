# App

#RES.DESC
In the newer v3 approach the input parameter is only name in the format of 'sub_domain.domain'.
##RES.DESC


#RES.COMM
resource "cloudfoundry_app" "my-app" {
  name       = "my-app"
  space      = "e6886bba-e263-4b52-aaf1-85d410f15fc8"
  buildpack = "nodejs_buildpack"
  memory     = 512
  path       = "something.zip"
  service_binding {
      service_instance = "xsuaa-tf"
  }
  routes = {
      route = my-app.hello.world.example..com"
  }
}
##RES.COMM


#RES.SAP
resource "cloudfoundry_app" "my-app" {
  name       = "my-app"
  space_name = "tf-space-1"
  org_name   = "PerformanceTeamBLR"
  buildpacks = ["nodejs_buildpack"]
  memory     = "512M"
  path       = "something.zip"
  service_bindings = [
    {
      service_instance = "xsuaa-tf"
    }
  ]
  routes = [
    {
      route = my-app.hello.world.example..com"
    }
  ]
}
##RES.SAP

---------------

#DS.DESC

##DS.DESC

#DS.COMM
data "cloudfoundry_app" "my-app" {
    name_or_id = "my-app"
    space      = "106a411e-3ea3-4d1b-a30d-54a6802bed27"
}
##DS.COMM

#DS.SAP
data "cloudfoundry_app" "my-app" {
  name  = "my-app"
  space_name = "tf-space-1"
  org_name   = "PerformanceTeamBLR"
}
##DS.SAP

