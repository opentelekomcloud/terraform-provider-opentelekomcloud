package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
	quotas.BookOne(t, quotas.Router)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1Exists(resourceVPCName, &vpc),
					resource.TestCheckResourceAttr(resourceVPCName, "name", "terraform_provider_test"),
					resource.TestCheckResourceAttr(resourceVPCName, "description", "simple description"),
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
					testAccCheckVpcV1Exists(resourceVPCName, &vpc),
					resource.TestCheckResourceAttr(resourceVPCName, "name", "terraform_provider_test1"),
					resource.TestCheckResourceAttr(resourceVPCName, "description", "simple description updated"),
					resource.TestCheckResourceAttr(resourceVPCName, "shared", "false"),
					resource.TestCheckResourceAttr(resourceVPCName, "tags.key", "value_update"),
				),
			},
		},
	})
}

func TestAccVpcV1_secondaryCidr(t *testing.T) {
	var vpc vpcs.Vpc
	t.Parallel()
	quotas.BookOne(t, quotas.Router)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV3BasicCidr,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1Exists(resourceVPCName, &vpc),
					resource.TestCheckResourceAttr(resourceVPCName, "name", "tf_acc_test_v3"),
					resource.TestCheckResourceAttr(resourceVPCName, "description", "simple description"),
					resource.TestCheckResourceAttr(resourceVPCName, "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceVPCName, "secondary_cidr", "23.9.0.0/16"),
					resource.TestCheckResourceAttr(resourceVPCName, "status", "OK"),
					resource.TestCheckResourceAttr(resourceVPCName, "shared", "true"),
					resource.TestCheckResourceAttr(resourceVPCName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceVPCName, "tags.key", "value"),
				),
			},
			{
				Config: testAccVpcV3UpdateCidr,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1Exists(resourceVPCName, &vpc),
					resource.TestCheckResourceAttr(resourceVPCName, "name", "tf_acc_test_v3"),
					resource.TestCheckResourceAttr(resourceVPCName, "description", "simple description updated"),
					resource.TestCheckResourceAttr(resourceVPCName, "secondary_cidr", "23.8.0.0/16"),
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
	quotas.BookOne(t, quotas.Router)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1Timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1Exists(resourceVPCName, &vpc),
				),
			},
		},
	})
}

func TestAccVpcV1_import(t *testing.T) {
	t.Parallel()
	quotas.BookOne(t, quotas.Router)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcV1Destroy,
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

func testAccCheckVpcV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.NetworkingV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV1 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpc_v1" {
			continue
		}

		_, err := vpcs.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("vpc still exists")
		}
	}

	return nil
}

func testAccCheckVpcV1Exists(n string, vpc *vpcs.Vpc) resource.TestCheckFunc {
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

		found, err := vpcs.Get(client, rs.Primary.ID).Extract()
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
  name        = "terraform_provider_test"
  description = "simple description"
  cidr        = "192.168.0.0/16"
  shared      = true

  tags = {
    foo = "bar"
    key = "value"
  }
}
`

const testAccVpcV1Update = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name        = "terraform_provider_test1"
  description = "simple description updated"
  cidr        = "192.168.0.0/16"
  shared      = false

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

const testAccVpcV3BasicCidr = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name           = "tf_acc_test_v3"
  description    = "simple description"
  cidr           = "192.168.0.0/16"
  secondary_cidr = "23.9.0.0/16"
  shared         = true

  tags = {
    foo = "bar"
    key = "value"
  }
}
`

const testAccVpcV3UpdateCidr = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name           = "tf_acc_test_v3"
  description    = "simple description updated"
  cidr           = "192.168.0.0/16"
  secondary_cidr = "23.8.0.0/16"
  shared         = false

  tags = {
    foo = "bar"
    key = "value_update"
  }
}
`
