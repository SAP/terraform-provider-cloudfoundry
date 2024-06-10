# cloudfoundry_org_quota

Gets information on a Cloud Foundry organization quota.  

|  SAP Cloudfoundry Provider | Community Cloudfoundry Provider  |
| -- | -- |
| <pre>data "cloudfoundry_org_quota" "org_quota" {</br>  name = "myquota"</br>}</br></pre>|<pre>data "cloudfoundry_org_quota" "org_quota" {</br>  name = "myquota"</br>}</br></pre> |  

## Differences
> [!NOTE]  
> 🔵 Required  🟢 Optional 🟠 Computed  🔴 Not present


| Attribute name|  SAP Cloudfoundry Provider(new)|  Community Provider(old) | Description
|---| ---| ---| ---| 
|orgs| 🟠 |🔴| 
|total_app_log_rate_limit|  🟠 |🔴|
|total_memory|  🟠 |🔴| 
|total_routes|  🟠 |🔴|
|total_services| 🟠 |🔴| 
|allow_paid_service_plans| 🟠 |🔴| 
|instance_memory | 🟠 |🔴| 
|total_app_instances| 🟠 |🔴| 
|total_services| 🟠 |🔴| 
|total_app_tasks| 🟠 |🔴| 
|total_private_domains| 🟠 |🔴| 
|total_route_ports| 🟠 |🔴| 
|total_service_keys| 🟠 |🔴| 