# cloudfoundry_buildpack

Provides a Cloud Foundry resource to manage buildpacks.

|  SAP Cloudfoundry Provider |Community Cloudfoundry Provider |
| -- | -- |
|  <pre>resource "cloudfoundry_buildpack" "bp" {</br>  name     = "hi"</br>  position = 1</br>  stack    = "cflinuxfs3"</br>  enabled  = false</br>  locked   = true</br>  labels   = { "hi" : "fi" }</br>  path     = "something.zip"</br>}</br></pre> |<pre>resource "cloudfoundry_buildpack" "bp" {</br>  name     = "hi"</br>  position = 1</br>  enabled  = false</br>  locked   = true</br>  labels   = { "hi" : "fi" }</br>  path     = "something.zip"</br>}</br></pre> |

## Differences
> [!NOTE]  
> 🔵 Required  🟢 Optional 🟠 Computed  🔴 Not present

| Attribute name| SAP Cloudfoundry Provider(new)|  Community Provider(old) | Description
|---| ---| ---| ---| 
|path| 🟢|  🔵  | - |
|stack |  🟢 |🔴| - |