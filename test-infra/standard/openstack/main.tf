variable "key_name" {}
variable "public_key" {}
variable "network_id" {}

variable "pool" {
  default = "public"
}

variable "flavor" {
  default = "m1.xlarge"
}

data "huaweicloud_images_image_v2" "packstack_standard" {
  name = "packstack-standard-ocata"
  most_recent = true
}

resource "huaweicloud_networking_floatingip_v2" "huaweicloud_acc_tests" {
  pool = "${var.pool}"
}

resource "huaweicloud_networking_secgroup_v2" "huaweicloud_acc_tests" {
  name = "huaweicloud_acc_tests"
  description = "Rules for openstack acceptance tests"
}

resource "huaweicloud_networking_secgroup_rule_v2" "huaweicloud_acc_tests_rule_1" {
  security_group_id = "${huaweicloud_networking_secgroup_v2.huaweicloud_acc_tests.id}"
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "tcp"
  port_range_min = 1
  port_range_max = 65535
  remote_ip_prefix = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule_v2" "huaweicloud_acc_tests_rule_2" {
  security_group_id = "${huaweicloud_networking_secgroup_v2.huaweicloud_acc_tests.id}"
  direction = "ingress"
  ethertype = "IPv6"
  protocol = "tcp"
  port_range_min = 1
  port_range_max = 65535
  remote_ip_prefix = "::/0"
}

resource "huaweicloud_networking_secgroup_rule_v2" "huaweicloud_acc_tests_rule_3" {
  security_group_id = "${huaweicloud_networking_secgroup_v2.huaweicloud_acc_tests.id}"
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "udp"
  port_range_min = 1
  port_range_max = 65535
  remote_ip_prefix = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule_v2" "huaweicloud_acc_tests_rule_4" {
  security_group_id = "${huaweicloud_networking_secgroup_v2.huaweicloud_acc_tests.id}"
  direction = "ingress"
  ethertype = "IPv6"
  protocol = "udp"
  port_range_min = 1
  port_range_max = 65535
  remote_ip_prefix = "::/0"
}

resource "huaweicloud_networking_secgroup_rule_v2" "huaweicloud_acc_tests_rule_5" {
  security_group_id = "${huaweicloud_networking_secgroup_v2.huaweicloud_acc_tests.id}"
  direction = "ingress"
  ethertype = "IPv4"
  protocol = "icmp"
  remote_ip_prefix = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule_v2" "huaweicloud_acc_tests_rule_6" {
  security_group_id = "${huaweicloud_networking_secgroup_v2.huaweicloud_acc_tests.id}"
  direction = "ingress"
  ethertype = "IPv6"
  protocol = "icmp"
  remote_ip_prefix = "::/0"
}

resource "huaweicloud_compute_instance_v2" "huaweicloud_acc_tests" {
  name = "huaweicloud_acc_tests"
  image_id = "${data.huaweicloud_images_image_v2.packstack_standard.id}"
  flavor_name = "${var.flavor}"
  key_pair = "${var.key_name}"

  security_groups = ["${huaweicloud_networking_secgroup_v2.huaweicloud_acc_tests.name}"]

  network {
    uuid = "${var.network_id}"
  }
}

resource "huaweicloud_compute_floatingip_associate_v2" "huaweicloud_acc_tests" {
  instance_id = "${huaweicloud_compute_instance_v2.huaweicloud_acc_tests.id}"
  floating_ip = "${huaweicloud_networking_floatingip_v2.huaweicloud_acc_tests.address}"
}

resource "null_resource" "rc_files" {
  provisioner "local-exec" {
    command = <<EOF
      while true ; do
        wget http://${huaweicloud_compute_floatingip_associate_v2.huaweicloud_acc_tests.floating_ip}/keystonerc_demo 2> /dev/null
        if [ $? = 0 ]; then
          break
        fi
        sleep 20
      done

      wget http://${huaweicloud_compute_floatingip_associate_v2.huaweicloud_acc_tests.floating_ip}/keystonerc_admin
    EOF
  }
}
