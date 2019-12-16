package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/security/groups"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/security/rules"
)

func TestAccNetworkingV2SecGroupRule_basic(t *testing.T) {
	var secgroup_1 groups.SecGroup
	var secgroup_2 groups.SecGroup
	var secgroup_rule_1 rules.SecGroupRule
	var secgroup_rule_2 rules.SecGroupRule

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(
						"opentelekomcloud_networking_secgroup_v2.secgroup_1", &secgroup_1),
					testAccCheckNetworkingV2SecGroupExists(
						"opentelekomcloud_networking_secgroup_v2.secgroup_2", &secgroup_2),
					testAccCheckNetworkingV2SecGroupRuleExists(
						"opentelekomcloud_networking_secgroup_rule_v2.secgroup_rule_1", &secgroup_rule_1),
					testAccCheckNetworkingV2SecGroupRuleExists(
						"opentelekomcloud_networking_secgroup_rule_v2.secgroup_rule_2", &secgroup_rule_2),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_networking_secgroup_rule_v2.secgroup_rule_1", "description", "test secgroup rule"),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroupRule_timeout(t *testing.T) {
	var secgroup_1 groups.SecGroup
	var secgroup_2 groups.SecGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRule_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(
						"opentelekomcloud_networking_secgroup_v2.secgroup_1", &secgroup_1),
					testAccCheckNetworkingV2SecGroupExists(
						"opentelekomcloud_networking_secgroup_v2.secgroup_2", &secgroup_2),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroupRule_numericProtocol(t *testing.T) {
	var secgroup_1 groups.SecGroup
	var secgroup_rule_1 rules.SecGroupRule

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRule_numericProtocol,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(
						"opentelekomcloud_networking_secgroup_v2.secgroup_1", &secgroup_1),
					testAccCheckNetworkingV2SecGroupRuleExists(
						"opentelekomcloud_networking_secgroup_rule_v2.secgroup_rule_1", &secgroup_rule_1),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_networking_secgroup_rule_v2.secgroup_rule_1", "protocol", "115"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2SecGroupRuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_networking_secgroup_rule_v2" {
			continue
		}

		_, err := rules.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Security group rule still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2SecGroupRuleExists(n string, security_group_rule *rules.SecGroupRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
		}

		found, err := rules.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Security group rule not found")
		}

		*security_group_rule = *found

		return nil
	}
}

const testAccNetworkingV2SecGroupRule_basic = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "opentelekomcloud_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group rule acceptance test"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_1" {
  description = "test secgroup rule"
  direction = "ingress"
  ethertype = "IPv4"
  port_range_max = 22
  port_range_min = 22
  protocol = "tcp"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_2" {
  direction = "ingress"
  ethertype = "IPv4"
  port_range_max = 80
  port_range_min = 80
  protocol = "tcp"
  remote_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_2.id}"
}
`

const testAccNetworkingV2SecGroupRule_lowerCaseCIDR = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_1" {
  direction = "ingress"
  ethertype = "IPv6"
  port_range_max = 22
  port_range_min = 22
  protocol = "tcp"
  remote_ip_prefix = "2001:558:FC00::/39"
  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
}
`

const testAccNetworkingV2SecGroupRule_timeout = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "opentelekomcloud_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group rule acceptance test"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_1" {
  direction = "ingress"
  ethertype = "IPv4"
  port_range_max = 22
  port_range_min = 22
  protocol = "tcp"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"

  timeouts {
    delete = "5m"
  }
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_2" {
  direction = "ingress"
  ethertype = "IPv4"
  port_range_max = 80
  port_range_min = 80
  protocol = "tcp"
  remote_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_2.id}"

  timeouts {
    delete = "5m"
  }
}
`

const testAccNetworkingV2SecGroupRule_protocols = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

# NOT SUPPORTED
#resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_ah" {
#  direction = "ingress"
#  ethertype = "IPv4"
#  protocol = "ah"
#  remote_ip_prefix = "0.0.0.0/0"
#  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
#}

# NOT SUPPORTED
#resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_dccp" {
#  direction = "ingress"
#  ethertype = "IPv4"
#  protocol = "dccp"
#  remote_ip_prefix = "0.0.0.0/0"
#  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
#}

# NOT SUPPORTED
#resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_egp" {
#  direction = "ingress"
#  ethertype = "IPv4"
#  protocol = "egp"
#  remote_ip_prefix = "0.0.0.0/0"
#  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
#}

# NOT SUPPORTED
#resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_esp" {
#  direction = "ingress"
#  ethertype = "IPv4"
#  protocol = "esp"
#  remote_ip_prefix = "0.0.0.0/0"
#  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
#}

# NOT SUPPORTED
#resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_gre" {
#  direction = "ingress"
#  ethertype = "IPv4"
#  protocol = "gre"
#  remote_ip_prefix = "0.0.0.0/0"
#  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
#}

# NOT SUPPORTED
#resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_igmp" {
#  direction = "ingress"
#  ethertype = "IPv4"
#  protocol = "igmp"
#  remote_ip_prefix = "0.0.0.0/0"
#  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
#}

# NOT SUPPORTED
#resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_ipv6_encap" {
#  direction = "ingress"
#  ethertype = "IPv6"
#  protocol = "ipv6-encap"
#  remote_ip_prefix = "::/0"
#  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
#}

# NOT SUPPORTED
#resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_ipv6_frag" {
#  direction = "ingress"
#  ethertype = "IPv6"
#  protocol = "ipv6-frag"
#  remote_ip_prefix = "::/0"
#  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
#}

# NOT SUPPORTED
#resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_ipv6_icmp" {
#  direction = "ingress"
#  ethertype = "IPv6"
#  protocol = "ipv6-icmp"
#  remote_ip_prefix = "::/0"
#  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
#}

# NOT SUPPORTED
#resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_ipv6_nonxt" {
#  direction = "ingress"
#  ethertype = "IPv6"
#  protocol = "ipv6-nonxt"
#  remote_ip_prefix = "::/0"
#  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
#}

# NOT SUPPORTED
#resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_ipv6_opts" {
#  direction = "ingress"
#  ethertype = "IPv6"
#  protocol = "ipv6-opts"
#  remote_ip_prefix = "::/0"
#  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
#}

# NOT SUPPORTED
#resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_ipv6_route" {
#  direction = "ingress"
#  ethertype = "IPv6"
#  protocol = "ipv6-route"
#  remote_ip_prefix = "::/0"
#  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
#}

# NOT SUPPORTED
#resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_ospf" {
#  direction = "ingress"
#  ethertype = "IPv4"
#  protocol = "ospf"
#  remote_ip_prefix = "0.0.0.0/0"
#  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
#}

# NOT SUPPORTED
#resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_pgm" {
#  direction = "ingress"
#  ethertype = "IPv4"
#  protocol = "pgm"
#  remote_ip_prefix = "0.0.0.0/0"
#  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
#}

# NOT SUPPORTED
#resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_rsvp" {
#  direction = "ingress"
#  ethertype = "IPv4"
#  protocol = "rsvp"
#  remote_ip_prefix = "0.0.0.0/0"
#  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
#}

# NOT SUPPORTED
#resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_sctp" {
#  direction = "ingress"
#  ethertype = "IPv4"
#  protocol = "sctp"
#  remote_ip_prefix = "0.0.0.0/0"
#  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
#}

# NOT SUPPORTED
#resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_udplite" {
#  direction = "ingress"
#  ethertype = "IPv4"
#  protocol = "udplite"
#  remote_ip_prefix = "0.0.0.0/0"
#  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
#}

# NOT SUPPORTED
#resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_vrrp" {
#  direction = "ingress"
#  ethertype = "IPv4"
#  protocol = "vrrp"
#  remote_ip_prefix = "0.0.0.0/0"
#  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
#}
`

const testAccNetworkingV2SecGroupRule_numericProtocol = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_1" {
  direction = "ingress"
  ethertype = "IPv4"
  port_range_max = 22
  port_range_min = 22
  protocol = "115"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
}
`
