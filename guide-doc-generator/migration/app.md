# App


## Resource

In the newer v3 approach the input parameter is only name in the format of 'sub_domain.domain'.


|  SAP Cloudfoundry Provider |Community Cloudfoundry Provider |
| -- | -- |
|  <pre>resource "cloudfoundry_app" "my-app" {</br>  name       = "my-app"</br>  space_name = "tf-space-1"</br>  org_name   = "PerformanceTeamBLR"</br>  buildpacks = ["nodejs_buildpack"]</br>  memory     = "512M"</br>  path       = "something.zip"</br>  service_bindings = [</br>    {</br>      service_instance = "xsuaa-tf"</br>    }</br>  ]</br>  routes = [</br>    {</br>      route = my-app.hello.world.example..com"</br>    }</br>  ]</br>}</br></pre> |<pre>resource "cloudfoundry_app" "my-app" {</br>  name       = "my-app"</br>  space      = "e6886bba-e263-4b52-aaf1-85d410f15fc8"</br>  buildpack = "nodejs_buildpack"</br>  memory     = 512</br>  path       = "something.zip"</br>  service_binding {</br>      service_instance = "xsuaa-tf"</br>  }</br>  routes = {</br>      route = my-app.hello.world.example..com"</br>  }</br>}</br></pre> |

<br/>

## Datasource





|  SAP Cloudfoundry Provider | Community Cloudfoundry Provider  |
| -- | -- |
| <pre>data "cloudfoundry_app" "my-app" {</br>  name  = "my-app"</br>  space_name = "tf-space-1"</br>  org_name   = "PerformanceTeamBLR"</br>}</br></pre>|<pre>data "cloudfoundry_app" "my-app" {</br>    name_or_id = "my-app"</br>    space      = "106a411e-3ea3-4d1b-a30d-54a6802bed27"</br>}</br></pre> |  
