package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v1/extensions/vpcs"
)

// PASS
func TestAccVpcV1_basic(t *testing.T) {
	var vpc vpcs.Vpc

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpcV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1Exists("opentelekomcloud_vpc_v1.vpc_1", &vpc),
				),
			},
			resource.TestStep{
				Config: testAccVpcV1_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_v1.vpc_1", "name", "terraform_provider_test1"),
				),
			},
		},
	})
}
// PASS
func TestAccVpcV1_timeout(t *testing.T) {
	var vpc vpcs.Vpc

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpcV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1Exists("opentelekomcloud_vpc_v1.vpc_1", &vpc),
				),
			},
		},
	})
}

func testAccCheckVpcV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	vpcClient, err := config.vpcV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud vpc client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpc_v1" {
			continue
		}

		_, err := vpcs.Get(vpcClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Vpc still exists")
		}
	}

	return nil
}

func testAccCheckVpcV1Exists(n string, router *vpcs.Vpc) resource.TestCheckFunc {
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
			return fmt.Errorf("Error creating OpenTelekomCloud vpc client: %s", err)
		}

		found, err := vpcs.Get(vpcClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("vpc not found")
		}

		*router = *found

		return nil
	}
}

const testAccVpcV1_basic = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
	name = "terraform_provider_test"
	cidr="192.168.0.0/16"
}
`

const testAccVpcV1_update = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
    name = "terraform_provider_test1"
	cidr="192.168.0.0/16"
}
`
const testAccVpcV1_timeout = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
	name = "terraform_provider_test"
	cidr="192.168.0.0/16"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
