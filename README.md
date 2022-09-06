# terraform-provider-glesys

## intro

This is a early stage implementation for using Terraform and the GleSYS API.
Please see https://github.com/glesys/glesys-go and https://github.com/GleSYS/API/wiki/API-Documentation for more information.

## Installation

### Using the terraform-provider-glesys

In your Terraform configuration, add this. Then run `terraform init`

```
terraform {
  required_providers {
    glesys = {
      source = "glesys/glesys"
      # version = "0.4.2" # If you want to specify a certain release.
    }
  }
}

provider "glesys" {
  # Configuration options
}
```

`$ terraform init`

### Run terraform

Instead of hardcoding credentials into your terraform templates.
Use environment variables for example.

`GLESYS_USERID="CL12345" GLESYS_TOKEN="ABC12345678" terraform plan`

## Local development
### Debian requirements

- golang >= 1.15
- git
- make
- terraform # Fetch the latest version from https://www.terraform.io/downloads

### Setup terraform-provider-glesys for local development

Clone the repo into a folder of your liking.

`$ git clone https://github.com/glesys/terraform-provider-glesys.git`

### Install the plugin

```
$ cd terraform-provider-glesys
$ make
==> Checking that code complies with gofmt requirements...
go install
go: finding github.com/hashicorp/terraform v0.12.9
...
$ mkdir -p ~/.terraform.d/local_dev
$ ln -s ~/go/bin/terraform-provider-glesys ~/.terraform.d/local_dev/
```

Setup dev_overrides

Create a file `provider_overrides.tfrc` for example:
```
provider_installation {
  dev_overrides {
    "glesys/glesys" = "/home/myuser/.terraform.d/local_dev"
  }
  direct {}
}
```

Run terraform with environment variable `TF_CLI_CONFIG_FILE`

`$ TF_CLI_CONFIG_FILE=/path/to/provider_overrides.tfrc terraform plan`

## Contribute

#### We love Pull Requests ♥

1. Fork the repo.
2. Make sure to run the tests to verify that you're starting with a clean slate.
3. Add a test for your change, make sure it fails. Refactoring existing code or
   improving documentation does not require new tests.
4. Make the changes and ensure the test pass.
5. Commit your changes, push to your fork and submit a Pull Request.
