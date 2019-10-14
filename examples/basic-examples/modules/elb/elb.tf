resource "opentelekomcloud_elb_loadbalancer" "elb" {
  name = "${var.elb_name}_external"
  type = "External"
  description = "${var.elb_desc}"
  vpc_id = "${var.vpc_id}"
  admin_state_up = "true"
  bandwidth = "${var.bw_size}"  
  }

 #resource "opentelekomcloud_elb_loadbalancer" "elb2" {
 # name = "${var.elb_name}_internal"
 # type = "Internal"
 # vpc_id = "${var.vpc_id}"
 # admin_state_up = "false"
 # vip_subnet_id  = "${var.vip_subnet_id}"
 # security_group_id = "${var.security_group_id}"
 # }

resource "opentelekomcloud_elb_listener" "listener" {
  name = "${var.elb_name}_test-elb-listener"
  description = "great listener"
  protocol = "TCP"
  backend_protocol = "TCP"
  protocol_port = 12345
  backend_port = 8080
  lb_algorithm = "roundrobin"
  loadbalancer_id = "${opentelekomcloud_elb_loadbalancer.elb.id}"
  tcp_timeout = 4 
  tcp_draining  = "true"
  tcp_draining_timeout = 30
  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
  }
resource "opentelekomcloud_elb_listener" "listener2" {
  name = "${var.elb_name}_test-elb-listener2"
  description = "great listener"
  protocol = "HTTPS"
  backend_protocol = "TCP"
  protocol_port = 4444
  backend_port = 8080
  lb_algorithm = "roundrobin"
  loadbalancer_id = "${opentelekomcloud_elb_loadbalancer.elb.id}"
  tcp_timeout = 4 
  tcp_draining  = "true"
  tcp_draining_timeout = 30
  certificate_id = "7085e1a5a9aa41bfae51f655728479f0"
  ssl_protocols  = "TLSv1.2"	   
  ssl_ciphers  = "Default"
  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
  
}


resource "opentelekomcloud_elb_health" "healthcheck" {
  listener_id = "${opentelekomcloud_elb_listener.listener.id}"
  healthcheck_protocol = "TCP"
  healthcheck_connect_port = 22
  healthy_threshold = 5
  healthcheck_timeout = 25
  healthcheck_interval = 3
  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

#resource "opentelekomcloud_elb_health" "healthcheck2" {
#  listener_id = "${opentelekomcloud_elb_listener.listener.id}"
#}

resource "opentelekomcloud_elb_backend" "backend" {
  address = "${var.ecs_ip}"
  listener_id = "${opentelekomcloud_elb_listener.listener2.id}"
  server_id = "${var.ecs_id}"
}






