# Create the server resource
resource "glesys_server" "vm" {
  bandwidth   = 100
  cpu         = 1
  datacenter  = "Falkenberg"
  description = "Server with extra disks"
  hostname    = "vmware-vm-fbg1-tf-extra-disk"
  memory      = 1024
  platform    = "VMware"
  storage     = 20
  template    = "Debian 12 64-bit"

  user {
    username = "alice"
    publickeys = [
      "ssh-ed25519 AAAAAAmykeyFFFFFF alice@example.com"
    ]
    password = "hunter4!"
  }
}

resource "glesys_server_disk" "data" {
  serverid = glesys_server.vm.id

  size = 100
  name = "Data disk"
}
