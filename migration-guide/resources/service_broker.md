# cloudfoundry_service_broker

Provides a Cloud Foundry resource for managing service brokers.

|  SAP Cloudfoundry Provider |Community Cloudfoundry Provider |
| -- | -- |
|  <pre>resource "cloudfoundry_service_broker" "mysql" { </br>  name     = "broker"</br>  url      = "example.broker.com"</br>  username = "test"</br>  password = "test"</br>}</br></pre> |<pre>resource "cloudfoundry_service_broker" "mysql" { </br>  name     = "broker"</br>  url      = "example.broker.com"</br>  username = "test"</br>  password = "test"</br>}</br></pre> |

## Differences
> [!NOTE]  
> 🔵 Required  🟢 Optional 🟠 Computed  🔴 Not present

| Attribute name| SAP Cloudfoundry Provider(new)|  Community Provider(old) | Description
|---| ---| ---| ---| 
|fail_when_catalog_not_accessible | 🔴|  🟢  | - |
|service_plans  |   🔴 |🟠| - |
|services  |   🔴 |🟠| - | 