package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/peerings"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceVPCPeeringName = "opentelekomcloud_vpc_peering_connection_v2.peering_1"

func TestAccVpcPeeringConnectionV2_basic(t *testing.T) {
	var peering peerings.Peering
	t.Parallel()
	quotas.BookMany(t, multipleRouters(2))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckOTCVpcPeeringConnectionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConnectionV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCVpcPeeringConnectionV2Exists(resourceVPCPeeringName, &peering),
					resource.TestCheckResourceAttr(resourceVPCPeeringName, "name", "opentelekomcloud_peering"),
					resource.TestCheckResourceAttr(resourceVPCPeeringName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccVpcPeeringConnectionV2Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVPCPeeringName, "name", "opentelekomcloud_peering_1"),
				),
			},
		},
	})
}

func TestAccVpcPeeringConnectionV2_import(t *testing.T) {
	t.Parallel()
	quotas.BookMany(t, multipleRouters(2))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckOTCVpcPeeringConnectionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConnectionV2Import,
			},
			{
				ResourceName:      resourceVPCPeeringName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVpcPeeringConnectionV2_timeout(t *testing.T) {
	var peering peerings.Peering
	t.Parallel()
	quotas.BookMany(t, multipleRouters(2))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckOTCVpcPeeringConnectionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConnectionV2Timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCVpcPeeringConnectionV2Exists("opentelekomcloud_vpc_peering_connection_v2.peering_1", &peering),
				),
			},
		},
	})
}

func testAccCheckOTCVpcPeeringConnectionV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	peeringClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud Peering client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpc_peering_connection_v2" {
			continue
		}

		_, err := peerings.Get(peeringClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("vpc Peering Connection still exists")
		}
	}

	return nil
}

func testAccCheckOTCVpcPeeringConnectionV2Exists(n string, peering *peerings.Peering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		peeringClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud Peering client: %s", err)
		}

		found, err := peerings.Get(peeringClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("vpc peering Connection not found")
		}

		*peering = *found

		return nil
	}
}

const testAccVpcPeeringConnectionV2Basic = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test_p"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "vpc_test_p1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering_1" {
  name = "opentelekomcloud_peering"
  vpc_id = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}
`
const testAccVpcPeeringConnectionV2Import = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test_p_imp"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "vpc_test_p_imp1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering_1" {
  name = "opentelekomcloud_peering_imp"
  vpc_id = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}
`
const testAccVpcPeeringConnectionV2Update = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test_p"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "vpc_test_p1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering_1" {
  name = "opentelekomcloud_peering_1"
  vpc_id = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}
`
const testAccVpcPeeringConnectionV2Timeout = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test_p_t"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "vpc_test_p_t1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering_1" {
  name = "opentelekomcloud_peering"
  vpc_id = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id

 timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
