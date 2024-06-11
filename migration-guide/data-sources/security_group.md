# cloudfoundry_security_group

Gets information on a Cloud Foundry application security group. Named as [`cloudfoundry_asg`](https://github.com/cloudfoundry-community/terraform-provider-cloudfoundry/blob/main/docs/data-sources/asg.md) in the community provider.

|  SAP Cloud Foundry Provider | Community Cloud Foundry Provider  |
| -- | -- |
| <pre>data "cloudfoundry_security_group" "public" {</br>  name = "public_networks"</br>}</br></pre>|<pre>data "cloudfoundry_asg" "public" {</br>    name = "public_networks"</br>}</br></pre> |  

## Differences

> [!NOTE]  
> 🔵 Required  🟢 Optional 🟠 Computed  🔴 Not present

| Attribute name | SAP Cloud Foundry Provider (new)|  Community Cloud Foundry Provider (old) | Description |
| --- | --- | --- | --- |
| globally_enabled_running | 🟠| 🔴 | - |
| globally_enabled_staging | 🟠| 🔴 | - |
| running_spaces | 🟠 | 🔴 | - |
| staging_spaces | 🟠 | 🔴 | - |
| rules | 🟠 | 🔴 | - |
| labels | 🟠 | 🔴 | - |
| annotations | 🟠 | 🔴 | - |
