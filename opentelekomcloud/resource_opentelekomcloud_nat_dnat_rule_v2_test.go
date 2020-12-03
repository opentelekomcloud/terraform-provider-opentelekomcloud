// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file at
//     https://www.github.com/huaweicloud/magic-modules
//
// ----------------------------------------------------------------------------

package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNatDnat_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNatDnatDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatDnat_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatDnatExists(),
				),
			},
		},
	})
}

func testAccNatDnat_basic() string {
	return fmt.Sprintf(`
resource "opentelekomcloud_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}

resource "opentelekomcloud_networking_router_interface_v2" "int_1" {
  subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
  router_id = opentelekomcloud_networking_router_v2.router_1.id
}

resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {
}

resource "opentelekomcloud_nat_gateway_v2" "nat_gw" {
  name   = "nat_gw"
  description = "test for terraform"
  spec = "1"
  internal_network_id = opentelekomcloud_networking_network_v2.network_1.id
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  depends_on = [opentelekomcloud_networking_router_interface_v2.int_1]
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = opentelekomcloud_networking_network_v2.network_1.id
  }
 depends_on = [opentelekomcloud_networking_router_interface_v2.int_1]
}

resource "opentelekomcloud_nat_dnat_rule_v2" "dnat" {
  floating_ip_id = opentelekomcloud_networking_floatingip_v2.fip_1.id
  nat_gateway_id = opentelekomcloud_nat_gateway_v2.nat_gw.id
  private_ip = opentelekomcloud_compute_instance_v2.instance_1.network.0.fixed_ip_v4
  internal_service_port = 993
  protocol = "tcp"
  external_service_port = 242
  depends_on = [opentelekomcloud_compute_instance_v2.instance_1]
}
`, OS_AVAILABILITY_ZONE)
}

func testAccCheckNatDnatDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.sdkClient(OS_REGION_NAME, "nat", serviceProjectLevel)
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_nat_dnat_rule_v2" {
			continue
		}

		url, err := replaceVarsForTest(rs, "dnat_rules/{id}")
		if err != nil {
			return err
		}
		url = client.ServiceURL(url)

		_, err = client.Get(
			url, nil,
			&golangsdk.RequestOpts{MoreHeaders: map[string]string{"Accept": "application/json"}})
		if err == nil {
			return fmt.Errorf("opentelekomcloud_nat_dnat_rule_v2 still exists at %s", url)
		}
	}

	return nil
}

func testAccCheckNatDnatExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		client, err := config.sdkClient(OS_REGION_NAME, "nat", serviceProjectLevel)
		if err != nil {
			return fmt.Errorf("Error creating sdk client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources["opentelekomcloud_nat_dnat_rule_v2.dnat"]
		if !ok {
			return fmt.Errorf("Error checking opentelekomcloud_nat_dnat_rule_v2.dnat exist, err=not found opentelekomcloud_nat_dnat_rule_v2.dnat")
		}

		url, err := replaceVarsForTest(rs, "dnat_rules/{id}")
		if err != nil {
			return fmt.Errorf("Error checking opentelekomcloud_nat_dnat_rule_v2.dnat exist, err=building url failed: %s", err)
		}
		url = client.ServiceURL(url)

		_, err = client.Get(
			url, nil,
			&golangsdk.RequestOpts{MoreHeaders: map[string]string{"Accept": "application/json"}})
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmt.Errorf("opentelekomcloud_nat_dnat_rule_v2.dnat is not exist")
			}
			return fmt.Errorf("Error checking opentelekomcloud_nat_dnat_rule_v2.dnat exist, err=send request failed: %s", err)
		}
		return nil
	}
}
