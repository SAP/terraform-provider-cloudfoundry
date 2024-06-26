---
page_title: "cloudfoundry_service_broker Resource - terraform-provider-cloudfoundry"
subcategory: ""
description: |-
  Provides a Cloud Foundry resource for managing service brokers
---

# cloudfoundry_service_broker (Resource)

Provides a Cloud Foundry resource for managing service brokers

## Example Usage

```terraform
resource "cloudfoundry_service_broker" "mysql" {
  name     = "broker"
  url      = "example.broker.com"
  username = "test"
  password = "test"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the service broker
- `password` (String, Sensitive) The password with which to authenticate against the service broker.
- `url` (String) URL of the service broker
- `username` (String) The username with which to authenticate against the service broker.

### Optional

- `annotations` (Map of String) The annotations associated with Cloud Foundry resources. Add as described [here](https://docs.cloudfoundry.org/adminguide/metadata.html#-view-metadata-for-an-object).
- `labels` (Map of String) The labels associated with Cloud Foundry resources. Add as described [here](https://docs.cloudfoundry.org/adminguide/metadata.html#-view-metadata-for-an-object).
- `space` (String) The GUID of the space the service broker is restricted to; omitted for globally available service brokers

### Read-Only

- `created_at` (String) The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.
- `id` (String) The GUID of the object.
- `updated_at` (String) The date and time when the resource was updated in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.

## Import

Import is supported using the following syntax:

```terraform
# terraform import cloudfoundry_service_broker.<resource_name> <service_broker_guid>

terraform import cloudfoundry_service_broker.service_broker 283f59d2-d660-45fb-9d96-b3e1aa92cfc7
```