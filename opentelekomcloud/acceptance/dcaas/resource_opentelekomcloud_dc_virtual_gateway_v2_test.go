package dcaas

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	virtualgateway "github.com/opentelekomcloud/gophertelekomcloud/openstack/dcaas/v2/virtual-gateway"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const vg = "opentelekomcloud_dc_virtual_gateway_v2.vgw_1"

func TestDirectConnectVirtualGatewayV2Resource_basic(t *testing.T) {
	gwName := fmt.Sprintf("dc-%s", acctest.RandString(5))
	var gateway virtualgateway.VirtualGateway
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDirectConnectVirtualGatewayV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualGatewayV2_basic(gwName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDirectConnectVirtualGatewayV2Exists(vg, &gateway),
					resource.TestCheckResourceAttr(vg, "name", gwName),
					resource.TestCheckResourceAttr(vg, "description", "acc test"),
					resource.TestCheckResourceAttrSet(vg, "asn"),
					resource.TestCheckResourceAttrSet(vg, "status"),
					resource.TestCheckResourceAttrSet(vg, "local_ep_group_id"),
					resource.TestCheckResourceAttr(vg, "local_ep_group.0.type", "cidr"),
					resource.TestCheckResourceAttr(vg, "local_ep_group.0.endpoints.#", "1"),
				),
			},
			{
				Config: testAccVirtualGatewayV2_update(gwName + "updated"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDirectConnectVirtualGatewayV2Exists(vg, &gateway),
					resource.TestCheckResourceAttr(vg, "name", gwName+"updated"),
					resource.TestCheckResourceAttr(vg, "description", "acc test updated"),
					resource.TestCheckResourceAttrSet(vg, "asn"),
					resource.TestCheckResourceAttrSet(vg, "status"),
					resource.TestCheckResourceAttrSet(vg, "local_ep_group_id"),
					resource.TestCheckResourceAttr(vg, "local_ep_group.0.type", "cidr"),
					resource.TestCheckResourceAttr(vg, "local_ep_group.0.endpoints.#", "2"),
				),
			},
			{
				ResourceName:      vg,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestDirectConnectVirtualGatewayV2ResourceIpv6_combined(t *testing.T) {
	gwName := fmt.Sprintf("dc-%s", acctest.RandString(5))
	var gateway virtualgateway.VirtualGateway
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDirectConnectVirtualGatewayV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualGatewayV2Ipv6_basic(gwName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDirectConnectVirtualGatewayV2Exists(vg, &gateway),
					resource.TestCheckResourceAttr(vg, "name", gwName),
					resource.TestCheckResourceAttrSet(vg, "local_ep_group_ipv6_id"),
					resource.TestCheckResourceAttr(vg, "local_ep_group_v6.0.type", "cidr"),
					resource.TestCheckResourceAttr(vg, "local_ep_group_v6.0.endpoints.#", "1"),
				),
			},
			{
				Config: testAccVirtualGatewayV2Ipv6_update(gwName + "updated"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDirectConnectVirtualGatewayV2Exists(vg, &gateway),
					resource.TestCheckResourceAttr(vg, "name", gwName+"updated"),
					resource.TestCheckResourceAttrSet(vg, "local_ep_group_ipv6_id"),
					resource.TestCheckResourceAttr(vg, "local_ep_group_v6.0.type", "cidr"),
					resource.TestCheckResourceAttr(vg, "local_ep_group_v6.0.endpoints.#", "2"),
				),
			},
			{
				Config: testAccVirtualGatewayV2_combined(gwName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDirectConnectVirtualGatewayV2Exists(vg, &gateway),
					resource.TestCheckResourceAttr(vg, "name", gwName),
					resource.TestCheckResourceAttr(vg, "description", "acc test"),
					resource.TestCheckResourceAttrSet(vg, "asn"),
					resource.TestCheckResourceAttrSet(vg, "status"),
					resource.TestCheckResourceAttrSet(vg, "local_ep_group_id"),
					resource.TestCheckResourceAttrSet(vg, "local_ep_group_ipv6_id"),
					resource.TestCheckResourceAttr(vg, "local_ep_group.0.type", "cidr"),
					resource.TestCheckResourceAttr(vg, "local_ep_group.0.endpoints.#", "1"),
					resource.TestCheckResourceAttr(vg, "local_ep_group_v6.0.type", "cidr"),
					resource.TestCheckResourceAttr(vg, "local_ep_group_v6.0.endpoints.#", "1"),
				),
			},
			{
				Config: testAccVirtualGatewayV2_combinedUpdate(gwName + "updated"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDirectConnectVirtualGatewayV2Exists(vg, &gateway),
					resource.TestCheckResourceAttr(vg, "name", gwName+"updated"),
					resource.TestCheckResourceAttr(vg, "description", "acc test updated"),
					resource.TestCheckResourceAttrSet(vg, "asn"),
					resource.TestCheckResourceAttrSet(vg, "status"),
					resource.TestCheckResourceAttrSet(vg, "local_ep_group_id"),
					resource.TestCheckResourceAttrSet(vg, "local_ep_group_ipv6_id"),
					resource.TestCheckResourceAttr(vg, "local_ep_group.0.type", "cidr"),
					resource.TestCheckResourceAttr(vg, "local_ep_group.0.endpoints.#", "2"),
					resource.TestCheckResourceAttr(vg, "local_ep_group_v6.0.type", "cidr"),
					resource.TestCheckResourceAttr(vg, "local_ep_group_v6.0.endpoints.#", "2"),
				),
			},
		},
	})
}

func testAccCheckDirectConnectVirtualGatewayV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.DCaaSV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating DCaaS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_dc_virtual_gateway_v2" {
			continue
		}

		_, err := virtualgateway.Get(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("virtual gateway still exists")
		}
		var errDefault404 golangsdk.ErrDefault404
		if !errors.As(err, &errDefault404) {
			return err
		}
	}
	return nil
}

func testAccCheckDirectConnectVirtualGatewayV2Exists(n string, gateway *virtualgateway.VirtualGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.DCaaSV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating DCaaS client: %s", err)
		}

		found, err := virtualgateway.Get(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("DCaaS Virtual Gateway not found")
		}

		*gateway = *found

		return nil
	}
}

func testAccVirtualGatewayV2_basic(name string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_dc_virtual_gateway_v2" "vgw_1" {
  vpc_id      = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  name        = "%s"
  description = "acc test"
  local_ep_group {
    name        = "tf_acc_eg_1"
    endpoints   = ["10.2.0.0/24"]
    description = "first"
  }
}
`, common.DataSourceSubnet, common.DataSourceProject, name)
}

func testAccVirtualGatewayV2_update(name string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_dc_virtual_gateway_v2" "vgw_1" {
  vpc_id      = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  name        = "%s"
  description = "acc test updated"
  local_ep_group {
    name        = "tf_acc_eg_1"
    endpoints   = ["10.2.0.0/24", "10.3.0.0/24"]
    description = "first"
  }
}
`, common.DataSourceSubnet, common.DataSourceProject, name)
}

func testAccVirtualGatewayV2Ipv6_basic(name string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_dc_virtual_gateway_v2" "vgw_1" {
  vpc_id      = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  name        = "%s"
  description = "acc test"
  local_ep_group_v6 {
    name        = "tf_acc_eg_1"
    endpoints   = ["2a07:8700:2:4::/64"]
    description = "first"
  }
}
`, common.DataSourceSubnet, common.DataSourceProject, name)
}

func testAccVirtualGatewayV2Ipv6_update(name string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_dc_virtual_gateway_v2" "vgw_1" {
  vpc_id      = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  name        = "%s"
  description = "acc test updated"
  local_ep_group_v6 {
    name        = "tf_acc_eg_1"
    endpoints   = ["2a07:8700:2:4::/64", "2a07:8700:2:54::/64"]
    description = "first"
  }
}
`, common.DataSourceSubnet, common.DataSourceProject, name)
}

func testAccVirtualGatewayV2_combined(name string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_dc_virtual_gateway_v2" "vgw_1" {
  vpc_id      = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  name        = "%s"
  description = "acc test"
  local_ep_group {
    name        = "tf_acc_eg_1"
    endpoints   = ["10.2.0.0/24"]
    description = "first"
  }

  local_ep_group_v6 {
    name        = "tf_acc_eg_1"
    endpoints   = ["2a07:8700:2:4::/64"]
    description = "first"
  }
}
`, common.DataSourceSubnet, common.DataSourceProject, name)
}

func testAccVirtualGatewayV2_combinedUpdate(name string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_dc_virtual_gateway_v2" "vgw_1" {
  vpc_id      = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  name        = "%s"
  description = "acc test updated"
  local_ep_group {
    name        = "tf_acc_eg_1"
    endpoints   = ["10.2.0.0/24", "10.3.0.0/24"]
    description = "first"
  }
  local_ep_group_v6 {
    name        = "tf_acc_eg_1"
    endpoints   = ["2a07:8700:2:4::/64", "2a07:8700:2:54::/64"]
    description = "first"
  }
}
`, common.DataSourceSubnet, common.DataSourceProject, name)
}
