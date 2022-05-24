package acceptance

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/vpnaas/siteconnections"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceSiteConnectionName = "opentelekomcloud_vpnaas_site_connection_v2.conn_1"

func TestAccVpnSiteConnectionV2_basic(t *testing.T) {
	var conn siteconnections.Connection
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckSiteConnectionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteConnectionV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteConnectionV2Exists(resourceSiteConnectionName, &conn),
					resource.TestCheckResourceAttrPtr(resourceSiteConnectionName, "ikepolicy_id", &conn.IKEPolicyID),
					resource.TestCheckResourceAttr(resourceSiteConnectionName, "admin_state_up", strconv.FormatBool(conn.AdminStateUp)),
					resource.TestCheckResourceAttr(resourceSiteConnectionName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceSiteConnectionName, "tags.key", "value"),
					resource.TestCheckResourceAttrPtr(resourceSiteConnectionName, "ipsecpolicy_id", &conn.IPSecPolicyID),
					resource.TestCheckResourceAttrPtr(resourceSiteConnectionName, "vpnservice_id", &conn.VPNServiceID),
					resource.TestCheckResourceAttrPtr(resourceSiteConnectionName, "local_ep_group_id", &conn.LocalEPGroupID),
					resource.TestCheckResourceAttrPtr(resourceSiteConnectionName, "local_id", &conn.LocalID),
					resource.TestCheckResourceAttrPtr(resourceSiteConnectionName, "peer_ep_group_id", &conn.PeerEPGroupID),
					resource.TestCheckResourceAttrPtr(resourceSiteConnectionName, "name", &conn.Name),
				),
			},
		},
	})
}

func testAccCheckSiteConnectionV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpnaas_site_connection_v2" {
			continue
		}
		_, err = siteconnections.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("site connection (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSiteConnectionV2Exists(n string, conn *siteconnections.Connection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.NetworkingV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
		}

		var found *siteconnections.Connection

		found, err = siteconnections.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		*conn = *found

		return nil
	}
}

var testAccSiteConnectionV2Basic = fmt.Sprintf(`
%s

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "tf_test_network"
  admin_state_up = true
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  network_id = opentelekomcloud_networking_network_v2.network_1.id
  cidr       = "192.168.199.0/24"
  ip_version = 4
}

resource "opentelekomcloud_networking_router_v2" "router_1" {
  name             = "my_router"
  external_gateway = data.opentelekomcloud_networking_network_v2.ext_network.id
}

resource "opentelekomcloud_networking_router_interface_v2" "router_interface_1" {
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
}

resource "opentelekomcloud_vpnaas_service_v2" "service_1" {
  router_id      = opentelekomcloud_networking_router_v2.router_1.id
  admin_state_up = false
}

resource "opentelekomcloud_vpnaas_ipsec_policy_v2" "policy_1" { }

resource "opentelekomcloud_vpnaas_ike_policy_v2" "policy_2" { }

resource "opentelekomcloud_vpnaas_endpoint_group_v2" "group_1" {
  type      = "cidr"
  endpoints = ["10.2.0.0/24", "10.3.0.0/24"]
}
resource "opentelekomcloud_vpnaas_endpoint_group_v2" "group_2" {
  type      = "subnet"
  endpoints = [opentelekomcloud_networking_subnet_v2.subnet_1.id]
}

resource "opentelekomcloud_vpnaas_site_connection_v2" "conn_1" {
  name              = "connection_1"
  ikepolicy_id      = opentelekomcloud_vpnaas_ike_policy_v2.policy_2.id
  ipsecpolicy_id    = opentelekomcloud_vpnaas_ipsec_policy_v2.policy_1.id
  vpnservice_id     = opentelekomcloud_vpnaas_service_v2.service_1.id
  psk               = "secret"
  peer_address      = "192.168.10.1"
  peer_id           = "192.168.10.1"
  local_ep_group_id = opentelekomcloud_vpnaas_endpoint_group_v2.group_2.id
  peer_ep_group_id  = opentelekomcloud_vpnaas_endpoint_group_v2.group_1.id

  tags = {
    foo = "bar"
    key = "value"
  }

  depends_on = ["opentelekomcloud_networking_router_interface_v2.router_interface_1"]
}
`, common.DataSourceExtNetwork)
