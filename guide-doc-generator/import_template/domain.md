# Domain

#RES.DESC
In the newer v3 approach the input parameter is only name in the format of 'sub_domain.domain'.
##RES.DESC


#RES.COMM
resource "cloudfoundry_domain" "sample" {
  sub_domain = "test"
  domain = "cfapps.stagingazure.hanavlab.ondemand.com"
  org         = "23919ba5-f9b6-4128-a1fb-69890818d25c"
}
##RES.COMM


#RES.SAP
resource "cloudfoundry_domain" "sample" {
  name  = "test.cfapps.stagingazure.hanavlab.ondemand.com"
  org         = "23919ba5-f9b6-4128-a1fb-69890818d25c"
  shared_orgs = ["537e7b58-b3e0-4464-9cad-2deae6120a80", "30edf44a-2d4c-432c-9680-9a61123edcf1"]
}
##RES.SAP

---------------

#DS.DESC

##DS.DESC

#DS.COMM
data "cloudfoundry_domain" "mydomain" {
    sub_domain = "test"
    domain = "cfapps.stagingazure.hanavlab.ondemand.com"
}
##DS.COMM

#DS.SAP
data "cloudfoundry_domain" "mydomain" {
  name = "test.cfapps.stagingazure.hanavlab.ondemand.com"
}
##DS.SAP

