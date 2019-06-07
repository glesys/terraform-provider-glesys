# Set some variables 
variable "datacenter" {
  default = "Falkenberg"
}

# Instance defaults
variable "description" {
  default = "tf-test-norrland-"
}

variable "bw" {
  default     = 100
  description = "Default bandwidth"
}

variable "mem" {
  default     = 2048
  description = "Default memory"
}

variable "storage" {
  default     = 20
  description = "Default storage"
}

variable "platform" {
  default = {
    Falkenberg = "OpenVZ"
    Stockholm  = "VMware"
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
    Stockholm  = "debian9"
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
  count     = 1
  bandwidth = var.bw
  cpu       = 2

  #datacenter = "${var.datacenter}"
  datacenter  = var.region[count.index]
  description = "${var.description}deb Host terraform"

  #hostname = "${var.description}deb8-${count.index}" # hostname with count index only
  hostname = "${var.description}deb-${var.region_short[count.index]}-${count.index}" # hostname with short regionname and count
  memory   = var.mem
  platform = var.platform[var.datacenter]
  storage  = var.storage

  #template = "${lookup(var.template, lookup(var.region_os, lookup()var.datacenter))}"
  template = var.template[var.region_os[var.region[count.index]]]
  password = "2hunter2!"
}

