# Space Role


## Resource

Intended to mimic functionality of  A role is a single resource in Cloudfoundry and the resource has been developed to replicate the behaviour of V3 API. For assigning multiple roles, one can create multiple resources to achieve function similar to the community provider with greater control by using for_each for a var with users.

V3 api has clear well-defined separate API's for roles and users and the innumerous endpoints for different roles at both space and org level have been removed. One just needs to specify the necessary `type` now for a particular user as a query while specifying the role. For specifying multiple role types for multiple users, one can create multiple resources now.


|  SAP Cloudfoundry Provider |Community Cloudfoundry Provider |
| -- | -- |
|  <pre>resource "cloudfoundry_space_role" "om1" {</br>  space    = "space-id"</br>  user     = "user-guid"</br>  type     = "space_manager"</br>}</br></br>resource "cloudfoundry_space_role" "od1" {</br>  for_each =  tolist(["username1","username2"])</br>  space    = "space-id"</br>  username = each.value</br>  type     = "space_developer"</br>}</br></br>resource "cloudfoundry_space_role" "oa1" {</br>  space    = "space-id"</br>  username = "username"</br>  type     = "space_auditor"</br>}</br></br>resource "cloudfoundry_space_role" "oa2" {</br>  space  = "space-id"</br>  user   = "user-guid"</br>  type   = "space_auditor"</br>}</br></pre> |<pre>resource "cloudfoundry_space_users" "ou1" {</br>  space              = "space-id"</br>  managers         = ["user-guid"]</br>  developers = ["username1","username2"]</br>  auditors         = ["user-guid", "username"]</br>}</br></br></pre> |

<br/>
