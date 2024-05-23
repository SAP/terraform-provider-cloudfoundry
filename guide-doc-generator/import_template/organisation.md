# Organisation


-----------------
#RES.DESC

##RES.DESC

------------------
#RES.COMM
resource "cloudfoundry_org" "org" {
    name = "tf-test"
    quota = cloudfoundry_quota.runaway.id
}
##RES.COMM

--------------------
#RES.SAP
resource "cloudfoundry_org" "org" {
  name      = "tf-test"
  suspended = false
}
##RES.SAP

---------------

#DS.DESC

##DS.DESC
----------------

#DS.COMM
data "cloudfoundry_org" "org" {
    name = "myorg"    
}
##DS.COMM
-----------------

#DS.SAP
data "cloudfoundry_org" "org" {
  name = "myorg"
}
##DS.SAP