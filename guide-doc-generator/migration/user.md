# user


## Resource




|  SAP Cloudfoundry Provider |Community Cloudfoundry Provider |
| -- | -- |
|  <pre>resource "cloudfoundry_user" "my_user" {</br>  id          = "test-user567"</br>  annotations = { purpose : "testing" }</br>}</br></pre> |<pre>resource "cloudfoundry_user" "admin-service-user" {</br></br>    name = "cf-admin"</br>    password = "Passw0rd"</br></br>    given_name = "John"</br>    family_name = "Doe"</br></br>    groups = [ "cloud_controller.admin", "scim.read", "scim.write" ]</br>}</br></pre> |

<br/>

## Datasource




|  SAP Cloudfoundry Provider | Community Cloudfoundry Provider  |
| -- | -- |
| <pre>data "cloudfoundry_user" "myuser" {</br>  name = "myuser"</br>}</br></pre>|<pre>data "cloudfoundry_user" "myuser" {</br>    name = "myuser"    </br>}</br></pre> |  
