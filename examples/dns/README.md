# DNS

## Importing existing resources

If you want to manage resources already created in the GleSYS Cloud.First create
a placeholder resource for the domain.

```
$ cat domain.tf
resource "glesys_dnsdomain" "mydomain" {
  name = "example.com"
}
```

```
$ terraform plan  # to verify that it looks correct
$ terraform import glesys_dnsdomain.mydomain 'example.com'

glesys_dnsdomain.mydomain: Importing from ID "example.com"...
glesys_dnsdomain.mydomain: Import prepared!
  Prepared glesys_dnsdomain for import
  glesys_dnsdomain.mydomain: Refreshing state... [id=example.com]

  Import successful!

  The resources that were imported are shown above. These resources are now in
  your Terraform state and will henceforth be managed by Terraform.
```

### Importing DNS records

Importing DNS records require the specific ID of the record.

The ID can be found by doing a manual request to the API

`curl -X POST --basic -u [API-USER]:[API-KEY] --data-urlencode "domainname=[data]" https://api.glesys.com/domain/listrecords`

To import the record

```
{'recordid': 12345678, 'domainname': 'example.com', 'host': 'mail', 'type': 'A', 'data': '172.16.0.1', 'ttl': 3600}
```

Prepare a placeholder resource for the record

```
cat domain.tf
...

resource "glesys_dnsdomain_record" "mail" {

}
```

```
terraform import glesys_dnsdomain_record.mail 'example.com,12345678'
```

The next run of 'terraform plan' will warn about missing parameters for the new
resource that has to be set. Add the missing (data, domain and host) parameters, and rerun 'terraform plan'.
