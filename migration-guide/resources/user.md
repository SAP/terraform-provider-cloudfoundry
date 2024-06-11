# cloudfoundry_user

Provides a Cloud Foundry resource for registering users. Intended to mimic functionality of [`cloudfoundry_user`](https://github.com/cloudfoundry-community/terraform-provider-cloudfoundry/blob/main/docs/resources/user.md) resource in the community provider. However, it does not create users in the origin DB with any roles.

|  SAP Cloud Foundry Provider |Community Cloud Foundry Provider |
| -- | -- |
|  <pre>resource "cloudfoundry_user" "my_user" {</br>  id          = "test-user567"</br>  annotations = { purpose : "testing" }</br>}</br></pre> |<pre>resource "cloudfoundry_user" "admin-service-user" {</br></br>    name = "cf-admin"</br>    password = "Passw0rd"</br></br>    given_name = "John"</br>    family_name = "Doe"</br></br>    groups = [ "cloud_controller.admin", "scim.read", "scim.write" ]</br>}</br></pre> |

## Differences

> [!NOTE]  
> 🔵 Required  🟢 Optional 🟠 Computed  🔴 Not present

TODO: Add differences between the providers for the resource.
