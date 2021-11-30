package acceptance

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/ports"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceNetworkingVIPName = "opentelekomcloud_networking_vip_v2.vip_1"

// TestAccNetworkingV2VIP_basic is basic acc test.
func TestAccNetworkingV2VIP_basic(t *testing.T) {
	var vip ports.Port
	t.Parallel()
	qts := vpcSubnetQuotas()
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2VIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2VIPConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2VIPExists(resourceNetworkingVIPName, &vip),
				),
			},
		},
	})
}

// testAccCheckNetworkingV2VIPDestroy checks destroy.
func testAccCheckNetworkingV2VIPDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_networking_vip_v2" {
			continue
		}

		_, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("VIP still exists")
		}
	}

	log.Printf("[DEBUG] testAccCheckNetworkingV2VIPDestroy success!")

	return nil
}

// testAccCheckNetworkingV2VIPExists checks exist.
func testAccCheckNetworkingV2VIPExists(n string, vip *ports.Port) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		networkingClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
		}

		found, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("VIP not found")
		}
		log.Printf("[DEBUG] test found is: %#v", found)
		*vip = *found

		return nil
	}
}

// testAccNetworkingV2VIPConfigBasic is used to create.
var testAccNetworkingV2VIPConfigBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1_vip"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1_vip"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}

resource "opentelekomcloud_networking_router_interface_v2" "router_interface_1" {
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
}

resource "opentelekomcloud_networking_router_v2" "router_1" {
  name = "router_1_vip"
  external_gateway = data.opentelekomcloud_networking_network_v2.ext_network.id
}

resource "opentelekomcloud_networking_vip_v2" "vip_1" {
  network_id = opentelekomcloud_networking_network_v2.network_1.id
  subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
}
`, common.DataSourceExtNetwork)
