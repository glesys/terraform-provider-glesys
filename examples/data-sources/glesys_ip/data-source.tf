# glesys_dnsdomain datasource
data "glesys_ip" "www-pub" {
  id = "192.0.2.10"
}

output "www-ptr" {
  value = data.glesys_ip.www-pub.ptr
}
