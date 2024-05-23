# space


-----------------
#RES.DESC

##RES.DESC

------------------
#RES.COMM
resource "cloudfoundry_space" "space" {
    name = "space"
    org  = "ca721b24-e24d-4171-83e1-1ef6bd836b38"
    quota = "dd457c79-f7c9-4828-862b-35843d3b646d"
    asgs = [ "ba10cc63-cc43-46b1-a00c-5f2a0d7d992e" ]
    allow_ssh = true
}
##RES.COMM

--------------------
#RES.SAP
resource "cloudfoundry_space" "space" {
  name      = "space"
  org       = "ca721b24-e24d-4171-83e1-1ef6bd836b38"
  allow_ssh = "true"
}
##RES.SAP

---------------

#DS.DESC
##DS.DESC
----------------

#DS.COMM
data "cloudfoundry_space" "space" {
    name = "myspace"
    org_name = "org"
}
##DS.COMM
-----------------

#DS.SAP
data "cloudfoundry_space" "space" {
  name = "myspace"
  org  = "ca721b24-e24d-4171-83e1-1ef6bd836b38"
}
##DS.SAP