package src

const WafdMain = `
##############
# NETWORK part
##############

resource "opentelekomcloud_vpc_v1" "vpc" {
  name   = var.environment
  cidr   = var.vpc_cidr
  shared = true
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet" {
  name          = var.environment
  vpc_id        = opentelekomcloud_vpc_v1.vpc.id
  cidr          = var.subnet_cidr
  gateway_ip    = var.subnet_gateway_ip
  primary_dns   = var.subnet_primary_dns
  secondary_dns = var.subnet_secondary_dns
}

data "opentelekomcloud_networking_secgroup_v2" "default_secgroup" {
  name = "default"
}

##################
# WAFD DOMAIN part
##################
resource "opentelekomcloud_waf_dedicated_policy_v1" "policy_1" {
  name = "policy_throttling"
}

resource "opentelekomcloud_waf_dedicated_domain_v1" "domain_1" {
  domain      = "www.wafd.throttling-test.com"
  keep_policy = true
  proxy       = true

  policy_id = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id

  server {
    client_protocol = "HTTP"
    server_protocol = "HTTP"
    address         = "10.1.0.10"
    port            = 8080
    type            = "ipv4"
    vpc_id          = opentelekomcloud_vpc_subnet_v1.subnet.vpc_id
  }
}

######################
# WAFD RULES part / 30
######################

resource "opentelekomcloud_waf_dedicated_cc_rule_v1" "rule_cc" {
  count        = 30
  policy_id    = opentelekomcloud_waf_dedicated_policy_v1.policy_1.id
  mode         = 0
  url          = "/abc_${count.index}"
  limit_num    = 10
  limit_period = 60
  lock_time    = 10
  tag_type     = "cookie"
  tag_index    = "sessionid"

  action {
    category     = "block"
    content_type = "application/json"
    content      = "{\"error\":\"forbidden\"}"
  }
}
`
