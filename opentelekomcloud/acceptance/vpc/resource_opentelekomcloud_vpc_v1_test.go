package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/vpcs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceVPCName = "opentelekomcloud_vpc_v1.vpc_1"

func TestAccVpcV1_basic(t *testing.T) {
	var vpc vpcs.Vpc
	t.Parallel()
	th.AssertNoErr(t, quotas.Router.Acquire())
	defer quotas.Router.Release()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckOTCVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCVpcV1Exists(resourceVPCName, &vpc),
					resource.TestCheckResourceAttr(resourceVPCName, "name", "terraform_provider_test"),
					resource.TestCheckResourceAttr(resourceVPCName, "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceVPCName, "status", "OK"),
					resource.TestCheckResourceAttr(resourceVPCName, "shared", "true"),
					resource.TestCheckResourceAttr(resourceVPCName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceVPCName, "tags.key", "value"),
				),
			},
			{
				Config: testAccVpcV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCVpcV1Exists(resourceVPCName, &vpc),
					resource.TestCheckResourceAttr(resourceVPCName, "name", "terraform_provider_test1"),
					resource.TestCheckResourceAttr(resourceVPCName, "shared", "false"),
					resource.TestCheckResourceAttr(resourceVPCName, "tags.key", "value_update"),
				),
			},
		},
	})
}

func TestAccVpcV1_timeout(t *testing.T) {
	var vpc vpcs.Vpc
	t.Parallel()
	th.AssertNoErr(t, quotas.Router.Acquire())
	defer quotas.Router.Release()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckOTCVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1Timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCVpcV1Exists(resourceVPCName, &vpc),
				),
			},
		},
	})
}

func TestAccVpcV1_import(t *testing.T) {
	t.Parallel()
	th.AssertNoErr(t, quotas.Router.Acquire())
	defer quotas.Router.Release()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckOTCVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1Import,
			},
			{
				ResourceName:      resourceVPCName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckOTCVpcV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	vpcClient, err := config.NetworkingV1Client(env.OS_REGION_NAME)
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
		vpcClient, err := config.NetworkingV1Client(env.OS_REGION_NAME)
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

const testAccVpcV1Basic = `
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

const testAccVpcV1Update = `
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

const testAccVpcV1Timeout = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "terraform_provider_test-t"
  cidr = "192.168.0.0/16"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

const testAccVpcV1Import = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name   = "terraform_provider_test-imp"
  cidr   = "192.168.0.0/16"
  shared = true

  tags = {
    foo = "bar"
    key = "value"
  }
}
`
