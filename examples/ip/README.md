# IP

## Reserve a new adress

To allocate a new adress specify the required values `datacenter`, `platform` and `version`.

```
$ cat ip.tf
resource "glesys_ip" "example" {
  datacenter = "Stockholm"
  platform   = "KVM"
  version    = 6
}
```

Changing any of these values will release the current address and reserve a new one.

## Importing already reserved addresses

Addresses that are already reserved can be imported into terraform.

Prepare a resource for the address.

```
$ cat ip.tf
resource "glesys_ip" "example" {
  address = "1.2.3.4"
}
```

Import the resource by specifying its address.

```
$ terraform import glesys_ip.example 1.2.3.4
glesys_ip.example: Importing from ID "1.2.3.4"...
glesys_ip.example: Import prepared!
  Prepared glesys_ip for import
glesys_ip.example: Refreshing state... [id=1.2.3.4]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

## Changing reverse pointer records

To manage reverse pointer records use the `ptr` attribute on the resource.
Note that the `ptr` attribute *must* end with a dot `.` character.

```
$ cat ip.tf
resource "glesys_ip" "example" {
  datacenter = "Stockholm"
  platform   = "KVM"
  version    = 6
  ptr        = "my.ptr."
}
```

Due to how terraform internals work it is not possible to reset a reverse pointer by removing the attribute from the terraform file. Possible workarounds for resetting a reverse pointer value:

1. Set the `ptr` attribute in terraform to what the default value would be. For example, for ip `1.2.3.4` the default value would be `1-2-3-4-static.glesys.net.`.
1. Reset the reverse pointer value manually via the web interface or the [API](https://github.com/GleSYS/API/wiki/API-Documentation#ipresetptr). After resetting the reverse pointer manually, run `terraform refresh` to reflect the change in the terraform state.
