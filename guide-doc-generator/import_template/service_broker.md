# Service Broker


-----------------
#RES.DESC

##RES.DESC

------------------
#RES.COMM
resource "cloudfoundry_service_broker" "mysql" { 
  name     = "broker"
  url      = "example.broker.com"
  username = "test"
  password = "test"
}
##RES.COMM

--------------------
#RES.SAP
resource "cloudfoundry_service_broker" "mysql" { 
  name     = "broker"
  url      = "example.broker.com"
  username = "test"
  password = "test"
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