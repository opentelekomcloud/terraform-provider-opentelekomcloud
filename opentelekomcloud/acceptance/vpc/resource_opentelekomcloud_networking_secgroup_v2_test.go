package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/security/groups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccNetworkingV2SecGroup_basic(t *testing.T) {
	var securityGroup groups.SecGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroup_basic,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SecGroupExists("opentelekomcloud_networking_secgroup_v2.secgroup_1", &securityGroup),
					testAccCheckNetworkingV2SecGroupRuleCount(&securityGroup, 2),
				),
			},
			{
				Config: testAccNetworkingV2SecGroup_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr("opentelekomcloud_networking_secgroup_v2.secgroup_1", "id", &securityGroup.ID),
					resource.TestCheckResourceAttr("opentelekomcloud_networking_secgroup_v2.secgroup_1", "name", "security_group_2"),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroup_noDefaultRules(t *testing.T) {
	var securityGroup groups.SecGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroup_noDefaultRules,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SecGroupExists("opentelekomcloud_networking_secgroup_v2.secgroup_1", &securityGroup),
					testAccCheckNetworkingV2SecGroupRuleCount(&securityGroup, 0),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroup_timeout(t *testing.T) {
	var securityGroup groups.SecGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroup_timeout,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SecGroupExists("opentelekomcloud_networking_secgroup_v2.secgroup_1", &securityGroup),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2SecGroupDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud Networkingv2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_networking_secgroup_v2" {
			continue
		}

		_, err := groups.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("security group still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2SecGroupRuleCount(
	sg *groups.SecGroup, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(sg.Rules) == count {
			return nil
		}

		return fmt.Errorf("unexpected number of rules in group %s. Expected %d, got %d",
			sg.ID, count, len(sg.Rules))
	}
}

const testAccNetworkingV2SecGroup_basic = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "security_group"
  description = "terraform security group acceptance test"
}
`

const testAccNetworkingV2SecGroup_update = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "security_group_2"
  description = "terraform security group acceptance test"
}
`

const testAccNetworkingV2SecGroup_noDefaultRules = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name                 = "security_group_1"
  description          = "terraform security group acceptance test"
  delete_default_rules = true
}
`

const testAccNetworkingV2SecGroup_timeout = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "security_group"
  description = "terraform security group acceptance test"

  timeouts {
    delete = "5m"
  }
}
`
