package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/security/groups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceNwSecGroupName = "opentelekomcloud_networking_secgroup_v2.secgroup_1"

func TestAccNetworkingV2SecGroup_basic(t *testing.T) {
	var securityGroup groups.SecGroup
	t.Parallel()
	th.AssertNoErr(t, quotas.SecurityGroup.Acquire())
	defer quotas.SecurityGroup.Release()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupBasic,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SecGroupExists(resourceNwSecGroupName, &securityGroup),
					testAccCheckNetworkingV2SecGroupRuleCount(&securityGroup, 2),
				),
			},
			{
				Config: testAccNetworkingV2SecGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceNwSecGroupName, "id", &securityGroup.ID),
					resource.TestCheckResourceAttr(resourceNwSecGroupName, "name", "security_group_2"),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroup_importBasic(t *testing.T) {
	t.Parallel()
	th.AssertNoErr(t, quotas.SecurityGroup.Acquire())
	defer quotas.SecurityGroup.Release()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupBasic,
			},
			{
				ResourceName:      resourceNwSecGroupName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkingV2SecGroup_noDefaultRules(t *testing.T) {
	var securityGroup groups.SecGroup
	t.Parallel()
	th.AssertNoErr(t, quotas.SecurityGroup.Acquire())
	defer quotas.SecurityGroup.Release()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupNoDefaultRules,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SecGroupExists(resourceNwSecGroupName, &securityGroup),
					testAccCheckNetworkingV2SecGroupRuleCount(&securityGroup, 0),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroup_timeout(t *testing.T) {
	var securityGroup groups.SecGroup
	t.Parallel()
	th.AssertNoErr(t, quotas.SecurityGroup.Acquire())
	defer quotas.SecurityGroup.Release()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupTimeout,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2SecGroupExists(resourceNwSecGroupName, &securityGroup),
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

const testAccNetworkingV2SecGroupBasic = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "security_group"
  description = "terraform security group acceptance test"
}
`

const testAccNetworkingV2SecGroupUpdate = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "security_group_2"
  description = "terraform security group acceptance test"
}
`

const testAccNetworkingV2SecGroupNoDefaultRules = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name                 = "security_group_1"
  description          = "terraform security group acceptance test"
  delete_default_rules = true
}
`

const testAccNetworkingV2SecGroupTimeout = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "security_group"
  description = "terraform security group acceptance test"

  timeouts {
    delete = "5m"
  }
}
`
