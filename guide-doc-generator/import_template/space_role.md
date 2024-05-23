# Space Role


-----------------
#RES.DESC
V3 api has clear well-defined separate API's for roles and users and the innumerous endpoints for different roles at both space and org level have been removed. One just needs to specify the necessary `type` now for a particular user as a query while specifying the role. For specifying multiple role types for multiple users, one can create multiple resources now.
##RES.DESC

------------------
#RES.COMM
resource "cloudfoundry_space_users" "ou1" {
  space              = "space-id"
  managers         = ["user-guid"]
  developers = ["username1","username2"]
  auditors         = ["user-guid", "username"]
}

##RES.COMM

--------------------
#RES.SAP
resource "cloudfoundry_space_role" "om1" {
  space    = "space-id"
  user     = "user-guid"
  type     = "space_manager"
}

resource "cloudfoundry_space_role" "od1" {
  for_each =  tolist(["username1","username2"])
  space    = "space-id"
  username = each.value
  type     = "space_developer"
}

resource "cloudfoundry_space_role" "oa1" {
  space    = "space-id"
  username = "username"
  type     = "space_auditor"
}

resource "cloudfoundry_space_role" "oa2" {
  space  = "space-id"
  user   = "user-guid"
  type   = "space_auditor"
}
##RES.SAP