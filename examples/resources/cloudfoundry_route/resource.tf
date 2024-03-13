resource "cloudfoundry_route" "bruh" {
  space  = "795a961c-6360-479a-9666-fff9cb906aad"
  domain = "440e24e5-ee11-41d9-a996-2ed0a1e2deea"
  host   = "testing123"
  path   = "/cart"
  destinations = [
    {
      app_id = "24a711f2-148b-4d48-b37a-90a66d6e842f"

    },
    {
      app_id           = "15a74002-cf3a-4bf2-b76f-fe96867c46ee"
      app_process_type = "web"
      port             = 36001
    },

  ]
} 