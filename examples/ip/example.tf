resource "glesys_ip" "example" {
  datacenter = "Stockholm"
  platform   = "KVM"
  version    = 4
}

resource "glesys_ip" "ptr_example" {
  address = "1.2.3.4"
  ptr     = "example.com."
}
