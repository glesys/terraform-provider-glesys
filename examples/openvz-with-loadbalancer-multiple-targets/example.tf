#Instance count
variable "instance_count" {
  default = 2
}

# Set some variables 
variable "datacenter" {
  default = "Falkenberg"
}

# Instance defaults
variable "description" {
  default = "tf-test-norrland-"
}

variable "bw" {
  default = 100
  description = "Default bandwidth"
}

variable "mem" {
  default = 2048
  description = "Default memory"
}

variable "storage" {
  default = 20
  description = "Default storage"
}

variable "platform" {
  default = {
    Falkenberg = "OpenVZ"
    Stockholm = "VMware"
  }
}

variable "region" {
  default = {
    "0" = "Falkenberg"
    "1" = "Stockholm"
  }
}

variable "region_short" {
  default = {
    "0" = "fbg"
    "1" = "sth"
  }
}

variable "region_os" {
  default = {
    Falkenberg = "debian9"
    Stockholm = "debian9"
  }
}

variable "template" {
  default = {
    debian8 = "Debian 8 64-bit"
    debian9 = "Debian 9 64-bit"
  }
}

# Create the server resource
resource "glesys_server" "deb_64" {
  count = "${var.instance_count}"
  bandwidth = "${var.bw}"
  cpu = 2
  #datacenter = "${var.datacenter}"
  datacenter = "${lookup(var.region, count.index)}"
  description = "${var.description}deb Host terraform"
  #hostname = "${var.description}deb8-${count.index}" # hostname with count index only
  hostname = "${var.description}deb-${lookup(var.region_short, count.index)}-${count.index}" # hostname with short regionname and count
  memory = "${var.mem}"
  platform = "${lookup(var.platform,var.datacenter)}"
  storage = "${var.storage}"
  #template = "${lookup(var.template, lookup(var.region_os, lookup()var.datacenter))}"
  template = "${lookup(var.template, lookup(var.region_os, lookup(var.region, count.index) ) )}"
  password = "hunter2!"
}


resource "glesys_loadbalancer" "mylb" {
  count = 1
  datacenter = "Falkenberg"
  name = "mylb-1"
}

resource "glesys_loadbalancer_backend" "mybackend" {
  loadbalancerid = "${glesys_loadbalancer.mylb.id}"
  name = "my-web-backend"
  mode = "tcp"
  connecttimeout = 20000
}

resource "glesys_loadbalancer_frontend" "myfrontend" {
  loadbalancerid = "${glesys_loadbalancer.mylb.id}"
  name = "web-fe91"
  backend = "${glesys_loadbalancer_backend.mybackend.id}"
  port = 80
}

resource "glesys_loadbalancer_target" "webb" {
  loadbalancerid = "${glesys_loadbalancer.mylb.id}"
  backend = "${glesys_loadbalancer_backend.mybackend.id}"
  count = "${var.instance_count}"
  name = "target-${count.index}"
  port = 8898
  targetip = "${element(glesys_server.deb_64.*.ipv4_address, count.index)}"
  weight = 15
  enabled = false
}
