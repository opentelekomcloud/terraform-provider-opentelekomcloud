resource "opentelekomcloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name = "test"
    size = 8
    share_type = "PER"
    charge_mode = "traffic"
  }
}

resource "opentelekomcloud_antiddos_v1" "antiddos_1" {
  floating_ip_id = "${opentelekomcloud_vpc_eip_v1.eip_1.id}"
  enable_l7 = true
  traffic_pos_id = 1
  http_request_pos_id = 2
  cleaning_access_pos_id = 1
  app_type_id = 0
}