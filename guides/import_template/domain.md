# Domain

#RES.DESC
In the newer v3 approach the input parameter is only name in the format of 'sub_domain.domain'.
##RES.DESC


#RES.COMM
resource "cloudfoundry_domain" "sample" {
  sub_domain = "test"
  domain = "cfapps.stagingazure.hanavlab.ondemand.com"
}
##RES.COMM


#RES.SAP
resource "cloudfoundry_domain" "sample" {
  name  = "test.cfapps.stagingazure.hanavlab.ondemand.com"
}
##RES.SAP

---------------

#DS.DESC

##DS.DESC

#DS.COMM
data "cloudfoundry_domain" "l" {
    sub_domain = "test"
    domain = "cfapps.stagingazure.hanavlab.ondemand.com"
}
##DS.COMM

#DS.SAP
data "cloudfoundry_domain" "mydomain" {
  name = "test.cfapps.stagingazure.hanavlab.ondemand.com"
}
##DS.SAP

