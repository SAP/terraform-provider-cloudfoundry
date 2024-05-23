# Service


-----------------
#RES.DESC

##RES.DESC

------------------
#RES.COMM

##RES.COMM

--------------------
#RES.SAP

##RES.SAP

---------------

#DS.DESC
##DS.DESC
----------------

#DS.COMM
data "cloudfoundry_service" "redis" {
    name = "p-redis"    
}
##DS.COMM
-----------------

#DS.SAP
data "cloudfoundry_service" "redis" {
  name = "p-redis"
}
##DS.SAP