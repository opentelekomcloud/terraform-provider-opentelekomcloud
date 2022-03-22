package acceptance

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/fwaas_v2/firewall_groups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/fw"
)

func TestAccFWFirewallGroupV2_basic(t *testing.T) {
	var epolicyID *string
	var ipolicyID *string

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWFirewallGroupV2Basic1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2("opentelekomcloud_fw_firewall_group_v2.fw_1", "", "", ipolicyID, epolicyID),
				),
			},
			{
				Config: testAccFWFirewallGroupV2Basic2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2(
						"opentelekomcloud_fw_firewall_group_v2.fw_1", "fw_1", "terraform acceptance test", ipolicyID, epolicyID),
				),
			},
		},
	})
}

func TestAccFWFirewallGroupV2_port0(t *testing.T) {
	var firewallGroup fw.FirewallGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWFirewallV2Port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("opentelekomcloud_fw_firewall_group_v2.fw_1", &firewallGroup),
					testAccCheckFWFirewallPortCount(&firewallGroup, 1),
				),
			},
		},
	})
}

func TestAccFWFirewallGroupV2_no_ports(t *testing.T) {
	var firewallGroup fw.FirewallGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWFirewallV2NoPorts,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("opentelekomcloud_fw_firewall_group_v2.fw_1", &firewallGroup),
					resource.TestCheckResourceAttr("opentelekomcloud_fw_firewall_group_v2.fw_1", "description", "firewall router test"),
					testAccCheckFWFirewallPortCount(&firewallGroup, 0),
				),
			},
		},
	})
}

func TestAccFWFirewallGroupV2_port_update(t *testing.T) {
	var firewallGroup fw.FirewallGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWFirewallV2Port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("opentelekomcloud_fw_firewall_group_v2.fw_1", &firewallGroup),
					testAccCheckFWFirewallPortCount(&firewallGroup, 1),
				),
			},
			{
				Config: testAccFWFirewallV2PortAdd,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("opentelekomcloud_fw_firewall_group_v2.fw_1", &firewallGroup),
					testAccCheckFWFirewallPortCount(&firewallGroup, 2),
				),
			},
		},
	})
}

func TestAccFWFirewallGroupV2_port_remove(t *testing.T) {
	var firewallGroup fw.FirewallGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWFirewallV2Port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("opentelekomcloud_fw_firewall_group_v2.fw_1", &firewallGroup),
					testAccCheckFWFirewallPortCount(&firewallGroup, 1),
				),
			},
			{
				Config: testAccFWFirewallV2PortRemove,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWFirewallGroupV2Exists("opentelekomcloud_fw_firewall_group_v2.fw_1", &firewallGroup),
					testAccCheckFWFirewallPortCount(&firewallGroup, 0),
				),
			},
		},
	})
}

func testAccCheckFWFirewallGroupV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_firewall_group" {
			continue
		}

		_, err = firewall_groups.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("firewall group (%s) still exists", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckFWFirewallGroupV2Exists(n string, firewallGroup *fw.FirewallGroup) resource.TestCheckFunc {
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
			return fmt.Errorf("exists) Error creating OpenTelekomCloud networking client: %s", err)
		}

		var found fw.FirewallGroup
		err = firewall_groups.Get(networkingClient, rs.Primary.ID).ExtractInto(&found)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("firewall group not found")
		}

		*firewallGroup = found

		return nil
	}
}

func testAccCheckFWFirewallPortCount(firewallGroup *fw.FirewallGroup, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(firewallGroup.PortIDs) != expected {
			return fmt.Errorf("expected %d Ports, got %d", expected, len(firewallGroup.PortIDs))
		}

		return nil
	}
}

func testAccCheckFWFirewallGroupV2(n, expectedName, expectedDescription string, ipolicyID *string, epolicyID *string) resource.TestCheckFunc {
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
			return fmt.Errorf("exists) Error creating OpenTelekomCloud networking client: %s", err)
		}

		var found *firewall_groups.FirewallGroup
		for i := 0; i < 5; i++ {
			// Firewall creation is asynchronous. Retry some times
			// if we get a 404 error. Fail on any other error.
			found, err = firewall_groups.Get(networkingClient, rs.Primary.ID).Extract()
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					time.Sleep(time.Second)
					continue
				}
				return err
			}
			break
		}

		switch {
		case found.Name != expectedName:
			err = fmt.Errorf("expected Name to be <%s> but found <%s>", expectedName, found.Name)
		case found.Description != expectedDescription:
			err = fmt.Errorf("expected Description to be <%s> but found <%s>",
				expectedDescription, found.Description)
		case found.IngressPolicyID == "":
			err = fmt.Errorf("ingress Policy should not be empty")
		case found.EgressPolicyID == "":
			err = fmt.Errorf("egress Policy should not be empty")
		case ipolicyID != nil && found.IngressPolicyID == *ipolicyID:
			err = fmt.Errorf("ingress Policy had not been correctly updated. Went from <%s> to <%s>",
				expectedName, found.Name)
		case epolicyID != nil && found.EgressPolicyID == *epolicyID:
			err = fmt.Errorf("egress Policy had not been correctly updated. Went from <%s> to <%s>",
				expectedName, found.Name)
		}

		if err != nil {
			return err
		}

		ipolicyID = &found.IngressPolicyID
		epolicyID = &found.EgressPolicyID

		return nil
	}
}

const testAccFWFirewallGroupV2Basic1 = `
resource "opentelekomcloud_fw_firewall_group_v2" "fw_1" {
  ingress_policy_id = opentelekomcloud_fw_policy_v2.policy_1.id
  egress_policy_id  = opentelekomcloud_fw_policy_v2.policy_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "opentelekomcloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
}
`

const testAccFWFirewallGroupV2Basic2 = `
resource "opentelekomcloud_fw_firewall_group_v2" "fw_1" {
  name              = "fw_1"
  description       = "terraform acceptance test"
  ingress_policy_id = opentelekomcloud_fw_policy_v2.policy_2.id
  egress_policy_id  = opentelekomcloud_fw_policy_v2.policy_2.id
  admin_state_up    = true
}

resource "opentelekomcloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
}

resource "opentelekomcloud_fw_policy_v2" "policy_2" {
  name = "policy_2"
}
`

var testAccFWFirewallV2Port = fmt.Sprintf(`
%s

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name        = "subnet_1"
  cidr        = "192.168.199.0/24"
  ip_version  = 4
  enable_dhcp = true
  network_id  = opentelekomcloud_networking_network_v2.network_1.id
}

resource "opentelekomcloud_networking_router_v2" "router_1" {
  name             = "router_1"
  admin_state_up   = "true"
  external_gateway = data.opentelekomcloud_networking_network_v2.ext_network.id
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  name           = "port_1"
  admin_state_up = "true"
  network_id     = opentelekomcloud_networking_network_v2.network_1.id

  fixed_ip {
    subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
    #ip_address = "192.168.199.23"
  }
}

resource "opentelekomcloud_networking_router_interface_v2" "router_interface_1" {
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  port_id   = opentelekomcloud_networking_port_v2.port_1.id
}

resource "opentelekomcloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
}

resource "opentelekomcloud_fw_firewall_group_v2" "fw_1" {
  name              = "firewall_1"
  description       = "firewall router test"
  ingress_policy_id = opentelekomcloud_fw_policy_v2.policy_1.id
  #egress_policy_id = opentelekomcloud_fw_policy_v2.policy_1.id
  ports = [
    opentelekomcloud_networking_port_v2.port_1.id
  ]
  depends_on = ["opentelekomcloud_networking_router_interface_v2.router_interface_1"]
}
`, common.DataSourceExtNetwork)

var testAccFWFirewallV2PortAdd = fmt.Sprintf(`
%s

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}

resource "opentelekomcloud_networking_network_v2" "network_2" {
  name           = "network_2"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_2" {
  name       = "subnet_2"
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_2.id
}

resource "opentelekomcloud_networking_router_v2" "router_1" {
  name             = "router_1"
  admin_state_up   = "true"
  external_gateway = data.opentelekomcloud_networking_network_v2.ext_network.id
}

resource "opentelekomcloud_networking_router_v2" "router_2" {
  name             = "router_2"
  admin_state_up   = "true"
  external_gateway = data.opentelekomcloud_networking_network_v2.ext_network.id
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  name           = "port_1"
  admin_state_up = "true"
  network_id     = opentelekomcloud_networking_network_v2.network_1.id

  fixed_ip {
    subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
    #ip_address = "192.168.199.23"
  }
}

resource "opentelekomcloud_networking_port_v2" "port_2" {
  name           = "port_2"
  admin_state_up = "true"
  network_id     = opentelekomcloud_networking_network_v2.network_2.id

  fixed_ip {
    subnet_id = opentelekomcloud_networking_subnet_v2.subnet_2.id
    #ip_address = "192.168.199.24"
  }
}

resource "opentelekomcloud_networking_router_interface_v2" "router_interface_1" {
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  port_id   = opentelekomcloud_networking_port_v2.port_1.id
}

resource "opentelekomcloud_networking_router_interface_v2" "router_interface_2" {
  router_id = opentelekomcloud_networking_router_v2.router_2.id
  port_id   = opentelekomcloud_networking_port_v2.port_2.id
}

resource "opentelekomcloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
}

resource "opentelekomcloud_fw_firewall_group_v2" "fw_1" {
  name              = "firewall_1"
  description       = "firewall router test"
  ingress_policy_id = opentelekomcloud_fw_policy_v2.policy_1.id
  egress_policy_id  = opentelekomcloud_fw_policy_v2.policy_1.id
  ports = [
    opentelekomcloud_networking_port_v2.port_1.id,
    opentelekomcloud_networking_port_v2.port_2.id
  ]
  depends_on = ["opentelekomcloud_networking_router_interface_v2.router_interface_1", "opentelekomcloud_networking_router_interface_v2.router_interface_2"]
}
`, common.DataSourceExtNetwork)

const testAccFWFirewallV2PortRemove = `
resource "opentelekomcloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
}

resource "opentelekomcloud_fw_firewall_group_v2" "fw_1" {
  name              = "firewall_1"
  description       = "firewall router test"
  ingress_policy_id = opentelekomcloud_fw_policy_v2.policy_1.id
  egress_policy_id  = opentelekomcloud_fw_policy_v2.policy_1.id
}
`

const testAccFWFirewallV2NoPorts = `
resource "opentelekomcloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
}

resource "opentelekomcloud_fw_firewall_group_v2" "fw_1" {
  name              = "firewall_1"
  description       = "firewall router test"
  ingress_policy_id = opentelekomcloud_fw_policy_v2.policy_1.id
  egress_policy_id  = opentelekomcloud_fw_policy_v2.policy_1.id
}
`
