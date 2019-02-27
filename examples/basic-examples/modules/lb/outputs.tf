

output "loadbalancerid"{
  value = "${opentelekomcloud_lb_loadbalancer_v2.lb_1.id}"
}

output "listenerid"{
  value = "${opentelekomcloud_lb_listener_v2.listener_1.id}"
}

