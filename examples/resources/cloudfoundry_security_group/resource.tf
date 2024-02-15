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
