# Organisation


## Resource




|  SAP Cloudfoundry Provider |Community Cloudfoundry Provider |
| -- | -- |
|  <pre>resource "cloudfoundry_org" "org" {</br>  name      = "tf-test"</br>  suspended = false</br>}</br></pre> |<pre>resource "cloudfoundry_org" "org" {</br>    name = "tf-test"</br>    quota = cloudfoundry_quota.runaway.id</br>}</br></pre> |

<br/>

## Datasource





|  SAP Cloudfoundry Provider | Community Cloudfoundry Provider  |
| -- | -- |
| <pre>data "cloudfoundry_org" "org" {</br>  name = "myorg"</br>}</br></pre>|<pre>data "cloudfoundry_org" "org" {</br>    name = "myorg"    </br>}</br></pre> |  
