package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/subnets"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceVPCSubnetName = "opentelekomcloud_vpc_subnet_v1.subnet_1"

func TestAccVpcSubnetV1Basic(t *testing.T) {
	var subnet subnets.Subnet
	t.Parallel()
	qts := vpcSubnetQuotas()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetV1Exists(resourceVPCSubnetName, &subnet),
					resource.TestCheckResourceAttr(resourceVPCSubnetName, "name", "subnet_name"),
					resource.TestCheckResourceAttr(resourceVPCSubnetName, "description", "some description"),
					resource.TestCheckResourceAttr(resourceVPCSubnetName, "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceVPCSubnetName, "gateway_ip", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceVPCSubnetName, "availability_zone", "eu-de-02"),
					resource.TestCheckResourceAttr(resourceVPCSubnetName, "ntp_addresses", "10.100.0.33,10.100.0.34"),
					resource.TestCheckResourceAttr(resourceVPCSubnetName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceVPCSubnetName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceVPCSubnetName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccVpcSubnetV1Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVPCSubnetName, "name", "subnet_name_update"),
					resource.TestCheckResourceAttr(resourceVPCSubnetName, "description", ""),
					resource.TestCheckResourceAttr(resourceVPCSubnetName, "ntp_addresses", "10.100.0.35,10.100.0.36"),
					resource.TestCheckResourceAttr(resourceVPCSubnetName, "tags.key", "value_update"),
				),
			},
		},
	})
}

func TestAccVpcSubnetV1Import(t *testing.T) {
	t.Parallel()
	qts := vpcSubnetQuotas()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetV1Import,
			},
			{
				ResourceName:      resourceVPCSubnetName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVpcSubnetV1Timeout(t *testing.T) {
	var subnet subnets.Subnet
	t.Parallel()
	qts := vpcSubnetQuotas()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetV1Timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetV1Exists(resourceVPCSubnetName, &subnet),
				),
			},
		},
	})
}

func TestAccVpcSubnetV1DnsList(t *testing.T) {
	var subnet subnets.Subnet
	t.Parallel()
	qts := vpcSubnetQuotas()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetV1DnsList,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetV1Exists(resourceVPCSubnetName, &subnet),
					resource.TestCheckResourceAttr(resourceVPCSubnetName, "dns_list.0", "100.125.4.25"),
					resource.TestCheckResourceAttr(resourceVPCSubnetName, "dns_list.1", "8.8.8.8"),
				),
			},
			{
				Config: testAccVpcSubnetV1DnsListUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetV1Exists(resourceVPCSubnetName, &subnet),
					resource.TestCheckResourceAttr(resourceVPCSubnetName, "dns_list.0", "100.125.4.25"),
					resource.TestCheckResourceAttr(resourceVPCSubnetName, "dns_list.1", "1.1.1.1"),
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
  name = "vpc_test_sn"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name              = "subnet_name"
  description       = "some description"
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
	testAccVpcSubnetV1Import = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test_imp"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name              = "subnet_name"
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
  name = "vpc_test_sn"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name              = "subnet_name_update"
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
  name = "vpc_test_t"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name              = "subnet_name"
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
  name = "vpc_name_dns"
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
	testAccVpcSubnetV1DnsListUpdate = `
resource "opentelekomcloud_vpc_v1" "vpc" {
  name = "vpc_name_dns"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_1" {
  name       = "subnet_name"
  vpc_id     = opentelekomcloud_vpc_v1.vpc.id
  cidr       = cidrsubnet(opentelekomcloud_vpc_v1.vpc.cidr, 8, 0)
  gateway_ip = cidrhost(cidrsubnet(opentelekomcloud_vpc_v1.vpc.cidr, 8, 0), 1)
  dns_list   = ["100.125.4.25", "1.1.1.1"]
}
`
)
