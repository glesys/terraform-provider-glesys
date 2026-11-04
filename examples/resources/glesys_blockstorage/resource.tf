resource "glesys_blockstorage" "datadisk" {
  name       = "extra storage"
  datacenter = "dc-fbg1"
  size       = 200
  tier       = "velocity"
  serverid   = "kvm123456"
}
