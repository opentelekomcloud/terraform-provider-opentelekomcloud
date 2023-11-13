#Test required params

resource "opentelekomcloud_vpc_eip_v1" "eip_2" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name       = "${var.bw_name}_2"
    size       = var.bw_size
    share_type = "PER"
  }
}

#create EIP,attach EIP
resource "opentelekomcloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name       = var.bw_name
    size       = var.bw_size
    share_type = "PER"
    #charge_mode = "traffic"
  }
}

resource "opentelekomcloud_vpc_eip_v1" "eip_lb" {
  publicip {
    type = "5_gray"
  }
  bandwidth {
    name        = "${var.bw_name}_lb"
    size        = var.bw_size
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "opentelekomcloud_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = opentelekomcloud_vpc_eip_v1.eip_1.publicip.0.ip_address
  instance_id = var.ecs_id
}
