# Buildpack


## Resource




|  SAP Cloudfoundry Provider |Community Cloudfoundry Provider |
| -- | -- |
|  <pre>resource "cloudfoundry_buildpack" "bp" {</br>  name     = "hi"</br>  position = 1</br>  stack    = "cflinuxfs3"</br>  enabled  = false</br>  locked   = true</br>  labels   = { "hi" : "fi" }</br>  path     = "something.zip"</br>}</br></pre> |<pre>resource "cloudfoundry_buildpack" "bp" {</br>  name     = "hi"</br>  position = 1</br>  enabled  = false</br>  locked   = true</br>  labels   = { "hi" : "fi" }</br>  path     = "something.zip"</br>}</br></pre> |

<br/>
