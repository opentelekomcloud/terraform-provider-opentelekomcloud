package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	VpcV3 "github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v3/vpcs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceVpcSecondaryCidrV3Name = "opentelekomcloud_vpc_secondary_cidr_v3.test"

func getVpcSecondaryCidrFunc(c *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.NetworkingV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud NetworkingV3 client: %w", err)
	}
	vpc, err := VpcV3.Get(client, state.Primary.ID)
	if err != nil {
		return nil, err
	}
	if len(vpc.SecondaryCidrs) == 0 {
		return nil, fmt.Errorf("vpc %s has no secondary CIDRs", state.Primary.ID)
	}
	return vpc, nil
}

func TestAccVpcSecondaryCidrV3_basic(t *testing.T) {
	t.Parallel()
	quotas.BookOne(t, quotas.Router)
	vpcName := tools.RandomString("tf-acc-sec-cidr-", 5)

	var vpc VpcV3.Vpc
	rc := common.InitResourceCheck(resourceVpcSecondaryCidrV3Name, &vpc, getVpcSecondaryCidrFunc)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSecondaryCidrV3OneCidr(vpcName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceVpcSecondaryCidrV3Name, "cidrs.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceVpcSecondaryCidrV3Name, "cidrs.*", "23.9.0.0/16"),
				),
			},
			{
				Config: testAccVpcSecondaryCidrV3ThreeCidrs(vpcName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceVpcSecondaryCidrV3Name, "cidrs.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceVpcSecondaryCidrV3Name, "cidrs.*", "23.9.0.0/16"),
					resource.TestCheckTypeSetElemAttr(resourceVpcSecondaryCidrV3Name, "cidrs.*", "23.10.0.0/16"),
					resource.TestCheckTypeSetElemAttr(resourceVpcSecondaryCidrV3Name, "cidrs.*", "23.11.0.0/16"),
				),
			},
			{
				Config: testAccVpcSecondaryCidrV3TwoCidrs(vpcName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceVpcSecondaryCidrV3Name, "cidrs.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceVpcSecondaryCidrV3Name, "cidrs.*", "23.9.0.0/16"),
					resource.TestCheckTypeSetElemAttr(resourceVpcSecondaryCidrV3Name, "cidrs.*", "23.11.0.0/16"),
				),
			},
			{
				Config: testAccVpcSecondaryCidrV3MaxCidrs(vpcName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceVpcSecondaryCidrV3Name, "cidrs.#", "5"),
				),
			},
			{
				ResourceName:      resourceVpcSecondaryCidrV3Name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccVpcSecondaryCidrV3OneCidr(vpcName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_secondary_cidr_v3" "test" {
  vpc_id = opentelekomcloud_vpc_v1.test.id
  cidrs  = ["23.9.0.0/16"]
}
`, vpcName)
}

func testAccVpcSecondaryCidrV3ThreeCidrs(vpcName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_secondary_cidr_v3" "test" {
  vpc_id = opentelekomcloud_vpc_v1.test.id
  cidrs = [
    "23.9.0.0/16",
    "23.10.0.0/16",
    "23.11.0.0/16",
  ]
}
`, vpcName)
}

func testAccVpcSecondaryCidrV3TwoCidrs(vpcName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_secondary_cidr_v3" "test" {
  vpc_id = opentelekomcloud_vpc_v1.test.id
  cidrs = [
    "23.9.0.0/16",
    "23.11.0.0/16",
  ]
}
`, vpcName)
}

func testAccVpcSecondaryCidrV3MaxCidrs(vpcName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_secondary_cidr_v3" "test" {
  vpc_id = opentelekomcloud_vpc_v1.test.id
  cidrs = [
    "23.9.0.0/16",
    "23.10.0.0/16",
    "23.11.0.0/16",
    "23.12.0.0/16",
    "23.13.0.0/16",
  ]
}
`, vpcName)
}
