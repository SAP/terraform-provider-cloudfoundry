# cloudfoundry_user

Gets information on Cloud Foundry users with a given username.

|  SAP Cloudfoundry Provider | Community Cloudfoundry Provider  |
| -- | -- |
| <pre>data "cloudfoundry_user" "myuser" {</br>  name = "myuser"</br>}</br></pre>|<pre>data "cloudfoundry_user" "myuser" {</br>    name = "myuser"    </br>}</br></pre> |  

## Differences
> [!NOTE]  
> 🔵 Required  🟢 Optional 🟠 Computed  🔴 Not present


| Attribute name|  SAP Cloudfoundry Provider(new)|  Community Provider(old) | Description
|---| ---| ---| ---| 
|org_id| 🔴 |🟢 | For fetching specific user under a particular org in current provider, one can set the `org` attribute in [`cloudfoundry_users`](https://github.com/SAP/terraform-provider-cloudfoundry/blob/migration_docs/docs/data-sources/users.md) resource and then from `users` attribute output, filter the user with `username` desired
|users| 🟠 |🔴| 
|id |  🔴|🟠|  The current provider returns multiple users if available with same user name in the `users` attribute unlike the community provider. Therefore the id is present in the respective user resources in `users` output