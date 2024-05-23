# space Quota


-----------------
#RES.DESC

##RES.DESC

------------------
#RES.COMM
resource "cloudfoundry_space_quota" "large" {
    name                     = "large"
    org                      = "ca721b24-e24d-4171-83e1-1ef6bd836b38"
    allow_paid_service_plans = false
    total_memory             = 51200
    total_routes             = 50
    total_services           = 200
          
    #Optionals

    instance_memory          = 2048
    total_service_keys       = 120
    total_app_instances      = 100
    total_route_ports        = 5
    total_app_tasks          = 10
}
##RES.COMM

--------------------
#RES.SAP

resource "cloudfoundry_space_quota" "large" {
  name                     = "large"
  allow_paid_service_plans = false
  org                      = "ca721b24-e24d-4171-83e1-1ef6bd836b38"
  
  #Optionals
  
  total_memory             = 51200
  total_routes             = 50
  total_services           = 200
  instance_memory          = 2048
  total_service_keys       = 120
  total_app_instances      = 100
  total_route_ports        = 5
  total_app_tasks          = 10

  total_app_log_rate_limit = 1000
  spaces = [
    "ba10cc63-cc43-46b1-a00c-5f2a0d7d992e",
  ]
}
##RES.SAP

---------------

#DS.DESC
##DS.DESC
----------------

#DS.COMM
data "cloudfoundry_space_quota" "my_space_quota" {
  name = "tf-test-do-not-delete"
  org  = "ca721b24-e24d-4171-83e1-1ef6bd836b38"
}
##DS.COMM
-----------------

#DS.SAP
data "cloudfoundry_space_quota" "my_space_quota" {
  name = "tf-test-do-not-delete"
  org  = "ca721b24-e24d-4171-83e1-1ef6bd836b38"
}
##DS.SAP