# Service instance


-----------------
#RES.DESC

##RES.DESC

------------------
#RES.COMM
resource "cloudfoundry_service_instance" "redis1" {
  name         = "pricing-grid"
  space        = "e6886bba-e263-4b52-aaf1-85d410f15fc8"
  tags         = ["terraform-test", "test1"]
  service_plan = data.cloudfoundry_service.redis.service_plans["shared-vm"]
}

resource "cloudfoundry_user_provided_service" "mq" {
  name = "mq-server"
  space = "e6886bba-e263-4b52-aaf1-85d410f15fc8"
  credentials = {
    "url" = "mq://localhost:9000"
    "username" = "admin"
    "password" = "admin"
  }
}
##RES.COMM

--------------------
#RES.SAP
resource "cloudfoundry_service_instance" "redis1" {
  name         = "pricing-grid"
  type         = "managed"
  space        = "e6886bba-e263-4b52-aaf1-85d410f15fc8"
  tags         = ["terraform-test", "test1"]
  service_plan = data.cloudfoundry_service.redis.service_plans["shared-vm"]

resource "cloudfoundry_service_instance" "mq" {
  name        = "mq"
  type        = "user-provided"
  space       = "e6886bba-e263-4b52-aaf1-85d410f15fc8"
  credentials = <<EOT
  {
    "url" = "mq://localhost:9000"
    "username" = "admin"
    "password" = "admin"
  }
  EOT
}
##RES.SAP

---------------

#DS.DESC
##DS.DESC
----------------

#DS.COMM
data "cloudfoundry_user_provided_service" "my-instance" {
    name  = "my-service-instance"
    space = "02c0cc92-6ecc-44b1-b7b2-096ca19ee143"
}

data "cloudfoundry_service_instance" "svc" {
    name_or_id = "managed-service-instance"
    space      = "02c0cc92-6ecc-44b1-b7b2-096ca19ee143"
}
##DS.COMM
-----------------

#DS.SAP
data "cloudfoundry_service_instance" "my-instance" {
  name  = "my-service-instance"
  space = "02c0cc92-6ecc-44b1-b7b2-096ca19ee143"
}

data "cloudfoundry_service_instance" "svc" {
  name  = "managed-service-instance"
  space = "02c0cc92-6ecc-44b1-b7b2-096ca19ee143"
}
##DS.SAP