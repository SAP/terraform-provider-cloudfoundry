# user


-----------------
#RES.DESC

##RES.DESC

------------------
#RES.COMM
resource "cloudfoundry_user" "my_user" {
  name    = "test"
  email       = "test@gmail.com"
  password    = "test123"
  given_name  = "test"
  family_name = "test"
}
##RES.COMM

--------------------
#RES.SAP
resource "cloudfoundry_user" "my_user" {
  username    = "test"
  email       = "test@gmail.com"
  password    = "test123"
  given_name  = "test"
  family_name = "test"
  annotations = { "purpose" : "testing", hi : "hello" }
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