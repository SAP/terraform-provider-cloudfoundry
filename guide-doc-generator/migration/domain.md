# Domain


## Resource

In the newer v3 approach the input parameter is only name in the format of 'sub_domain.domain'.


|  SAP Cloudfoundry Provider |Community Cloudfoundry Provider |
| -- | -- |
|  <pre>resource "cloudfoundry_domain" "sample" {</br>  name  = "test.cfapps.stagingazure.hanavlab.ondemand.com"</br>  org         = "23919ba5-f9b6-4128-a1fb-69890818d25c"</br>  shared_orgs = ["537e7b58-b3e0-4464-9cad-2deae6120a80", "30edf44a-2d4c-432c-9680-9a61123edcf1"]</br>}</br></pre> |<pre>resource "cloudfoundry_domain" "sample" {</br>  sub_domain = "test"</br>  domain = "cfapps.stagingazure.hanavlab.ondemand.com"</br>  org         = "23919ba5-f9b6-4128-a1fb-69890818d25c"</br>}</br></pre> |

<br/>

## Datasource





|  SAP Cloudfoundry Provider | Community Cloudfoundry Provider  |
| -- | -- |
| <pre>data "cloudfoundry_domain" "mydomain" {</br>  name = "test.cfapps.stagingazure.hanavlab.ondemand.com"</br>}</br></pre>|<pre>data "cloudfoundry_domain" "mydomain" {</br>    sub_domain = "test"</br>    domain = "cfapps.stagingazure.hanavlab.ondemand.com"</br>}</br></pre> |  
