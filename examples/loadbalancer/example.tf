resource "glesys_loadbalancer" "mylb" {
  count      = 1
  datacenter = "Falkenberg"
  name       = "mylb-1"
}

resource "glesys_loadbalancer_backend" "mybackend" {
  loadbalancerid = glesys_loadbalancer.mylb[0].id
  name           = "my-web-backend"
  mode           = "tcp"
  connecttimeout = 20000
}

resource "glesys_loadbalancer_frontend" "myfrontend" {
  loadbalancerid = glesys_loadbalancer.mylb[0].id
  name           = "web-fe91"
  backend        = glesys_loadbalancer_backend.mybackend.id
  port           = 80
}

resource "glesys_loadbalancer_target" "web111" {
  loadbalancerid = glesys_loadbalancer.mylb[0].id
  backend        = glesys_loadbalancer_backend.mybackend.id

  name     = "web111"
  port     = 8898
  targetip = "172.16.0.10"
  weight   = 15

  enabled = false
}
