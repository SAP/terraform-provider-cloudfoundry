# Service instance


## Resource




|  SAP Cloudfoundry Provider |Community Cloudfoundry Provider |
| -- | -- |
|  <pre>resource "cloudfoundry_service_instance" "redis1" {</br>  name         = "pricing-grid"</br>  type         = "managed"</br>  space        = "e6886bba-e263-4b52-aaf1-85d410f15fc8"</br>  tags         = ["terraform-test", "test1"]</br>  service_plan = data.cloudfoundry_service.redis.service_plans["shared-vm"]</br></br>resource "cloudfoundry_service_instance" "mq" {</br>  name        = "mq"</br>  type        = "user-provided"</br>  space       = "e6886bba-e263-4b52-aaf1-85d410f15fc8"</br>  credentials = <<EOT</br>  {</br>    "url" = "mq://localhost:9000"</br>    "username" = "admin"</br>    "password" = "admin"</br>  }</br>  EOT</br>}</br></pre> |<pre>resource "cloudfoundry_service_instance" "redis1" {</br>  name         = "pricing-grid"</br>  space        = "e6886bba-e263-4b52-aaf1-85d410f15fc8"</br>  tags         = ["terraform-test", "test1"]</br>  service_plan = data.cloudfoundry_service.redis.service_plans["shared-vm"]</br>}</br></br>resource "cloudfoundry_user_provided_service" "mq" {</br>  name = "mq-server"</br>  space = "e6886bba-e263-4b52-aaf1-85d410f15fc8"</br>  credentials = {</br>    "url" = "mq://localhost:9000"</br>    "username" = "admin"</br>    "password" = "admin"</br>  }</br>}</br></pre> |

<br/>

## Datasource




|  SAP Cloudfoundry Provider | Community Cloudfoundry Provider  |
| -- | -- |
| <pre>data "cloudfoundry_service_instance" "my-instance" {</br>  name  = "my-service-instance"</br>  space = "02c0cc92-6ecc-44b1-b7b2-096ca19ee143"</br>}</br></br>data "cloudfoundry_service_instance" "svc" {</br>  name  = "managed-service-instance"</br>  space = "02c0cc92-6ecc-44b1-b7b2-096ca19ee143"</br>}</br></pre>|<pre>data "cloudfoundry_user_provided_service" "my-instance" {</br>    name  = "my-service-instance"</br>    space = "02c0cc92-6ecc-44b1-b7b2-096ca19ee143"</br>}</br></br>data "cloudfoundry_service_instance" "svc" {</br>    name_or_id = "managed-service-instance"</br>    space      = "02c0cc92-6ecc-44b1-b7b2-096ca19ee143"</br>}</br></pre> |  
