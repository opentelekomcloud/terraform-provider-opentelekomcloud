package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/subnets"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccVpcSubnetV1Basic(t *testing.T) {
	var subnet subnets.Subnet
	resourceName := "opentelekomcloud_vpc_subnet_v1.subnet_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetV1Exists(resourceName, &subnet),
					resource.TestCheckResourceAttr(resourceName, "name", "opentelekomcloud_subnet"),
					resource.TestCheckResourceAttr(resourceName, "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "gateway_ip", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "eu-de-02"),
					resource.TestCheckResourceAttr(resourceName, "ntp_addresses", "10.100.0.33,10.100.0.34"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccVpcSubnetV1Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "opentelekomcloud_subnet_1"),
					resource.TestCheckResourceAttr(resourceName, "ntp_addresses", "10.100.0.35,10.100.0.36"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
				),
			},
		},
	})
}

func TestAccVpcSubnetV1Timeout(t *testing.T) {
	var subnet subnets.Subnet
	resourceName := "opentelekomcloud_vpc_subnet_v1.subnet_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetV1Timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetV1Exists(resourceName, &subnet),
				),
			},
		},
	})
}

func TestAccVpcSubnetV1DnsList(t *testing.T) {
	var subnet subnets.Subnet
	resourceName := "opentelekomcloud_vpc_subnet_v1.subnet_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetV1DnsList,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetV1Exists(resourceName, &subnet),
				),
			},
		},
	})
}

func testAccCheckVpcSubnetV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.NetworkingV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV1 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpc_subnet_v1" {
			continue
		}

		_, err := subnets.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("subnet still exists")
		}
	}

	return nil
}
func testAccCheckVpcSubnetV1Exists(n string, subnet *subnets.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.NetworkingV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud NetworkingV1 client: %w", err)
		}

		found, err := subnets.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("subnet not found")
		}

		*subnet = *found

		return nil
	}
}

const (
	testAccVpcSubnetV1Basic = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name              = "opentelekomcloud_subnet"
  cidr              = "192.168.0.0/16"
  gateway_ip        = "192.168.0.1"
  vpc_id            = opentelekomcloud_vpc_v1.vpc_1.id
  availability_zone = "eu-de-02"
  ntp_addresses     = "10.100.0.33,10.100.0.34"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`
	testAccVpcSubnetV1Update = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name              = "opentelekomcloud_subnet_1"
  cidr              = "192.168.0.0/16"
  gateway_ip        = "192.168.0.1"
  vpc_id            = opentelekomcloud_vpc_v1.vpc_1.id
  availability_zone = "eu-de-02"
  ntp_addresses     = "10.100.0.35,10.100.0.36"

  tags = {
    foo = "bar"
    key = "value_update"
  }
}
`

	testAccVpcSubnetV1Timeout = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name              = "opentelekomcloud_subnet"
  cidr              = "192.168.0.0/16"
  gateway_ip        = "192.168.0.1"
  vpc_id            = opentelekomcloud_vpc_v1.vpc_1.id
  availability_zone = "eu-de-02"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

	testAccVpcSubnetV1DnsList = `
resource "opentelekomcloud_vpc_v1" "vpc" {
  name = "vpc_name"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name       = "subnet_name"
  vpc_id     = opentelekomcloud_vpc_v1.vpc.id
  cidr       = cidrsubnet(opentelekomcloud_vpc_v1.vpc.cidr, 8, 0)
  gateway_ip = cidrhost(cidrsubnet(opentelekomcloud_vpc_v1.vpc.cidr, 8, 0), 1)
  dns_list   = ["100.125.4.25", "8.8.8.8"]
}
`
)
