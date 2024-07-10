# user


## Resource




|  SAP Cloudfoundry Provider |Community Cloudfoundry Provider |
| -- | -- |
|  <pre>resource "cloudfoundry_user" "my_user" {</br>  username    = "test"</br>  email       = "test@gmail.com"</br>  password    = "test123"</br>  given_name  = "test"</br>  family_name = "test"</br>  annotations = { "purpose" : "testing", hi : "hello" }</br>}</br></pre> |<pre>resource "cloudfoundry_user" "my_user" {</br>  name    = "test"</br>  email       = "test@gmail.com"</br>  password    = "test123"</br>  given_name  = "test"</br>  family_name = "test"</br>}</br></pre> |

<br/>

## Datasource




|  SAP Cloudfoundry Provider | Community Cloudfoundry Provider  |
| -- | -- |
| <pre>data "cloudfoundry_user" "myuser" {</br>  name = "myuser"</br>}</br></pre>|<pre>data "cloudfoundry_user" "myuser" {</br>    name = "myuser"    </br>}</br></pre> |  
