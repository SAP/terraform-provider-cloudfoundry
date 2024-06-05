# Routes


-----------------
#RES.DESC
Domain ID and Space ID are used instead of instead of domain and space in the new provider. The new resource uses `destinations` instead og `target` to specifiy route mappings.
##RES.DESC

------------------
#RES.COMM
resource "cloudfoundry_route" "default" {
    domain = data.cloudfoundry_domain.apps.domain.id
    space = data.cloudfoundry_space.dev.id
    hostname = "myapp"
}
##RES.COMM

--------------------
#RES.SAP

resource "cloudfoundry_route" "bruh" {
  space  = "795a961c-6360-479a-9666-fff9cb906aad"
  domain = "440e24e5-ee11-41d9-a996-2ed0a1e2deea"

  #Optional 

  host   = "myapp"
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

##RES.SAP

---------------

#DS.DESC
The new provider also supports filtereing routes on the basis of `space` in additions to `org`.
##DS.DESC
----------------

#DS.COMM

data "cloudfoundry_route" "my-route" {
    domain   = "domain-id"
    hostname = "my-host"
}
##DS.COMM
-----------------

#DS.SAP
data "cloudfoundry_route" "bruh" {
  domain = "a25ca0c1-353a-40f9-bcf4-d2a0adf4112b"
  space - "b45da1f2-353a-40f9-bcf4-d2a0adf4112b"
}
##DS.SAP