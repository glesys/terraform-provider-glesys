# Vagrant

Run terraform-provider-glesys in a virtual machine!

## Getting started

https://learn.hashicorp.com/collections/vagrant/getting-started

### tl;dr

Assuming you have vagrant installed, and a decent hypervisor available.

```
$ vagrant up
$ vagrant ssh
$ cd code/terraform-provider-glesys/examples/vmware
$ terraform init
$ GLESYS_USERID=CL12345 GLESYS_TOKEN=abcsecrettoken111 terraform apply
```
