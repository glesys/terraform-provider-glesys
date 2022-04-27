FROM golang:1.18.1-alpine as build-provider

COPY . /usr/src/myapp/

WORKDIR /usr/src/myapp

RUN go build -v -o terraform-provider-glesys_v0.3.2 -buildvcs=false

FROM hashicorp/terraform:1.1.9

RUN mkdir -p /root/.terraform.d/plugins/github.com/norrland/glesys/0.3.2/linux_amd64/

COPY --from=build-provider /usr/src/myapp/terraform-provider-glesys_v0.3.2 /root/.terraform.d/plugins/github.com/norrland/glesys/0.3.2/linux_amd64/

WORKDIR /home

ENTRYPOINT ["/bin/sh"]
