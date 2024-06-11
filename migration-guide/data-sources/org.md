# cloudfoundry_org

Gets information on a Cloud Foundry organization.

|  SAP Cloud Foundry Provider | Community Cloud Foundry Provider  |
| -- | -- |
| <pre>data "cloudfoundry_org" "org" {</br>  name = "myorg"</br>}</br></pre>|<pre>data "cloudfoundry_org" "org" {</br>    name = "myorg"    </br>}</br></pre> |  

## Differences

> [!NOTE]  
> 🔵 Required  🟢 Optional 🟠 Computed  🔴 Not present

| Attribute name | SAP Cloud Foundry Provider (new)|  Community Cloud Foundry Provider (old) | Description |
| --- | --- | --- | --- |
| quota | 🟠| 🔴 | - |
| suspended | 🟠 | 🔴 | - |
