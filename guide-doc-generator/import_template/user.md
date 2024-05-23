# user


-----------------
#RES.DESC

##RES.DESC

------------------
#RES.COMM
resource "cloudfoundry_user" "admin-service-user" {

    name = "cf-admin"
    password = "Passw0rd"

    given_name = "John"
    family_name = "Doe"

    groups = [ "cloud_controller.admin", "scim.read", "scim.write" ]
}
##RES.COMM

--------------------
#RES.SAP
resource "cloudfoundry_user" "my_user" {
  id          = "test-user567"
  annotations = { purpose : "testing" }
}
##RES.SAP

---------------

#DS.DESC
##DS.DESC
----------------

#DS.COMM
data "cloudfoundry_user" "myuser" {
    name = "myuser"    
}
##DS.COMM
-----------------

#DS.SAP
data "cloudfoundry_user" "myuser" {
  name = "myuser"
}
##DS.SAP