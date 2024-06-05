# Org Role


-----------------
#RES.DESC
Org Roles have been replaced with Org role to match the change in the v3 api definitoins. Roles of `type` are assigned to the individual user.
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