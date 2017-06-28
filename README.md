# terraform-provider-glesys

## intro

Use at your own risk. This is a PoC for working with terraform and glesys/glesys-go.

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
