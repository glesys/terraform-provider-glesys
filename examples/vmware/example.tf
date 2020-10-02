# Set some variables
variable "datacenter" {
  default = "Falkenberg"
}

# Instance defaults
variable "description" {
  default = "tf-vmware-"
}

variable "bw" {
  default     = 100
  description = "Default bandwidth"
}

variable "mem" {
  default     = 1024
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

variable "templatelist" {
  default = {
    "0" = "Centos 7 64-bit"
    "1" = "Centos 8 64-bit"
    "2" = "Debian 8 64-bit"
    "3" = "Debian 9 64-bit"
    "4" = "Debian 10 64-bit"
    "5" = "Ubuntu 14.04 LTS 64-bit"
    "6" = "Ubuntu 16.04 LTS 64-bit"
    "7" = "Ubuntu 18.04 LTS 64-bit"
    "8" = "Ubuntu 20.04 LTS 64-bit"
  }
}

# Create the server resource
resource "glesys_server" "vmware" {
  count       = 1
  bandwidth   = var.bw
  cpu         = 1
  datacenter  = "Falkenberg"
  description = "${var.description}-example-${count.index}"
  hostname    = "vmware-vm-${count.index}"
  memory      = var.mem
  platform    = "VMware"
  storage     = var.storage
  template    = var.templatelist[count.index]
  password    = "hunter2!"
}

