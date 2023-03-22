package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/ipgroups"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	elbv3 "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/elb/v3"
)

const resourceIpGroupName = "opentelekomcloud_lb_ipgroup_v3.group_1"

func TestAccLBV3IpGroup_basic(t *testing.T) {
	var ipgroup ipgroups.IpGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV3IpGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV3IpGroupConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV3IpGroupExists(resourceIpGroupName, &ipgroup),
					resource.TestCheckResourceAttr(resourceIpGroupName, "name", "group_1"),
					resource.TestCheckResourceAttr(resourceIpGroupName, "description", "some interesting description"),
					resource.TestCheckResourceAttr(resourceIpGroupName, "ip_list.#", "2"),
					resource.TestCheckResourceAttr(resourceIpGroupName, "ip_list.0.ip", "192.168.10.10"),
					resource.TestCheckResourceAttr(resourceIpGroupName, "ip_list.0.description", "first"),
					resource.TestCheckResourceAttr(resourceIpGroupName, "ip_list.1.ip", "192.168.10.11"),
					resource.TestCheckResourceAttr(resourceIpGroupName, "ip_list.1.description", "second"),
				),
			},
			{
				Config: testAccLBV3IpGroupConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV3IpGroupExists(resourceIpGroupName, &ipgroup),
					resource.TestCheckResourceAttr(resourceIpGroupName, "name", "group_1"),
					resource.TestCheckResourceAttr(resourceIpGroupName, "description", "description update"),
					resource.TestCheckResourceAttr(resourceIpGroupName, "ip_list.#", "3"),
				),
			},
		},
	})
}

func TestAccLBV3IpGroup_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV3IpGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV3IpGroupConfigBasic,
			},
			{
				ResourceName:      resourceIpGroupName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLBV3IpGroupDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ElbV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf(elbv3.ErrCreateClient, err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_lb_ipgroup_v3" {
			continue
		}

		_, err := ipgroups.Get(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("listener still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV3IpGroupExists(n string, listener *ipgroups.IpGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.ElbV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf(elbv3.ErrCreateClient, err)
		}

		found, err := ipgroups.Get(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("listener not found")
		}

		*listener = *found

		return nil
	}
}

const testAccLBV3IpGroupConfigBasic = `
resource "opentelekomcloud_lb_ipgroup_v3" "group_1" {
  name        = "group_1"
  description = "some interesting description"

  ip_list {
    ip          = "192.168.10.10"
	description = "first"
  }
  ip_list {
    ip          = "192.168.10.11"
	description = "second"
  }
}
`

const testAccLBV3IpGroupConfigUpdate = `
resource "opentelekomcloud_lb_ipgroup_v3" "group_1" {
  name        = "group_1"
  description = "description update"

  ip_list {
    ip          = "192.168.50.10"
	description = "one"
  }
  ip_list {
    ip          = "192.168.100.10"
	description = "two"
  }
  ip_list {
    ip          = "192.168.150.10"
	description = "three"
  }
}
`
