package acceptance

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/ports"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/vpc"
)

// TestAccNetworkingV2VIPAssociate_basic is basic acc test.
func TestAccNetworkingV2VIPAssociate_basic(t *testing.T) {
	t.Skip("this test produces dangling resources")
	var vip ports.Port
	var port1 ports.Port
	var port2 ports.Port
	t.Parallel()
	qts := vpcSubnetQuotas()
	qts = append(qts,
		&quotas.ExpectedQuota{Q: quotas.Volume, Count: 2},
		&quotas.ExpectedQuota{Q: quotas.VolumeSize, Count: 4 + 4},
		&quotas.ExpectedQuota{Q: quotas.Server, Count: 2},
	)
	quotas.BookMany(t, qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2VIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccNetworkingV2VIPAssociateConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					// testAccCheckNetworkingV2PortExists("opentelekomcloud_networking_port_v2.port_1", &port1),
					// testAccCheckNetworkingV2PortExists("opentelekomcloud_networking_port_v2.port_2", &port2),
					testAccCheckNetworkingV2VIPExists("opentelekomcloud_networking_vip_v2.vip_1", &vip),
					testAccCheckNetworkingV2VIPAssociateAssociated(&port1, &vip),
					testAccCheckNetworkingV2VIPAssociateAssociated(&port2, &vip),
				),
			},
		},
	})
}

// testAccCheckNetworkingV2VIPAssociateDestroy checks destroy.
func testAccCheckNetworkingV2VIPAssociateDestroy(s *terraform.State) error { // nolint:unused
	config := common.TestAccProvider.Meta().(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_networking_vip_associate_v2" {
			continue
		}

		vipid, portids, err := vpc.ParseNetworkingVIPAssociateID(rs.Primary.ID)
		if err != nil {
			return err
		}

		vipport, err := ports.Get(networkingClient, vipid).Extract()
		if err != nil {
			// If the error is a 404, then the vip port does not exist,
			// and therefore the floating IP cannot be associated to it.
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}

		// port by port
		for _, portid := range portids {
			p, err := ports.Get(networkingClient, portid).Extract()
			if err != nil {
				// If the error is a 404, then the port does not exist,
				// and therefore the floating IP cannot be associated to it.
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return nil
				}
				return err
			}

			// But if the port and vip still exists
			for _, ip := range p.FixedIPs {
				for _, addresspair := range vipport.AllowedAddressPairs {
					if ip.IPAddress == addresspair.IPAddress {
						return fmt.Errorf("VIP %s is still associated to port %s", vipid, portid)
					}
				}
			}
		}
	}

	log.Printf("[DEBUG] testAccCheckNetworkingV2VIPAssociateDestroy success!")
	return nil
}

func testAccCheckNetworkingV2VIPAssociateAssociated(p *ports.Port, vip *ports.Port) resource.TestCheckFunc { // nolint:unused
	return func(s *terraform.State) error {
		config := common.TestAccProvider.Meta().(*cfg.Config)
		networkingClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
		}

		p, err := ports.Get(networkingClient, p.ID).Extract()
		if err != nil {
			// If the error is a 404, then the port does not exist,
			// and therefore the VIP cannot be associated to it.
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}

		vipport, err := ports.Get(networkingClient, vip.ID).Extract()
		if err != nil {
			// If the error is a 404, then the vip port does not exist,
			// and therefore the VIP cannot be associated to it.
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}

		for _, ip := range p.FixedIPs {
			for _, addresspair := range vipport.AllowedAddressPairs {
				if ip.IPAddress == addresspair.IPAddress {
					log.Printf("[DEBUG] testAccCheckNetworkingV2VIPAssociateAssociated success!")
					return nil
				}
			}
		}

		return fmt.Errorf("VIP %s was not attached to port %s", vipport.ID, p.ID)
	}
}

// TestAccNetworkingV2VIPAssociateConfigBasic is used to create.
var TestAccNetworkingV2VIPAssociateConfigBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_vip_ass_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name       = "subnet_vip_ass_1"
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}

resource "opentelekomcloud_networking_router_interface_v2" "router_interface_1" {
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
}

resource "opentelekomcloud_networking_router_v2" "router_1" {
  name             = "router_vip_ass_1"
  external_gateway = data.opentelekomcloud_networking_network_v2.ext_network.id
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  name           = "port_1"
  admin_state_up = "true"
  network_id     = opentelekomcloud_networking_network_v2.network_1.id
  fixed_ip {
    subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
  }
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_vip_ass_1"
  security_groups = ["default"]

  network {
    port = opentelekomcloud_networking_port_v2.port_1.id
  }
}

resource "opentelekomcloud_networking_port_v2" "port_2" {
  name           = "port_2"
  admin_state_up = "true"
  network_id     = opentelekomcloud_networking_network_v2.network_1.id
  fixed_ip {
    subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
  }
}

resource "opentelekomcloud_compute_instance_v2" "instance_2" {
  name            = "instance_vip_ass_2"
  security_groups = ["default"]

  network {
    port = opentelekomcloud_networking_port_v2.port_2.id
  }
}

resource "opentelekomcloud_networking_vip_v2" "vip_1" {
  network_id = opentelekomcloud_networking_network_v2.network_1.id
  subnet_id  = opentelekomcloud_networking_subnet_v2.subnet_1.id
}

resource "opentelekomcloud_networking_vip_associate_v2" "vip_associate_1" {
  vip_id   = opentelekomcloud_networking_vip_v2.vip_1.id
  port_ids = [opentelekomcloud_networking_port_v2.port_1.id, opentelekomcloud_networking_port_v2.port_2.id]
}
`, common.DataSourceExtNetwork)
