package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceDnatRuleName = "opentelekomcloud_nat_dnat_rule_v2.dnat"

func TestAccNatDnat_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acc.TestAccPreCheck(t) },
		ProviderFactories: acc.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNatDnatDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatDnatBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatDnatExists(),
				),
			},
		},
	})
}

func TestAccNatDnatRule_withPort(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acc.TestAccPreCheck(t) },
		ProviderFactories: acc.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNatDnatDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatDnatRuleWithPort,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatDnatExists(),
				),
			},
		},
	})
}

var testAccNatDnatRuleWithPort = fmt.Sprintf(`
%s

resource "opentelekomcloud_networking_port_v2" "this" {
  name       = "test"
  network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  fixed_ip {
    subnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.subnet_id
  }
}

resource "opentelekomcloud_nat_gateway_v2" "this" {
  name                = "test"
  spec                = "1"
  router_id           = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  internal_network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
}

resource "opentelekomcloud_networking_floatingip_v2" "eip" {}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  security_groups   = ["default"]
  availability_zone = "%s"
  network {
    port = opentelekomcloud_networking_port_v2.this.id
  }
}

resource "opentelekomcloud_nat_dnat_rule_v2" "dnat" {
  floating_ip_id        = opentelekomcloud_networking_floatingip_v2.eip.id
  nat_gateway_id        = opentelekomcloud_nat_gateway_v2.this.id
  external_service_port = 80
  protocol              = "tcp"
  port_id               = opentelekomcloud_networking_port_v2.this.id
  internal_service_port = 80
  depends_on            = [opentelekomcloud_compute_instance_v2.instance_1]
}
`, acc.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccNatDnatBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}

resource "opentelekomcloud_nat_gateway_v2" "nat_gw" {
  name                = "nat_gw"
  description         = "test for terraform"
  spec                = "1"
  internal_network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  router_id           = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  security_groups   = ["default"]
  availability_zone = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}

resource "opentelekomcloud_nat_dnat_rule_v2" "dnat" {
  floating_ip_id        = opentelekomcloud_networking_floatingip_v2.fip_1.id
  nat_gateway_id        = opentelekomcloud_nat_gateway_v2.nat_gw.id
  private_ip            = opentelekomcloud_compute_instance_v2.instance_1.network.0.fixed_ip_v4
  internal_service_port = 993
  protocol              = "tcp"
  external_service_port = 242
}
`, acc.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

func testAccCheckNatDnatDestroy(s *terraform.State) error {
	config := acc.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.NatV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating sdk client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_nat_dnat_rule_v2" {
			continue
		}

		url, err := common.ReplaceVarsForTest(rs, "dnat_rules/{id}")
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
		config := acc.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.NatV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating sdk client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources["opentelekomcloud_nat_dnat_rule_v2.dnat"]
		if !ok {
			return fmt.Errorf("error checking opentelekomcloud_nat_dnat_rule_v2.dnat exist, err=not found opentelekomcloud_nat_dnat_rule_v2.dnat")
		}

		url, err := common.ReplaceVarsForTest(rs, "dnat_rules/{id}")
		if err != nil {
			return fmt.Errorf("error checking opentelekomcloud_nat_dnat_rule_v2.dnat exist, err=building url failed: %s", err)
		}
		url = client.ServiceURL(url)

		_, err = client.Get(
			url, nil,
			&golangsdk.RequestOpts{MoreHeaders: map[string]string{"Accept": "application/json"}})
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmt.Errorf("opentelekomcloud_nat_dnat_rule_v2.dnat is not exist")
			}
			return fmt.Errorf("error checking opentelekomcloud_nat_dnat_rule_v2.dnat exist, err=send request failed: %s", err)
		}
		return nil
	}
}
