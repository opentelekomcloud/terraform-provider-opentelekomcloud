
resource "opentelekomcloud_lb_loadbalancer_v2" "lb" {
  vip_subnet_id = var.subnetid
}

#binding eip for lb
resource "opentelekomcloud_networking_floatingip_associate_v2" "fip_1" {
  floating_ip = var.lb_eip
  port_id     = opentelekomcloud_lb_loadbalancer_v2.lb.vip_port_id
}

resource "opentelekomcloud_lb_loadbalancer_v2" "lb_1" {
  vip_subnet_id         = var.subnetid
  name                  = var.elb_name
  admin_state_up        = true
  description           = var.elb_desc
  loadbalancer_provider = "vlb"
}

resource "opentelekomcloud_lb_listener_v2" "listener_1" {
  protocol        = "HTTP"
  protocol_port   = 80
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v2.lb_1.id
}

resource "opentelekomcloud_lb_listener_v2" "listener_2" {
  protocol        = "HTTP"
  protocol_port   = 80
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v2.lb.id
  name            = "${var.elb_name}_listener"
  #default_pool_id  = opentelekomcloud_lb_pool_v2.pool_2.id
  description               = var.elb_desc
  default_tls_container_ref = "035bf725552a4836a1e809e69ac6b243" #serve
  sni_container_refs        = []
  admin_state_up            = true
}

resource "opentelekomcloud_lb_pool_v2" "pool_2" {
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v2.lb.id
}

resource "opentelekomcloud_lb_pool_v2" "pool_1" {
  name        = "${var.elb_name}_pool"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = opentelekomcloud_lb_listener_v2.listener_1.id
  description = var.elb_desc
  persistence {
    type        = "APP_COOKIE"
    cookie_name = "testCookie"
  }
}

resource "opentelekomcloud_lb_member_v2" "member_1" {
  name           = "${var.elb_name}_member"
  weight         = "2"
  admin_state_up = true
  address        = var.ecs_ip
  protocol_port  = 80
  pool_id        = opentelekomcloud_lb_pool_v2.pool_1.id
  subnet_id      = var.subnetid
}

resource "opentelekomcloud_lb_monitor_v2" "monitor_1" {
  pool_id        = opentelekomcloud_lb_pool_v2.pool_1.id
  type           = "HTTP"
  delay          = 5
  timeout        = 10
  max_retries    = 3
  name           = "${var.elb_name}_monitor"
  http_method    = "GET"
  expected_codes = "200"
  url_path       = "/"
}
