resource "opentelekomcloud_dns_zone_v2" "zone" {
  name = "test.hanmeina."
}

resource "opentelekomcloud_dns_zone_v2" "private_zone" {
  name        = var.zone_name
  description = var.zone_desc
  ttl         = 3000
  #type = "private"
  type = var.zone_type
  router {
    router_id     = var.vpc_id
    router_region = var.region
  }
}

resource "opentelekomcloud_dns_recordset_v2" "example_com" {
  zone_id     = opentelekomcloud_dns_zone_v2.private_zone.id
  name        = var.recordset_name
  description = "${var.recordset_desc}_test"
  ttl         = 3000
  type        = var.recordset_type
  #records = [var.records]
  records = [
    "10.1.0.9",
    "10.1.0.11"
  ]
}
