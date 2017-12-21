resource "opentelekomcloud_vpc_v1" "vpc_v1" {
  count = "${var.instance_count}"
  name = "${var.project}-vpc${format("%02d", count.index+1)}"
  cidr = "${var.vpc_cidr}"
}

