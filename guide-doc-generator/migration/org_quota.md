# Org Quota


## Resource

In the v3 API only `allow_paid_service_plans` and `name` are required fields. Orgs can be bound in the new provider.



|  SAP Cloudfoundry Provider |Community Cloudfoundry Provider |
| -- | -- |
|  <pre></br>resource "cloudfoundry_org_quota" "large" {</br>  name                     = "large"</br>  allow_paid_service_plans = false</br></br>  #Optionals</br>  </br>  total_memory             = 51200</br>  total_routes             = 50</br>  total_services           = 200</br>  instance_memory          = 2048</br>  total_service_keys       = 120</br>  total_app_instances      = 100</br>  total_route_ports        = 5</br>  total_private_domains    = 40</br>  total_app_tasks          = 10</br></br>  total_app_log_rate_limit = 1000</br>  orgs = [</br>    "ba10cc63-cc43-46b1-a00c-5f2a0d7d992e",</br>  ]</br>}</br></br></pre> |<pre>resource "cloudfoundry_org_quota" "large" {</br>    name                     = "large"</br>    allow_paid_service_plans = false</br>    total_memory             = 51200</br>    total_routes             = 50</br>    total_services           = 200</br>          </br>    #Optionals</br></br>    instance_memory          = 2048</br>    total_service_keys       = 120</br>    total_app_instances      = 100</br>    total_route_ports        = 5</br>    total_private_domains    = 40</br>    total_app_tasks          = 10</br>}</br></br></pre> |

<br/>

## Datasource





|  SAP Cloudfoundry Provider | Community Cloudfoundry Provider  |
| -- | -- |
| <pre>data "cloudfoundry_org_quota" "org_quota" {</br>  name = "myquota"</br>}</br></pre>|<pre>data "cloudfoundry_org_quota" "org_quota" {</br>  name = "myquota"</br>}</br></pre> |  
