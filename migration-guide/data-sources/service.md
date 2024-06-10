# cloudfoundry_service 

Get Service Offering and its related plans.

|  SAP Cloudfoundry Provider | Community Cloudfoundry Provider  |
| -- | -- |
| <pre>data "cloudfoundry_service" "redis" {</br>  name = "p-redis"</br>}</br></pre>|<pre>data "cloudfoundry_service" "redis" {</br>    name = "p-redis"    </br>}</br></pre> | 

## Differences
> [!NOTE]  
> 🔵 Required  🟢 Optional 🟠 Computed  🔴 Not present

| Attribute name|  SAP Cloudfoundry Provider(new)|  Community Provider(old) | Description
|---| ---| ---| ---| 
|service_broker_guid|  🔴 | 🟢 | `service_broker_guid`attribute has been changed to `service_broker` 
|service_broker|  🟢 | 🔴 | 
|service_broker_name|  🔴  |  🟠 |
|space|  🔴  |  🟢 |

