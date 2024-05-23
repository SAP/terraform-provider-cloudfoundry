# space Quota


## Resource




|  SAP Cloudfoundry Provider |Community Cloudfoundry Provider |
| -- | -- |
|  <pre></br>resource "cloudfoundry_space_quota" "large" {</br>  name                     = "large"</br>  allow_paid_service_plans = false</br>  org                      = "ca721b24-e24d-4171-83e1-1ef6bd836b38"</br>  </br>  #Optionals</br>  </br>  total_memory             = 51200</br>  total_routes             = 50</br>  total_services           = 200</br>  instance_memory          = 2048</br>  total_service_keys       = 120</br>  total_app_instances      = 100</br>  total_route_ports        = 5</br>  total_app_tasks          = 10</br></br>  total_app_log_rate_limit = 1000</br>  spaces = [</br>    "ba10cc63-cc43-46b1-a00c-5f2a0d7d992e",</br>  ]</br>}</br></pre> |<pre>resource "cloudfoundry_space_quota" "large" {</br>    name                     = "large"</br>    org                      = "ca721b24-e24d-4171-83e1-1ef6bd836b38"</br>    allow_paid_service_plans = false</br>    total_memory             = 51200</br>    total_routes             = 50</br>    total_services           = 200</br>          </br>    #Optionals</br></br>    instance_memory          = 2048</br>    total_service_keys       = 120</br>    total_app_instances      = 100</br>    total_route_ports        = 5</br>    total_app_tasks          = 10</br>}</br></pre> |

<br/>

## Datasource




|  SAP Cloudfoundry Provider | Community Cloudfoundry Provider  |
| -- | -- |
| <pre>data "cloudfoundry_space_quota" "my_space_quota" {</br>  name = "tf-test-do-not-delete"</br>  org  = "ca721b24-e24d-4171-83e1-1ef6bd836b38"</br>}</br></pre>|<pre>data "cloudfoundry_space_quota" "my_space_quota" {</br>  name = "tf-test-do-not-delete"</br>  org  = "ca721b24-e24d-4171-83e1-1ef6bd836b38"</br>}</br></pre> |  
