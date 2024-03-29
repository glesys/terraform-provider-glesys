---
page_title: "glesys_networkadapter Resource - terraform-provider-glesys"
subcategory: ""
description: |-
  Create a networkadapter attached to a VMware server.
---
# glesys_networkadapter (Resource)
Create a networkadapter attached to a VMware server.
## Example Usage
```terraform
resource "glesys_networkadapter" "netadapter2" {
  serverid = "wps123456"
  networkid = "vl123456"
  bandwidth = 200
}
```
<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `serverid` (String) Server ID to connect the adapter to

### Optional

- `adaptertype` (String) `VMXNET 3` (default) or `E1000`
- `bandwidth` (Number) adapter bandwidth
- `networkid` (String) Network ID to connect to. Defaults to `internet`.

### Read-Only

- `id` (String) The ID of this resource.
- `name` (String) Network Adapter name

