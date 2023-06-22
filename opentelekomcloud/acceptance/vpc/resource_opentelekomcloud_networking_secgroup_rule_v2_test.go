package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/security/groups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/security/rules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceNwSGRuleName = "opentelekomcloud_networking_secgroup_rule_v2.secgroup_rule_1"

func TestAccNetworkingV2SecGroupRule_basic(t *testing.T) {
	var secgroup1 groups.SecGroup
	var secgroup2 groups.SecGroup
	var secgroupRule1 rules.SecGroupRule
	var secgroupRule2 rules.SecGroupRule
	t.Parallel()
	qts := quotas.MultipleQuotas{
		{Q: quotas.SecurityGroup, Count: 2},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRuleBasic,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SecGroupExists(resourceNwSecGroupName, &secgroup1),
					TestAccCheckNetworkingV2SecGroupExists(
						"opentelekomcloud_networking_secgroup_v2.secgroup_2", &secgroup2),
					testAccCheckNetworkingV2SecGroupRuleExists(resourceNwSGRuleName, &secgroupRule1),
					testAccCheckNetworkingV2SecGroupRuleExists(
						"opentelekomcloud_networking_secgroup_rule_v2.secgroup_rule_2", &secgroupRule2),
					resource.TestCheckResourceAttr(resourceNwSGRuleName, "description", "test secgroup rule"),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroupRule_importBasic(t *testing.T) {
	t.Parallel()
	quotas.BookOne(t, quotas.SecurityGroup)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRuleNumericProtocol,
			},
			{
				ResourceName:      resourceNwSGRuleName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkingV2SecGroupRule_timeout(t *testing.T) {
	var secgroup_1 groups.SecGroup
	var secgroup_2 groups.SecGroup
	t.Parallel()
	qts := quotas.MultipleQuotas{
		{Q: quotas.SecurityGroup, Count: 2},
	}
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRuleTimeout,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SecGroupExists(resourceNwSecGroupName, &secgroup_1),
					TestAccCheckNetworkingV2SecGroupExists(
						"opentelekomcloud_networking_secgroup_v2.secgroup_2", &secgroup_2),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroupRule_numericProtocol(t *testing.T) {
	var secgroup1 groups.SecGroup
	var secgroupRule1 rules.SecGroupRule
	t.Parallel()
	quotas.BookOne(t, quotas.SecurityGroup)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRuleNumericProtocol,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SecGroupExists(resourceNwSecGroupName, &secgroup1),
					testAccCheckNetworkingV2SecGroupRuleExists(resourceNwSGRuleName, &secgroupRule1),
					resource.TestCheckResourceAttr(resourceNwSGRuleName, "protocol", "115"),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroupRule_noPorts(t *testing.T) {
	var secgroup1 groups.SecGroup
	var secgroupRule1 rules.SecGroupRule
	t.Parallel()
	quotas.BookOne(t, quotas.SecurityGroup)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRuleNoPorts,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SecGroupExists(resourceNwSecGroupName, &secgroup1),
					testAccCheckNetworkingV2SecGroupRuleExists(resourceNwSGRuleName, &secgroupRule1),
					resource.TestCheckResourceAttr(resourceNwSGRuleName, "direction", "egress"),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroupRule_ICMP(t *testing.T) {
	var secgroup1 groups.SecGroup
	var secgroupRule1 rules.SecGroupRule
	t.Parallel()
	quotas.BookOne(t, quotas.SecurityGroup)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRuleICMPEchoReply,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SecGroupExists(resourceNwSecGroupName, &secgroup1),
					testAccCheckNetworkingV2SecGroupRuleExists(resourceNwSGRuleName, &secgroupRule1),
					resource.TestCheckResourceAttr(resourceNwSGRuleName, "direction", "ingress"),
					resource.TestCheckResourceAttr(resourceNwSGRuleName, "port_range_min", "0"),
					resource.TestCheckResourceAttr(resourceNwSGRuleName, "port_range_max", "0"),
				),
			},
			{
				Config: testAccNetworkingV2SecGroupRuleICMPAll,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SecGroupExists(resourceNwSecGroupName, &secgroup1),
					testAccCheckNetworkingV2SecGroupRuleExists(resourceNwSGRuleName, &secgroupRule1),
					resource.TestCheckResourceAttr(resourceNwSGRuleName, "direction", "ingress"),
					resource.TestCheckResourceAttr(resourceNwSGRuleName, "port_range_min", "0"),
					resource.TestCheckResourceAttr(resourceNwSGRuleName, "port_range_max", "255"),
				),
			},
			{
				Config: testAccNetworkingV2SecGroupRuleICMPEcho,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SecGroupExists(resourceNwSecGroupName, &secgroup1),
					testAccCheckNetworkingV2SecGroupRuleExists(resourceNwSGRuleName, &secgroupRule1),
					resource.TestCheckResourceAttr(resourceNwSGRuleName, "direction", "ingress"),
					resource.TestCheckResourceAttr(resourceNwSGRuleName, "port_range_min", "8"),
					resource.TestCheckResourceAttr(resourceNwSGRuleName, "port_range_max", "0"),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroupRule_noProtocolError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNetworkingV2SecGroupRuleError,
				ExpectError: regexp.MustCompile(`"port_range_min": all of .+`),
			},
		},
	})
}

func testAccCheckNetworkingV2SecGroupRuleDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_networking_secgroup_rule_v2" {
			continue
		}

		_, err := rules.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("security group rule still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2SecGroupRuleExists(n string, securityGroupRule *rules.SecGroupRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		networkingClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
		}

		found, err := rules.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("security group rule not found")
		}

		*securityGroupRule = *found

		return nil
	}
}

const testAccNetworkingV2SecGroupRuleBasic = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "opentelekomcloud_networking_secgroup_v2" "secgroup_2" {
  name        = "secgroup_2"
  description = "terraform security group rule acceptance test"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_1" {
  description       = "test secgroup rule"
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_max    = 22
  port_range_min    = 22
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_1.id
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_2" {
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_max    = 80
  port_range_min    = 80
  protocol          = "tcp"
  remote_group_id   = opentelekomcloud_networking_secgroup_v2.secgroup_1.id
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_2.id
}
`

const testAccNetworkingV2SecGroupRuleTimeout = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "opentelekomcloud_networking_secgroup_v2" "secgroup_2" {
  name        = "secgroup_2"
  description = "terraform security group rule acceptance test"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_1" {
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_max    = 22
  port_range_min    = 22
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_1.id

  timeouts {
    delete = "5m"
  }
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_2" {
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_max    = 80
  port_range_min    = 80
  protocol          = "tcp"
  remote_group_id   = opentelekomcloud_networking_secgroup_v2.secgroup_1.id
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_2.id

  timeouts {
    delete = "5m"
  }
}
`

const testAccNetworkingV2SecGroupRuleNumericProtocol = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_1" {
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_max    = 22
  port_range_min    = 22
  protocol          = "115"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_1.id
}
`

const testAccNetworkingV2SecGroupRuleNoPorts = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_1" {
  direction         = "egress"
  ethertype         = "IPv4"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_1.id
}
`

const testAccNetworkingV2SecGroupRuleError = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_1" {
  direction         = "egress"
  ethertype         = "IPv4"
  remote_ip_prefix  = "0.0.0.0/0"
  port_range_min    = 0
  port_range_max    = 22
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_1.id
}
`

const testAccNetworkingV2SecGroupRuleICMPEcho = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_1" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  port_range_min    = 8
  port_range_max    = 0
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_1.id
}
`

const testAccNetworkingV2SecGroupRuleICMPAll = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_1" {
  description       = "all ICMP ports"
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  port_range_min    = 0
  port_range_max    = 255
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_1.id
}
`

const testAccNetworkingV2SecGroupRuleICMPEchoReply = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "terraform security group rule acceptance test"
}

resource "opentelekomcloud_networking_secgroup_rule_v2" "secgroup_rule_1" {
  description       = "echo ICMP"
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  port_range_min    = 0
  port_range_max    = 0
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_1.id
}
`
