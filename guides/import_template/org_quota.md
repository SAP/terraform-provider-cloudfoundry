# Org Quota


-----------------
#RES.DESC
In the v3 API only `allow_paid_service_plans` and `name` are required fields. Orgs can be bound in the new provider.

##RES.DESC

------------------
#RES.COMM
resource "cloudfoundry_org_quota" "large" {
    name = "large"
    allow_paid_service_plans = false
    instance_memory = 2048
    total_memory = 51200
    total_app_instances = 100
    total_routes = 50
    total_services = 200
    total_route_ports = 5
}

##RES.COMM

--------------------
#RES.SAP

resource "cloudfoundry_org_quota" "large" {
  name                     = "large"
  allow_paid_service_plans = false

  #Optionals
  
  instance_memory          = 2048
  total_memory             = 51200
  total_app_instances      = 100
  total_routes             = 50
  total_services           = 200
  total_service_keys       = 120
  total_private_domains    = 40
  total_app_tasks          = 10
  total_route_ports        = 5
  total_app_log_rate_limit = 1000
  orgs = [
    "ba10cc63-cc43-46b1-a00c-5f2a0d7d992e",
  ]
}

##RES.SAP