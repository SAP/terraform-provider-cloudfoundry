# Service binding


-----------------
#RES.DESC

##RES.DESC

------------------
#RES.COMM
resource "cloudfoundry_service_key" "redis1-key1" {
  name = "hifi"
  service_instance = "e9ec29ca-993d-42e2-9c5b-cb17b1972cce"
}
##RES.COMM

--------------------
#RES.SAP
resource "cloudfoundry_service_credential_binding" "scb1" {
  type             = "key"
  name             = "hifi"
  service_instance = "e9ec29ca-993d-42e2-9c5b-cb17b1972cce"
}

resource "cloudfoundry_service_credential_binding" "scb7" {
  type             = "app"
  name             = "hifi"
  service_instance = "e9ec29ca-993d-42e2-9c5b-cb17b1972cce"
  app              = "ec6ac2b3-fb79-43c4-9734-000d4299bd59"
}
##RES.SAP

---------------

#DS.DESC
##DS.DESC
----------------

#DS.COMM
data "cloudfoundry_service_key" "my-key" {
    name             = "my-service-key"
    service_instance = "e9ec29ca-993d-42e2-9c5b-cb17b1972cce"
}
##DS.COMM
-----------------

#DS.SAP
data "cloudfoundry_service_credential_binding" "scbd" {
  service_instance = "5e2976bb-332e-41e1-8be3-53baafea9296"
  app              = "ec6ac2b3-fb79-43c4-9734-000d4299bd59"
}

data "cloudfoundry_service_credential_binding" "my-key" {
    name             = "my-service-key"
    service_instance = "e9ec29ca-993d-42e2-9c5b-cb17b1972cce"
}
##DS.SAP