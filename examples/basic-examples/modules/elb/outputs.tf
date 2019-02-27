

output "loadbalancerid"{
  value = "${opentelekomcloud_elb_loadbalancer.elb.id}"
}

output "listenerid"{
  value = "${opentelekomcloud_elb_listener.listener.id}"
}

