# space


## Resource




|  SAP Cloudfoundry Provider |Community Cloudfoundry Provider |
| -- | -- |
|  <pre>resource "cloudfoundry_space" "space" {</br>  name      = "space"</br>  org       = "ca721b24-e24d-4171-83e1-1ef6bd836b38"</br>  allow_ssh = "true"</br>}</br></pre> |<pre>resource "cloudfoundry_space" "space" {</br>    name = "space"</br>    org  = "ca721b24-e24d-4171-83e1-1ef6bd836b38"</br>    quota = "dd457c79-f7c9-4828-862b-35843d3b646d"</br>    asgs = [ "ba10cc63-cc43-46b1-a00c-5f2a0d7d992e" ]</br>    allow_ssh = true</br>}</br></pre> |

<br/>

## Datasource




|  SAP Cloudfoundry Provider | Community Cloudfoundry Provider  |
| -- | -- |
| <pre>data "cloudfoundry_space" "space" {</br>  name = "myspace"</br>  org  = "ca721b24-e24d-4171-83e1-1ef6bd836b38"</br>}</br></pre>|<pre>data "cloudfoundry_space" "space" {</br>    name = "myspace"</br>    org_name = "org"</br>}</br></pre> |  
