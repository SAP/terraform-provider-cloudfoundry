# Buildpack


-----------------
#RES.DESC

##RES.DESC

------------------
#RES.COMM
resource "cloudfoundry_buildpack" "bp" {
  name     = "hi"
  position = 1
  enabled  = false
  locked   = true
  labels   = { "hi" : "fi" }
  path     = "something.zip"
}
##RES.COMM

--------------------
#RES.SAP
resource "cloudfoundry_buildpack" "bp" {
  name     = "hi"
  position = 1
  stack    = "cflinuxfs3"
  enabled  = false
  locked   = true
  labels   = { "hi" : "fi" }
  path     = "something.zip"
}
##RES.SAP

---------------

#DS.DESC

##DS.DESC
----------------

#DS.COMM
##DS.COMM
-----------------

#DS.SAP
##DS.SAP