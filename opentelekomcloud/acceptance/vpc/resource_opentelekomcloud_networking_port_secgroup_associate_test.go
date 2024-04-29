package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/ports"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourcePortAssociateName = "opentelekomcloud_networking_port_secgroup_associate_v2.associate"

func getPortResourceFunc(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Networking v2 client: %s", err)
	}
	return ports.Get(client, state.Primary.Attributes["port_id"]).Extract()
}

func TestAccNetworkingV2PortAssociate_basic(t *testing.T) {
	var port ports.Port
	rc := common.InitResourceCheck(
		resourcePortAssociateName,
		&port,
		getPortResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPortAssociate_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					testAccCheckNetworkingV2PortSecGroupAssociateCountSecurityGroups(&port, 2),
				),
			},
			{
				ResourceName:            resourcePortAssociateName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "security_group_ids"},
			},
		},
	})
}

func testAccCheckNetworkingV2PortSecGroupAssociateCountSecurityGroups(port *ports.Port, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(port.SecurityGroups) != expected {
			return fmt.Errorf("expected %d Security Groups, got %d", expected, len(port.SecurityGroups))
		}

		return nil
	}
}

const testAccNetworkingV2PortSecGroupAssociate = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "acc_network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "opentelekomcloud_networking_secgroup_v2" "secgroup_2" {
  name        = "secgroup_2"
  description = "terraform security group acceptance test"
}

resource "opentelekomcloud_networking_port_v2" "port" {
  name           = "port_1"
  network_id     = opentelekomcloud_networking_network_v2.network_1.id
  admin_state_up = "true"
}
`

func testAccPortAssociate_basic() string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_networking_port_secgroup_associate_v2" "associate" {
  port_id = opentelekomcloud_networking_port_v2.port.id
  force   = "false"
  security_group_ids = [
    opentelekomcloud_networking_secgroup_v2.secgroup_1.id,
  ]
}
`, testAccNetworkingV2PortSecGroupAssociate)
}
