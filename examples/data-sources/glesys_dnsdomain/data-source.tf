# glesys_dnsdomain datasource
data "glesys_dnsdomain" "example" {
  name = "example.com"
}

output "domain_ttl" {
  value = data.glesys_dnsdomain.example.ttl
}
