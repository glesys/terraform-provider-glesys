FROM golang:1.14-alpine as build-provider

COPY . /usr/src/myapp/

WORKDIR /usr/src/myapp

RUN go build -v -o terraform-provider-glesys_v0.4.0

FROM hashicorp/terraform:light

RUN mkdir -p /root/.terraform.d/plugins/github.com/norrland/glesys/0.4.0/linux_amd64/

COPY --from=build-provider /usr/src/myapp/terraform-provider-glesys_v0.4.0 /root/.terraform.d/plugins/github.com/norrland/glesys/0.4.0/linux_amd64/

WORKDIR /home

ENTRYPOINT ["/bin/sh"]
