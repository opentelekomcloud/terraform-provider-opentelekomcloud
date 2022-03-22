package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/networks"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/subnets"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/vpc"
)

const resourceNwNetworkName = "opentelekomcloud_networking_network_v2.network_1"

func TestAccNetworkingV2Network_basic(t *testing.T) {
	var network networks.Network

	t.Parallel()
	quotas.BookOne(t, quotas.Network)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkBasic,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2NetworkExists(resourceNwNetworkName, &network),
				),
			},
			{
				Config: testAccNetworkingV2NetworkUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNwNetworkName, "name", "network_2"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkImport,
			},

			{
				ResourceName:      resourceNwNetworkName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkingV2Network_netstack(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var router routers.Router

	t.Parallel()
	qts := vpcSubnetQuotas()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkNetstack,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2NetworkExists(resourceNwNetworkName, &network),
					TestAccCheckNetworkingV2SubnetExists(resourceNwSubnetName, &subnet),
					TestAccCheckNetworkingV2RouterExists(resourceRouterName, &router),
					TestAccCheckNetworkingV2RouterInterfaceExists("opentelekomcloud_networking_router_interface_v2.ri_1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_timeout(t *testing.T) {
	var network networks.Network
	t.Parallel()
	quotas.BookOne(t, quotas.Network)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkTimeout,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2NetworkExists(resourceNwNetworkName, &network),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_multipleSegmentMappings(t *testing.T) {
	var network networks.Network
	t.Parallel()
	quotas.BookOne(t, quotas.Network)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkMultipleSegmentMappings,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckNetworkingV2NetworkExists(resourceNwNetworkName, &network),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2NetworkDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_networking_network_v2" {
			continue
		}

		_, id := vpc.ExtractValFromNid(rs.Primary.ID)
		_, err := networks.Get(networkingClient, id).Extract()
		if err == nil {
			return fmt.Errorf("network still exists")
		}
	}

	return nil
}

const testAccNetworkingV2NetworkBasic = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}
`

const testAccNetworkingV2NetworkImport = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_1_imp"
  admin_state_up = "true"
}
`

const testAccNetworkingV2NetworkUpdate = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_2"
  # Can't do this to a network on OTC
  #admin_state_up = "false"
}
`

const testAccNetworkingV2NetworkNetstack = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  cidr       = "192.168.10.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}

resource "opentelekomcloud_networking_router_v2" "router_1" {
  name = "router_1"
}

resource "opentelekomcloud_networking_router_interface_v2" "ri_1" {
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
}
`

const testAccNetworkingV2NetworkTimeout = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

const testAccNetworkingV2NetworkMultipleSegmentMappings = `
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1"
  segments = [
    {
      segmentation_id = 2
      network_type    = "vxlan"
    }
  ]
  admin_state_up = "true"
}
`
