resource "glesys_dnsdomain" "mydomain" {
  name = "example.com"
}

resource "glesys_dnsdomain_record" "www" {
  domain = glesys_dnsdomain.mydomain.id
  data = "127.0.0.1"
  host = "www"
  type = "A"
}
