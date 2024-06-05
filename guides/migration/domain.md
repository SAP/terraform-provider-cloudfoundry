# Domain


## Resource

In the newer v3 approach the input parameter is only name in the format of 'sub_domain.domain'.


| Community Cloudfoundry Provider | SAP Cloudfoundry Provider |
| -- | -- |
| <pre>resource "cloudfoundry_domain" "sample" {</br>  sub_domain = "test"</br>  domain = "cfapps.stagingazure.hanavlab.ondemand.com"</br>}</br></pre> | <pre>resource "cloudfoundry_domain" "sample" {</br>  name  = "test.cfapps.stagingazure.hanavlab.ondemand.com"</br>}</br></pre> |

<br/>

## Datasource





| Community Cloudfoundry Provider | SAP Cloudfoundry Provider |
| -- | -- |
| <pre>data "cloudfoundry_domain" "l" {</br>    sub_domain = "test"</br>    domain = "cfapps.stagingazure.hanavlab.ondemand.com"</br>}</br></pre> | <pre>data "cloudfoundry_domain" "mydomain" {</br>  name = "test.cfapps.stagingazure.hanavlab.ondemand.com"</br>}</br></pre> |
