package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/subnets"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccOTCVpcSubnetV1Basic(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckOTCVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOTCVpcSubnetV1Basic,
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
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_subnet_v1.subnet_1", "ntp_addresses", "10.100.0.33,10.100.0.34"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_subnet_v1.subnet_1", "tags.foo", "bar"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_subnet_v1.subnet_1", "tags.key", "value"),
				),
			},
			{
				Config: testAccOTCVpcSubnetV1Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_subnet_v1.subnet_1", "name", "opentelekomcloud_subnet_1"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_subnet_v1.subnet_1", "ntp_addresses", "10.100.0.35,10.100.0.36"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_subnet_v1.subnet_1", "tags.key", "value_update"),
				),
			},
		},
	})
}

func TestAccOTCVpcSubnetV1Timeout(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckOTCVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOTCVpcSubnetV1Timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCVpcSubnetV1Exists("opentelekomcloud_vpc_subnet_v1.subnet_1", &subnet),
				),
			},
		},
	})
}

func TestAccOTCVpcSubnetV1DnsList(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckOTCVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOTCVpcSubnetV1DnsList,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCVpcSubnetV1Exists("opentelekomcloud_vpc_subnet_v1.subnet_1", &subnet),
				),
			},
		},
	})
}

func testAccCheckOTCVpcSubnetV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	subnetClient, err := config.NetworkingV1Client(env.OS_REGION_NAME)
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

		config := common.TestAccProvider.Meta().(*cfg.Config)
		subnetClient, err := config.NetworkingV1Client(env.OS_REGION_NAME)
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

const (
	testAccOTCVpcSubnetV1Basic = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name = "opentelekomcloud_subnet"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = opentelekomcloud_vpc_v1.vpc_1.id
  availability_zone = "eu-de-02"
  ntp_addresses = "10.100.0.33,10.100.0.34"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`
	testAccOTCVpcSubnetV1Update = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name = "opentelekomcloud_subnet_1"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = opentelekomcloud_vpc_v1.vpc_1.id
  availability_zone = "eu-de-02"
  ntp_addresses = "10.100.0.35,10.100.0.36"

  tags = {
    foo = "bar"
    key = "value_update"
  }
}
`

	testAccOTCVpcSubnetV1Timeout = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name = "opentelekomcloud_subnet"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = opentelekomcloud_vpc_v1.vpc_1.id
  availability_zone = "eu-de-02"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

	testAccOTCVpcSubnetV1DnsList = `
resource "opentelekomcloud_vpc_v1" "vpc" {
  name       = "vpc_name"
  cidr       = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name          = "subnet_name"
  vpc_id        = opentelekomcloud_vpc_v1.vpc.id
  cidr          = cidrsubnet(opentelekomcloud_vpc_v1.vpc.cidr, 8, 0)
  gateway_ip    = cidrhost(cidrsubnet(opentelekomcloud_vpc_v1.vpc.cidr, 8, 0), 1)
  dns_list = ["100.125.4.25", "8.8.8.8"]
}
`
)
