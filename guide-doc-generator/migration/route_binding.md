# Service binding


## Resource




|  SAP Cloudfoundry Provider |Community Cloudfoundry Provider |
| -- | -- |
|  <pre>resource "cloudfoundry_service_route_binding" "srb" {</br>  service_instance = "3a8588f9-f846-444f-ab9e-48282f06449b"</br>  route            = "3966c2fb-d84d-462d-82a5-a81cf7cdab20"</br>  labels           = { "hi" : "fi" }</br>}</br></pre> |<pre>resource "cloudfoundry_route_service_binding" "srb" {</br>  service_instance = "3a8588f9-f846-444f-ab9e-48282f06449b"</br>  route            = "3966c2fb-d84d-462d-82a5-a81cf7cdab20"</br>}</br></pre> |

<br/>
