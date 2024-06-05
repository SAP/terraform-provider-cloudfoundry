# Routes


## Resource

Domain ID and Space ID are used instead of instead of domain and space in the new provider. The new resource uses `destinations` instead og `target` to specifiy route mappings.


| Community Cloudfoundry Provider | SAP Cloudfoundry Provider |
| -- | -- |
| <pre>resource "cloudfoundry_route" "default" {</br>    domain = data.cloudfoundry_domain.apps.domain.id</br>    space = data.cloudfoundry_space.dev.id</br>    hostname = "myapp"</br>}</br></pre> | <pre></br>resource "cloudfoundry_route" "bruh" {</br>  space  = "795a961c-6360-479a-9666-fff9cb906aad"</br>  domain = "440e24e5-ee11-41d9-a996-2ed0a1e2deea"</br></br>  #Optional </br></br>  host   = "myapp"</br>  path   = "/cart"</br>  destinations = [</br>    {</br>      app_id = "24a711f2-148b-4d48-b37a-90a66d6e842f"</br></br>    },</br>    {</br>      app_id           = "15a74002-cf3a-4bf2-b76f-fe96867c46ee"</br>      app_process_type = "web"</br>      port             = 36001</br>    },</br></br>  ]</br>}</br></br></pre> |

<br/>

## Datasource


The new provider also supports filtereing routes on the basis of `space` in additions to `org`.


| Community Cloudfoundry Provider | SAP Cloudfoundry Provider |
| -- | -- |
| <pre></br>data "cloudfoundry_route" "my-route" {</br>    domain   = "domain-id"</br>    hostname = "my-host"</br>}</br></pre> | <pre>data "cloudfoundry_route" "bruh" {</br>  domain = "a25ca0c1-353a-40f9-bcf4-d2a0adf4112b"</br>  space - "b45da1f2-353a-40f9-bcf4-d2a0adf4112b"</br>}</br></pre> |
