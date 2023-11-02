module "network" {
  source = "../modules/network"
}

resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer" {
  router_id   = module.network.shared_subnet.vpc_id
  network_ids = [module.network.shared_subnet.network_id]

  availability_zones = ["eu-de-01"]
}

resource "opentelekomcloud_lb_ipgroup_v3" "group_1" {
  name        = "group_1"
  description = "some interesting description 1"

  ip_list {
    ip          = "192.168.10.10"
    description = "first"
  }
  depends_on = [opentelekomcloud_lb_loadbalancer_v3.loadbalancer]
}

resource "opentelekomcloud_lb_listener_v3" "listener" {
  loadbalancer_id     = opentelekomcloud_lb_loadbalancer_v3.loadbalancer.id
  protocol            = "HTTP"
  protocol_port       = 8080
  advanced_forwarding = true
  sni_match_algo      = "wildcard"

  insert_headers {
    forwarded_host = true
  }

  ip_group {
    id     = opentelekomcloud_lb_ipgroup_v3.group_1.id
    enable = true
  }
  depends_on = [
    opentelekomcloud_lb_loadbalancer_v3.loadbalancer,
    opentelekomcloud_lb_ipgroup_v3.group_1
  ]
}

resource "opentelekomcloud_lb_pool_v3" "pool" {
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.loadbalancer.id
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "HTTP"

  depends_on = [opentelekomcloud_lb_listener_v3.listener]
}

resource "opentelekomcloud_lb_policy_v3" "policy" {
  name             = "policy_updated"
  description      = "some interesting description"
  action           = "REDIRECT_TO_POOL"
  listener_id      = opentelekomcloud_lb_listener_v3.listener.id
  redirect_pool_id = opentelekomcloud_lb_pool_v3.pool.id
  position         = 37

  depends_on = [opentelekomcloud_lb_listener_v3.listener]
}

resource "opentelekomcloud_lb_rule_v3" "rule" {
  type         = "PATH"
  compare_type = "EQUAL_TO"
  value        = "/test"
  policy_id    = opentelekomcloud_lb_policy_v3.policy.id

  conditions {
    value = "/home"
  }
  depends_on = [opentelekomcloud_lb_listener_v3.listener]
}
