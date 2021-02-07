resource "glesys_server" "kvm" {
  count = 1
  datacenter = "Stockholm"
  memory = 1024
  storage = 20
  cpu = 1
  bandwidth = 100

  hostname = "www1"

  platform = "KVM"
  template = "debian-10"

  user {
        username = "alice"
        publickeys = [
          "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINOCh8br7CwZDMGmINyJgBip943QXgkf7XdXrDMJf5Dl alice@example.com",
          "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOfN4dBsS2p1UX+DP6RicdxAYCCeRK8mzCldCS0W9A+5 alice@ws.example.com"
        ]
        password = "hunter3!"
  }
  user {
        username = "bob"
        publickeys = ["ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINOCh8br7CwZDMGmINyJgBip943QXgkf7XdXrDMJf5Dl bob@example.com"]
        password = "hunter333!"
  }
}
