# Org Role


-----------------
#RES.DESC
V3 api has clear well-defined separate API's for roles and users and the innumerous endpoints for different roles at both space and org level have been removed. One just needs to specify the necessary `type` now for a particular user as a query while specifying the role. For specifying multiple role types for multiple users, one can create multiple resources now.
##RES.DESC

------------------
#RES.COMM
resource "cloudfoundry_org_users" "ou1" {
  org              = "organization-id"
  managers         = ["user-guid"]
  billing_managers = ["username1","username2"]
  auditors         = ["user-guid", "username"]
}

##RES.COMM

--------------------
#RES.SAP
resource "cloudfoundry_org_role" "om1" {
  for_each =  tolist(["username1","username2"])
  org      = "organization-id"
  username = each.value
  type     = "organization_billing_manager"
}

resource "cloudfoundry_org_role" "oa1" {
  org              = "organization-id"
  username = "username"
  type     = "organization_auditor"
}

##RES.SAP