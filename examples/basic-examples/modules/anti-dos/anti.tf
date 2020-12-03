
resource "opentelekomcloud_vpc_eip_v1" "eip_anti" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name       = "test_anti"
    size       = 5
    share_type = "PER"
  }
}
resource "opentelekomcloud_antiddos_v1" "myantiddos" {
  floating_ip_id         = opentelekomcloud_vpc_eip_v1.eip_anti.id
  enable_l7              = true
  traffic_pos_id         = 1
  http_request_pos_id    = 3
  cleaning_access_pos_id = 2
  app_type_id            = 0
}
