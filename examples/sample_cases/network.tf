resource "opentelekomcloud_networking_network_v2" "network" {
  count          = var.instance_count
  name           = "${var.project}-network"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet" {
  name            = "${var.project}-subnet"
  count           = var.instance_count
  network_id      = opentelekomcloud_networking_network_v2.network.id
  cidr            = var.subnet_cidr
  ip_version      = 4
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}

resource "opentelekomcloud_networking_router_v2" "router" {
  count            = var.instance_count
  name             = "${var.project}-router"
  admin_state_up   = "true"
  external_gateway = "0a2228f2-7f8a-45f1-8e09-9039e1d09975"
}

resource "opentelekomcloud_networking_router_interface_v2" "interface" {
  count     = var.instance_count
  router_id = opentelekomcloud_networking_router_v2.router.id
  subnet_id = opentelekomcloud_networking_subnet_v2.subnet.id
}

