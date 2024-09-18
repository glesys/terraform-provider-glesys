### Basic segment

resource "glesys_privatenetwork_segment" "testseg" {
  privatenetworkid = "pn-123ab"
  datacenter       = "dc-fbg1"
  ipv4subnet       = "10.0.0.0/24"
  name             = "seg-1"
  platform         = "kvm"
}

### PrivateNetwork, Segment And NetworkAdapter Example

resource "glesys_privatenetwork" "privatenet" {
  name = "vm-internal"
}

// Segment attached to 'privatenet'
resource "glesys_privatenetwork_segment" "seg-kvm-fbg" {
  privatenetworkid = glesys_privatenetwork.privatenet.id

  name       = "int-fbg"
  datacenter = "dc-fbg1"
  ipv4subnet = "10.2.0.0/24"
  platform   = "kvm"
}

resource "glesys_server" "myvm" {
  platform   = "KVM"
  datacenter = "Falkenberg"
  # ...
}

// NetworkAdapter attached to KVM vm 'myvm'
resource "glesys_networkadapter" "kvm-nic2" {
  serverid = glesys_server.myvm.id

  name      = "pn-internal"
  bandwidth = 1000
  networkid = glesys_privatenetwork_segment.seg-kvm-fbg.id

}
