# Service binding


## Resource




|  SAP Cloudfoundry Provider |Community Cloudfoundry Provider |
| -- | -- |
|  <pre>resource "cloudfoundry_service_credential_binding" "scb1" {</br>  type             = "key"</br>  name             = "hifi"</br>  service_instance = "e9ec29ca-993d-42e2-9c5b-cb17b1972cce"</br>}</br></br>resource "cloudfoundry_service_credential_binding" "scb7" {</br>  type             = "app"</br>  name             = "hifi"</br>  service_instance = "e9ec29ca-993d-42e2-9c5b-cb17b1972cce"</br>  app              = "ec6ac2b3-fb79-43c4-9734-000d4299bd59"</br>}</br></pre> |<pre>resource "cloudfoundry_service_key" "redis1-key1" {</br>  name = "hifi"</br>  service_instance = "e9ec29ca-993d-42e2-9c5b-cb17b1972cce"</br>}</br></pre> |

<br/>

## Datasource




|  SAP Cloudfoundry Provider | Community Cloudfoundry Provider  |
| -- | -- |
| <pre>data "cloudfoundry_service_credential_binding" "scbd" {</br>  service_instance = "5e2976bb-332e-41e1-8be3-53baafea9296"</br>  app              = "ec6ac2b3-fb79-43c4-9734-000d4299bd59"</br>}</br></br>data "cloudfoundry_service_credential_binding" "my-key" {</br>    name             = "my-service-key"</br>    service_instance = "e9ec29ca-993d-42e2-9c5b-cb17b1972cce"</br>}</br></pre>|<pre>data "cloudfoundry_service_key" "my-key" {</br>    name             = "my-service-key"</br>    service_instance = "e9ec29ca-993d-42e2-9c5b-cb17b1972cce"</br>}</br></pre> |  
