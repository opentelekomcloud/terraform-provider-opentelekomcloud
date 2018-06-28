package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v1/subnets"
)

func TestAccOTCVpcSubnetV1_basic(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOTCVpcSubnetV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCVpcSubnetV1Exists("opentelekomcloud_vpc_subnet_v1.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_subnet_v1.subnet_1", "name", "opentelekomcloud_subnet"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_subnet_v1.subnet_1", "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_subnet_v1.subnet_1", "gateway_ip", "192.168.0.1"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_subnet_v1.subnet_1", "availability_zone", "eu-de-02"),
				),
			},
			resource.TestStep{
				Config: testAccOTCVpcSubnetV1_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_subnet_v1.subnet_1", "name", "opentelekomcloud_subnet_1"),
				),
			},
		},
	})
}

// PASS
func TestAccOTCVpcSubnetV1_timeout(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOTCVpcSubnetV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCVpcSubnetV1Exists("opentelekomcloud_vpc_subnet_v1.subnet_1", &subnet),
				),
			},
		},
	})
}

func testAccCheckOTCVpcSubnetV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	subnetClient, err := config.networkingV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud vpc client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpc_subnet_v1" {
			continue
		}

		_, err := subnets.Get(subnetClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Subnet still exists")
		}
	}

	return nil
}
func testAccCheckOTCVpcSubnetV1Exists(n string, subnet *subnets.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		subnetClient, err := config.networkingV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud Vpc client: %s", err)
		}

		found, err := subnets.Get(subnetClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Subnet not found")
		}

		*subnet = *found

		return nil
	}
}

const testAccOTCVpcSubnetV1_basic = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name = "opentelekomcloud_subnet"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
  availability_zone = "eu-de-02"

}
`
const testAccOTCVpcSubnetV1_update = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name = "opentelekomcloud_subnet_1"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
  availability_zone = "eu-de-02"
 }
`

const testAccOTCVpcSubnetV1_timeout = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name = "opentelekomcloud_subnet"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
  availability_zone = "eu-de-02"

 timeouts {
    create = "5m"
    delete = "5m"
  }

}
`
