---
page_title: "Provider: Glesys"
description: |-
---

# GleSYS Provider

The GleSYS Provider is used to interact with the resources supported by GleSYS.
The provider needs to be configured with the proper credentials before it can be
used.

Use the navigation to the left to read about the available resources.

## Example Usage

Configure terraform: (example.tf)
```hcl
terraform {
  required_providers {
    glesys = {
      source = "glesys/glesys"
      version = "~> 0.4.4"
    }
  }
}

# Configure provider
provider "glesys" {
  token  = "ABC123"
  userid = "CL12345"
}

# Create a server
resource "glesys_server" "www" {
  # ...
}
```
Execute in your shell:
```sh
$ terraform init
$ terraform plan
```

## Authentication

### Static credentials

!> **Warning** Hard-coding credentials into any Terraform configuration is not
recommended, and risks secret leakage should this file ever be committed to a
public version control system.

Static credentials can be provided by adding `token` and `userid` in-line in the
Glesys provider block:

Usage:

```hcl
provider "glesys" {
  token  = "ABC123"
  userid = "CL12345"
}
```

### Environment variables

```hcl
provider "glesys" {}
```

Usage:

```sh
$ export GLESYS_TOKEN="ABC123"
$ export GLESYS_USERID="CL12345"
$ terraform plan
```

## Argument Reference

* `token` - (Optional) API Key for the Glesys API. Alternatively, this can be
  set using environment variables:
  * `GLESYS_TOKEN`
* `userid` (Optional) UserId for the Glesys API. The project name, such as
  CL12345. Can also be set using environment variables:
  * `GLESYS_USERID`
