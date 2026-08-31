package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/vpcs"
	VpcV3 "github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v3/vpcs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceVPCName = "opentelekomcloud_vpc_v1.vpc_1"

func getVpcV1ResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.VpcV1Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud VPC v1 client: %w", err)
	}

	return vpcs.Get(client, state.Primary.ID)
}

func TestAccVpcV1_basic(t *testing.T) {
	var vpc vpcs.Vpc
	t.Parallel()
	quotas.BookOne(t, quotas.Router)
	rc := common.InitResourceCheck(resourceVPCName, &vpc, getVpcV1ResourceFunc)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1Basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceVPCName, "name", "terraform_provider_test"),
					resource.TestCheckResourceAttr(resourceVPCName, "description", "simple description"),
					resource.TestCheckResourceAttr(resourceVPCName, "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceVPCName, "status", "OK"),
					resource.TestCheckResourceAttr(resourceVPCName, "shared", "true"),
					resource.TestCheckResourceAttr(resourceVPCName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttrSet(resourceVPCName, "tenant_id"),
					resource.TestCheckResourceAttrSet(resourceVPCName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceVPCName, "updated_at"),
					resource.TestCheckResourceAttr(resourceVPCName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceVPCName, "tags.key", "value"),
				),
			},
			{
				Config: testAccVpcV1Update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceVPCName, "name", "terraform_provider_test1"),
					resource.TestCheckResourceAttr(resourceVPCName, "description", "simple description updated"),
					resource.TestCheckResourceAttr(resourceVPCName, "shared", "false"),
					resource.TestCheckResourceAttr(resourceVPCName, "tags.key", "value_update"),
				),
			},
			{
				ResourceName:      resourceVPCName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVpcV1_secondaryCidr(t *testing.T) {
	var vpc vpcs.Vpc
	t.Parallel()
	quotas.BookOne(t, quotas.Router)
	rc := common.InitResourceCheck(resourceVPCName, &vpc, getVpcV1ResourceFunc)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV3BasicCidr,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
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
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceVPCName, "name", "tf_acc_test_v3"),
					resource.TestCheckResourceAttr(resourceVPCName, "description", "simple description updated"),
					resource.TestCheckResourceAttr(resourceVPCName, "secondary_cidr", "23.8.0.0/16"),
					resource.TestCheckResourceAttr(resourceVPCName, "shared", "false"),
					resource.TestCheckResourceAttr(resourceVPCName, "tags.key", "value_update"),
				),
			},
			{
				Config: testAccVpcV3RemoveCidr,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceVPCName, "name", "tf_acc_test_v3"),
					resource.TestCheckResourceAttr(resourceVPCName, "description", "simple description updated"),
					resource.TestCheckResourceAttr(resourceVPCName, "secondary_cidr", ""),
					resource.TestCheckResourceAttr(resourceVPCName, "shared", "false"),
					resource.TestCheckResourceAttr(resourceVPCName, "tags.key", "value_update"),
				),
			},
		},
	})
}

func TestAccVpcV1_readIgnoresExternalCidrs(t *testing.T) {
	const injectedCidr = "23.9.0.0/16"

	var vpc vpcs.Vpc
	t.Parallel()
	quotas.BookOne(t, quotas.Router)
	rc := common.InitResourceCheck(resourceVPCName, &vpc, getVpcV1ResourceFunc)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1ReadIsolation,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					injectSecondaryCidr(resourceVPCName, injectedCidr),
				),
			},
			{
				Config:   testAccVpcV1ReadIsolation,
				PlanOnly: true,
			},
		},
	})
}

func injectSecondaryCidr(resourceName, cidr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}
		client, err := common.TestAccProvider.Meta().(*cfg.Config).NetworkingV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating NetworkingV3 client: %w", err)
		}
		_, err = VpcV3.AddSecondaryCidr(client, rs.Primary.ID, VpcV3.CidrOpts{
			Vpc: &VpcV3.AddExtendCidrOption{ExtendCidrs: []string{cidr}},
		})
		if err != nil {
			return fmt.Errorf("error injecting secondary CIDR %s on VPC %s: %w", cidr, rs.Primary.ID, err)
		}
		return nil
	}
}

func TestAccVpcV1_timeout(t *testing.T) {
	var vpc vpcs.Vpc
	t.Parallel()
	quotas.BookOne(t, quotas.Router)
	rc := common.InitResourceCheck(resourceVPCName, &vpc, getVpcV1ResourceFunc)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1Timeout,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
		},
	})
}

const testAccVpcV1Basic = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name                  = "terraform_provider_test"
  description           = "simple description"
  cidr                  = "192.168.0.0/16"
  shared                = true
  enterprise_project_id = "0"

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

const testAccVpcV3RemoveCidr = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name        = "tf_acc_test_v3"
  description = "simple description updated"
  cidr        = "192.168.0.0/16"
  shared      = false

  tags = {
    foo = "bar"
    key = "value_update"
  }
}
`

const testAccVpcV1ReadIsolation = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "tf_acc_test_read_isolation"
  cidr = "192.168.0.0/16"
}
`
