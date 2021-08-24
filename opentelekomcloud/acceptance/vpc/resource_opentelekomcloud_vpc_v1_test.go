package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/vpcs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccVpcV1_basic(t *testing.T) {
	var vpc vpcs.Vpc

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckOTCVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCVpcV1Exists("opentelekomcloud_vpc_v1.vpc_1", &vpc),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_v1.vpc_1", "name", "terraform_provider_test"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_v1.vpc_1", "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_v1.vpc_1", "status", "OK"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_v1.vpc_1", "shared", "true"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_v1.vpc_1", "tags.foo", "bar"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_v1.vpc_1", "tags.key", "value"),
				),
			},
		},
	})
}

func TestAccVpcV1_update(t *testing.T) {
	var vpc vpcs.Vpc

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckOTCVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCVpcV1Exists("opentelekomcloud_vpc_v1.vpc_1", &vpc),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_v1.vpc_1", "name", "terraform_provider_test"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_v1.vpc_1", "shared", "true"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_v1.vpc_1", "tags.key", "value"),
				),
			},
			{
				Config: testAccVpcV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCVpcV1Exists("opentelekomcloud_vpc_v1.vpc_1", &vpc),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_v1.vpc_1", "name", "terraform_provider_test1"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_v1.vpc_1", "shared", "false"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_v1.vpc_1", "tags.key", "value_update"),
				),
			},
		},
	})
}

func TestAccVpcV1_timeout(t *testing.T) {
	var vpc vpcs.Vpc

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckOTCVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCVpcV1Exists("opentelekomcloud_vpc_v1.vpc_1", &vpc),
				),
			},
		},
	})
}

func testAccCheckOTCVpcV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	vpcClient, err := config.NetworkingV1Client(env.OsRegionName)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud vpc client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpc_v1" {
			continue
		}

		_, err := vpcs.Get(vpcClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("vpc still exists")
		}
	}

	return nil
}

func testAccCheckOTCVpcV1Exists(n string, vpc *vpcs.Vpc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		vpcClient, err := config.NetworkingV1Client(env.OsRegionName)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud vpc client: %s", err)
		}

		found, err := vpcs.Get(vpcClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("vpc not found")
		}

		*vpc = *found

		return nil
	}
}

const testAccVpcV1_basic = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name   = "terraform_provider_test"
  cidr   = "192.168.0.0/16"
  shared = true

  tags = {
    foo = "bar"
    key = "value"
  }
}
`

const testAccVpcV1_update = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name   = "terraform_provider_test1"
  cidr   = "192.168.0.0/16"
  shared = false

  tags = {
    foo = "bar"
    key = "value_update"
  }
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
