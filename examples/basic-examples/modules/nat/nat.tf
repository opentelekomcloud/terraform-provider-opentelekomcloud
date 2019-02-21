resource "opentelekomcloud_nat_gateway_v2" "nat_1" {
  name = "${var.nat_name}"
  description = "${var.nat_desc}"
  spec = "3"
  router_id = "${var.vpc_id}"
  internal_network_id = "${var.subnet1_id}"
}

resource "opentelekomcloud_nat_snat_rule_v2" "snat_1" {
  nat_gateway_id = "${opentelekomcloud_nat_gateway_v2.nat_1.id}"
  network_id = "${var.subnet1_id}"
  floating_ip_id = "${var.eip_id}"
}