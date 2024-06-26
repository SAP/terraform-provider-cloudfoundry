# cloudfoundry_stack

Gets information on a Cloud Foundry stack.

|  SAP Cloudfoundry Provider | Community Cloudfoundry Provider  |
| -- | -- |
| <pre>data "cloudfoundry_stack" "mystack" {</br>    name = "my_custom_stack"</br>}</br></pre>|<pre>data "cloudfoundry_stack" "mystack" {</br>    name = "my_custom_stack"</br>}</br></pre> |  

## Differences
> [!NOTE]  
> 🔵 Required  🟢 Optional 🟠 Computed  🔴 Not present

| Attribute name|  SAP Cloudfoundry Provider(new)| Community Provider(old) |Description
|---| ---| ---| ---| 
|build_rootfs_image | 🟠|🔴|  - |
|run_rootfs_image | 🟠|🔴|  - |
|default|  🟠|🔴 | - |