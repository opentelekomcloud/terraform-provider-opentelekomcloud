package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v1/extensions/vpcs"
)

// PASS
func TestAccNetworkingV1Vpc_basic(t *testing.T) {
	var vpc_group vpcs.Vpc

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV1VpcDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV1Vpc_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV1VpcExists(
						"opentelekomcloud_networking_vpc_v1.networking_vpc_v1", &vpc_group),
					testAccCheckNetworkingV1VpcCount(&vpc_group, 2),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV1Vpc_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(
						"opentelekomcloud_networking_vpc_v1.networking_vpc_v1", "id", &vpc_group.ID),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_networking_vpc_v1.networking_vpc_v1", "name", "modifyVpc"),
				),
			},
		},
	})
}

// PASS
func TestAccNetworkingV1Vpc_noDefaultRules(t *testing.T) {
	var vpc_group vpcs.Vpc

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV1VpcDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV1Vpc_noDefaultRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV1VpcExists(
						"opentelekomcloud_networking_vpc_v1.networking_vpc_v1", &vpc_group),
					testAccCheckNetworkingV1VpcCount(&vpc_group, 0),
				),
			},
		},
	})
}

// PASS
func TestAccNetworkingV1Vpc_timeout(t *testing.T) {
	var vpc_group vpcs.Vpc

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV1VpcDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV1Vpc_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV1VpcExists(
						"opentelekomcloud_networking_vpc_v1.networking_vpc_v1", &vpc_group),
				),
			},
		},
	})
}

func testAccCheckNetworkingV1VpcDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	vpcClient, err := config.vpcV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_networking_secgroup_v2" {
			continue
		}

		_, err := vpcs.Get(vpcClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Security group still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV1VpcExists(n string, vpc *vpcs.Vpc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		vpcClient, err := config.vpcV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
		}

		found, err := vpcs.Get(vpcClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Security group not found")
		}

		*vpc = *found

		return nil
	}
}

func testAccCheckNetworkingV1VpcCount(
	sg *vpcs.Vpc, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		/*if len(sg.Rules) == count {
			return nil
		}*/

		return fmt.Errorf("Unexpected number of rules in group %s. Expected %d, got %d",
			sg.ID, count/*, len(sg.Rules)*/)
	}
}

const testAccNetworkingV1Vpc_basic = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name = "testVpc"

}
`

const testAccNetworkingV1Vpc_update = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name = "modifyVpc"

}
`

const testAccNetworkingV1Vpc_noDefaultRules = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
	name = "security_group_1"

	delete_default_rules = true
}
`

const testAccNetworkingV1Vpc_timeout = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name = "testVpc"


  timeouts {
    delete = "5m"
  }
}
`
