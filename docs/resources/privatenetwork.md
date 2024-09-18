---
page_title: "glesys_privatenetwork Resource - terraform-provider-glesys"
subcategory: ""
description: |-
  Create a PrivateNetwork resource.
---
# glesys_privatenetwork (Resource)
Create a PrivateNetwork resource.
## Example Usage
```terraform
resource "glesys_privatenetwork" "test" {
  name = "mynet"
}
```
<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) PrivateNetwork name

### Read-Only

- `id` (String) The ID of this resource.
- `ipv6aggregate` (String) IPv6Aggregate for the PrivateNetwork.
