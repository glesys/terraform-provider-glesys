[![Build Status](https://travis-ci.org/norrland/terraform-provider-glesys.svg?branch=master)](https://travis-ci.org/norrland/terraform-provider-glesys)

# terraform-provider-glesys

## intro

This is a PoC for working with terraform and glesys/glesys-go.
Please see https://github.com/glesys/glesys-go and https://github.com/GleSYS/API/wiki/API-Documentation for more information.

## Installation / Local development

Clone the repo into a folder of your liking.

### Link the terraform-provider-glesys to your GOPATH.

`ln -s terraform-provider-glesys $GOPATH/src/github.com/norrland/terraform-provider-glesys`

### Fetch the dependencies

`cd terraform-provider-glesys && go get -d`

### build and use the provider

`terraform-provider-glesys$ go build -o terraform-provider-glesys . `

### Run terraform

Instead of hardcoding credentials into your terraform templates.
Use environment variables for example.

`GLESYS_USERID="CL12345" GLESYS_TOKEN="ABC12345678" terraform plan`
