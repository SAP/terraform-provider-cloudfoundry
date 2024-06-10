# cloudfoundry_org (Resource)

Provides a Cloud Foundry resource for managing Cloud Foundry organizations


|  SAP Cloudfoundry Provider |Community Cloudfoundry Provider |
| -- | -- |
|  <pre>resource "cloudfoundry_org" "org" {</br>  name      = "tf-test"</br>  suspended = false</br>}</br></pre> |<pre>resource "cloudfoundry_org" "org" {</br>    name = "tf-test"</br>    quota = cloudfoundry_quota.runaway.id</br>}</br></pre> |

## Differences
> [!NOTE]  
> 🔵 Required  🟢 Optional 🟠 Computed  🔴 Not present

| Attribute name|  SAP Cloudfoundry Provider(new)|  Community Provider(old) | Description
|---| ---| ---| ---| 
|quota|  🟠| 🟢 | One cannot set quota as it is a read-only attribute in the current provider. For setting quota use resource [`cloudfoundry_org_quota`](https://github.com/SAP/terraform-provider-cloudfoundry/blob/main/docs/resources/org_quota.md).
|suspended|   🟢 |🔴  |  
|delete_recursive_allowed | 🔴  | 🟢| V3 API by default follows recursive deletion.
