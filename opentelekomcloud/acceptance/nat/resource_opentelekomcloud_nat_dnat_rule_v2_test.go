package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/dnatrules"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceDnatRuleName = "opentelekomcloud_nat_dnat_rule_v2.dnat"

func TestAccNatDnat_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNatDnatDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatDnatBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatDnatExists(resourceDnatRuleName),
					resource.TestCheckResourceAttr(resourceDnatRuleName, "internal_service_port", "993"),
					resource.TestCheckResourceAttr(resourceDnatRuleName, "external_service_port", "242"),
					resource.TestCheckResourceAttr(resourceDnatRuleName, "protocol", "tcp"),
					resource.TestCheckResourceAttrSet(resourceDnatRuleName, "private_ip"),
				),
			},
		},
	})
}

func TestAccNatDnatRule_withPort(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNatDnatDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatDnatRuleWithPort,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatDnatExists(resourceDnatRuleName),
					resource.TestCheckResourceAttr(resourceDnatRuleName, "internal_service_port", "80"),
					resource.TestCheckResourceAttr(resourceDnatRuleName, "protocol", "tcp"),
					resource.TestCheckResourceAttrSet(resourceDnatRuleName, "port_id"),
				),
			},
		},
	})
}

func TestAccNatDnat_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNatDnatDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatDnatBasic,
			},

			{
				ResourceName:      resourceDnatRuleName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNatDnatDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.NatV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NATv2 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_nat_dnat_rule_v2" {
			continue
		}

		_, err := dnatrules.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("dnat rule still exists")
		}
	}

	return nil
}

func testAccCheckNatDnatExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.NatV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud NATv2 client: %w", err)
		}

		found, err := dnatrules.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("dnat rule not found")
		}

		return nil
	}
}

var testAccNatDnatRuleWithPort = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_networking_port_v2" "this" {
  name       = "dnat_rule_port"
  network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  fixed_ip {
    subnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.subnet_id
  }
}

resource "opentelekomcloud_nat_gateway_v2" "this" {
  name                = "dnat_rule_gw"
  spec                = "1"
  router_id           = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  internal_network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
}

resource "opentelekomcloud_networking_floatingip_v2" "eip" {}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  security_groups   = ["default"]
  availability_zone = "%s"
  image_id          = data.opentelekomcloud_images_image_v2.latest_image.id

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
`, common.DataSourceImage, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)

var testAccNatDnatBasic = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}

resource "opentelekomcloud_nat_gateway_v2" "nat_gw" {
  name                = "dnat_rule_gw"
  description         = "test for terraform"
  spec                = "1"
  internal_network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  router_id           = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  security_groups   = ["default"]
  availability_zone = "%s"
  image_id          = data.opentelekomcloud_images_image_v2.latest_image.id

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
`, common.DataSourceImage, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)
