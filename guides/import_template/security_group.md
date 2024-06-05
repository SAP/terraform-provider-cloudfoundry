# Application Security Group


-----------------
#RES.DESC
Application Security Group `asg`  is now identified as `security_group`. The newer resource also exposes some additional parameters that are introduced as part of the v3 specification.
##RES.DESC

------------------
#RES.COMM
resource "cloudfoundry_asg" "my_security_group" {
  name = "rmq-service"

  rule {
    protocol = "tcp"
    destination = "192.168.1.100"
    ports = "1883,8883"
    log = true
  }
  rule {
    protocol = "tcp"
    destination = "192.168.1.101"
    ports = "5671-5672"
    log = true
  }
}
##RES.COMM

--------------------
#RES.SAP
resource "cloudfoundry_security_group" "my_security_group" {
  name                     = "tf-test"
  globally_enabled_running = false
  globally_enabled_staging = false
  rules = [{
    protocol    = "tcp"
    destination = "192.168.1.100"
    ports       = "1883,8883"
    log         = true
    }, {
    protocol    = "udp"
    destination = "192.168.1.100"
    ports       = "1883,8883"
    log         = false
    },
    {
      protocol    = "icmp"
      type        = 0
      code        = 0
      destination = "192.168.1.100"
      log         = false
  }]
  staging_spaces = ["3bc20dc4-1870-4835-8308-dda2d766e61e", "e6886bba-e263-4b52-aaf1-85d410f15fc8"]
  running_spaces = ["e6886bba-e263-4b52-aaf1-85d410f15fc8"]

}
##RES.SAP