# Org Role


## Resource

Org Roles have been replaced with Org role to match the change in the v3 api definitoins. Roles of `type` are assigned to the individual user.


| Community Cloudfoundry Provider | SAP Cloudfoundry Provider |
| -- | -- |
| <pre>resource "cloudfoundry_org_users" "ou1" {</br>  org              = "organization-id"</br>  managers         = ["user-guid"]</br>  billing_managers = ["username"]</br>  auditors         = ["user-guid", "username"]</br>}</br></br></pre> | <pre>resource "cloudfoundry_org_role" "om1" {</br>  org              = "organization-id"</br>  username = "username"</br>  type     = "organization_billing_manager"</br>}</br></br>resource "cloudfoundry_org_role" "oa1" {</br>  org              = "organization-id"</br>  username = "username"</br>  type     = "organization_auditor"</br>}</br></br></pre> |

<br/>
