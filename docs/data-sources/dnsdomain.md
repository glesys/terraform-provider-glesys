---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "glesys_dnsdomain Data Source - terraform-provider-glesys"
subcategory: ""
description: |-
  Get information about a DNS Domain associated with your GleSYS Project.
---

# glesys_dnsdomain (Data Source)

Get information about a DNS Domain associated with your GleSYS Project.

## Example Usage

```terraform
# glesys_dnsdomain datasource
data "glesys_dnsdomain" "example" {
  name = "example.com"
}

output "domain_ttl" {
  value = data.glesys_dnsdomain.example.ttl
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) name of the domain

### Read-Only

- `expire` (Number) expire ttl of the domain.
- `id` (String) The ID of this resource.
- `minimum` (Number) minimum ttl of the domain.
- `refresh` (Number) refresh ttl of the domain.
- `retry` (Number) retry ttl of the domain.
- `ttl` (Number) ttl of the domain.
