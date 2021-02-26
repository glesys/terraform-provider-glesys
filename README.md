# terraform-provider-glesys

## intro

This is a early stage implementation for using Terraform and the GleSYS API.
Please see https://github.com/glesys/glesys-go and https://github.com/GleSYS/API/wiki/API-Documentation for more information.

## Installation / Local development

### Debian requirements

- golang >= 1.14
- git
- make
- terraform 0.12.x # Fetch the latest 0.12 version from https://releases.hashicorp.com/terraform/

#### install terraform

```
$ curl -O https://releases.hashicorp.com/terraform/0.12.29/terraform_0.12.29_linux_amd64.zip
$ unzip terraform_0.12.29_linux_amd64.zip
$ mv terraform ~/bin
```

### Setup terraform-provider-glesys

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
$ mkdir -p ~/.terraform.d/plugins
$ ln -s ~s ~/go/bin/terraform-provider-glesys ~/.terraform.d/plugins/
```

### Run terraform

Instead of hardcoding credentials into your terraform templates.
Use environment variables for example.

`GLESYS_USERID="CL12345" GLESYS_TOKEN="ABC12345678" terraform plan`
