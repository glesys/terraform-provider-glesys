resource "glesys_network" "network-example-ams" {
  count = 1
  datacenter = "Falkenberg"
  description = "tf-test-fbg-${count.index}"
}
resource "glesys_network" "network-example-ams" {
  count = 1
  datacenter = "Amsterdam"
  description = "tf-test-ams-${count.index}"
}

