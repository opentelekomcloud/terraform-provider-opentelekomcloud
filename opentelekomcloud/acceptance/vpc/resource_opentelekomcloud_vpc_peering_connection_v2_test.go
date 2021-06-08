package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/peerings"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccOTCVpcPeeringConnectionV2_basic(t *testing.T) {
	var peering peerings.Peering

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckOTCVpcPeeringConnectionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOTCVpcPeeringConnectionV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCVpcPeeringConnectionV2Exists("opentelekomcloud_vpc_peering_connection_v2.peering_1", &peering),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_peering_connection_v2.peering_1", "name", "opentelekomcloud_peering"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_peering_connection_v2.peering_1", "status", "ACTIVE"),
				),
			},
			{
				Config: testAccOTCVpcPeeringConnectionV2_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_peering_connection_v2.peering_1", "name", "opentelekomcloud_peering_1"),
				),
			},
		},
	})
}

func TestAccOTCVpcPeeringConnectionV2_timeout(t *testing.T) {
	var peering peerings.Peering

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckOTCVpcPeeringConnectionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOTCVpcPeeringConnectionV2_timeout,
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
			return fmt.Errorf("Vpc Peering Connection still exists")
		}
	}

	return nil
}

func testAccCheckOTCVpcPeeringConnectionV2Exists(n string, peering *peerings.Peering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
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
			return fmt.Errorf("Vpc peering Connection not found")
		}

		*peering = *found

		return nil
	}
}

const testAccOTCVpcPeeringConnectionV2_basic = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "vpc_test1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering_1" {
  name = "opentelekomcloud_peering"
  vpc_id = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}
`
const testAccOTCVpcPeeringConnectionV2_update = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "vpc_test1"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering_1" {
  name = "opentelekomcloud_peering_1"
  vpc_id = opentelekomcloud_vpc_v1.vpc_1.id
  peer_vpc_id = opentelekomcloud_vpc_v1.vpc_2.id
}
`
const testAccOTCVpcPeeringConnectionV2_timeout = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
  name = "vpc_test1"
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
