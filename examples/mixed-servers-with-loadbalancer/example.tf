# This example sets up KVM:s in Sweden, Norway and Finland. It also creates a round-robin load balancer for the KVM's



# Server: Sweden 1 - Falkenberg

resource "glesys_server" "sweden1" {
  count = 1
  datacenter = "Falkenberg"
  memory = 1024
  storage = 20
  cpu = 1
  bandwidth = 100

  hostname = "sweden1"

  platform = "KVM"
  template = "debian-10"

  user {
        username = "irken"
        publickeys = [
          "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINOCh8br7CwZDMGmINyJgBip943QXgkf7XdXrDMJf5Dl irken@example.com",
        ]
        password = "S01.E09.A.Room.With.A.Moose"
  }
}

# Server: Sweden 2 - Stockholm

resource "glesys_server" "sweden2" {
  count = 1
  datacenter = "Stockholm"
  memory = 1024
  storage = 20
  cpu = 1
  bandwidth = 100

  hostname = "sweden2"

  platform = "KVM"
  template = "debian-10"

  user {
    username = "irken"
    publickeys = [
      "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINOCh8br7CwZDMGmINyJgBip943QXgkf7XdXrDMJf5Dl irken@example.com",
    ]
    password = "S01.E09.A.Room.With.A.Moose"
  }
}

# Server: Finland 1 - Oulu

resource "glesys_server" "finland1" {
  count = 1
  datacenter = "Oulu"
  memory = 1024
  storage = 20
  cpu = 1
  bandwidth = 100

  hostname = "finland1"

  platform = "KVM"
  template = "debian-10"

  user {
    username = "irken"
    publickeys = [
      "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINOCh8br7CwZDMGmINyJgBip943QXgkf7XdXrDMJf5Dl irken@example.com",
    ]
    password = "S01.E09.A.Room.With.A.Moose"
  }
}

# Server: Norway 1 - Oslo

resource "glesys_server" "norway1" {
  count = 1
  datacenter = "Oslo"
  memory = 1024
  storage = 20
  cpu = 1
  bandwidth = 100
  password="23m3iour2h323nu3h2e!"

  hostname = "norway1"

  platform = "VMware"
  template = "Debian 11 64-bit"

  user {
    username = "irken"
    publickeys = [
      "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINOCh8br7CwZDMGmINyJgBip943QXgkf7XdXrDMJf5Dl irken@example.com",
    ]
    password = "S01.E09.A.Room.With.A.Moose"
  }
}

# Server: Netherlands 1 - Amsterdam

resource "glesys_server" "netherlands1" {
  count = 1
  datacenter = "Amsterdam"
  memory = 1024
  storage = 20
  cpu = 1
  bandwidth = 100
  password="23m3iour2h323nu3h2e!"

  hostname = "netherlands1"

  platform = "VMware"
  template = "Debian 11 64-bit"

  user {
    username = "irken"
    publickeys = [
      "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINOCh8br7CwZDMGmINyJgBip943QXgkf7XdXrDMJf5Dl irken@example.com",
    ]
    password = "S01.E09.A.Room.With.A.Moose"
  }
}

# Loadbalancer: eu-north

resource "glesys_loadbalancer" "eu-north" {
  datacenter = "Falkenberg"
  name       = "eu-north"
}

resource "glesys_loadbalancer_backend" "eu-north-be" {
  depends_on = [
    glesys_loadbalancer.eu-north
  ]
  loadbalancerid = glesys_loadbalancer.eu-north.id
  name           = "eu-north-be"
  mode           = "tcp"
  connecttimeout = 20000
}


resource "glesys_loadbalancer_frontend" "eu-north-fe" {
  depends_on = [
    glesys_loadbalancer.eu-north,
    glesys_loadbalancer_backend.eu-north-be
  ]
  loadbalancerid = glesys_loadbalancer.eu-north.id
  name           = "eu-north-fe"
  backend        = glesys_loadbalancer_backend.eu-north-be.id
  port           = 80
}



resource "glesys_loadbalancer_target" "sweden1" {
  depends_on = [
    glesys_server.sweden1,
    glesys_loadbalancer.eu-north,
    glesys_loadbalancer_backend.eu-north-be,
    glesys_loadbalancer_frontend.eu-north-fe
  ]
  loadbalancerid = glesys_loadbalancer.eu-north.id
  backend        = glesys_loadbalancer_backend.eu-north-be.id
  name     = "sweden1"
  port     = 8080
  targetip = element(glesys_server.sweden1.*.ipv4_address, 0)
  weight   = 15
  enabled = true
}

resource "glesys_loadbalancer_target" "sweden2" {
  depends_on = [
    glesys_server.sweden2,
    glesys_loadbalancer.eu-north,
    glesys_loadbalancer_backend.eu-north-be,
    glesys_loadbalancer_frontend.eu-north-fe
  ]
  loadbalancerid = glesys_loadbalancer.eu-north.id
  backend        = glesys_loadbalancer_backend.eu-north-be.id
  name     = "sweden2"
  port     = 8080
  targetip = element(glesys_server.sweden2.*.ipv4_address, 0)
  weight   = 15
  enabled = true
}


resource "glesys_loadbalancer_target" "norway1" {
  depends_on = [
    glesys_server.norway1,
    glesys_loadbalancer.eu-north,
    glesys_loadbalancer_backend.eu-north-be,
    glesys_loadbalancer_frontend.eu-north-fe
  ]
  loadbalancerid = glesys_loadbalancer.eu-north.id
  backend        = glesys_loadbalancer_backend.eu-north-be.id
  name     = "norway1"
  port     = 8080
  targetip = element(glesys_server.norway1.*.ipv4_address, 0)
  weight   = 15
  enabled = true
}


resource "glesys_loadbalancer_target" "finland1" {
  depends_on = [
    glesys_server.finland1,
    glesys_loadbalancer.eu-north,
    glesys_loadbalancer_backend.eu-north-be,
    glesys_loadbalancer_frontend.eu-north-fe
  ]
  loadbalancerid = glesys_loadbalancer.eu-north.id
  backend        = glesys_loadbalancer_backend.eu-north-be.id
  name     = "finland1"
  port     = 8080
  weight   = 15
  targetip = element(glesys_server.finland1.*.ipv4_address, 0)
  enabled = true
}

resource "glesys_loadbalancer_target" "netherlands1" {
  depends_on = [
    glesys_server.netherlands1,
    glesys_loadbalancer.eu-north,
    glesys_loadbalancer_backend.eu-north-be,
    glesys_loadbalancer_frontend.eu-north-fe
  ]
  loadbalancerid = glesys_loadbalancer.eu-north.id
  backend        = glesys_loadbalancer_backend.eu-north-be.id
  name     = "netherlands1"
  port     = 8080
  weight   = 15
  targetip = element(glesys_server.netherlands1.*.ipv4_address, 0)
  enabled = true
}


