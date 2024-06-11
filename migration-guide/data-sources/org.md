# cloudfoundry_org

Gets information on a Cloud Foundry organization.

|  SAP Cloudfoundry Provider | Community Cloudfoundry Provider  |
| -- | -- |
| <pre>data "cloudfoundry_org" "org" {</br>  name = "myorg"</br>}</br></pre>|<pre>data "cloudfoundry_org" "org" {</br>    name = "myorg"    </br>}</br></pre> |  

## Differences

> [!NOTE]  
> 🔵 Required  🟢 Optional 🟠 Computed  🔴 Not present


| Attribute name|  SAP Cloudfoundry Provider(new)|  Community Provider(old) | Description
|---| ---| ---| ---| 
|quota|  🟠| 🔴 | 
|suspended | 🟠 |🔴|  
