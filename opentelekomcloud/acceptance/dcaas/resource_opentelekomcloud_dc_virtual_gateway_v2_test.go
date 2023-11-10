package dcaas

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	dcaas "github.com/opentelekomcloud/gophertelekomcloud/openstack/dcaas/v2/direct-connect"
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

func testAccCheckDirectConnectVirtualGatewayV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.DCaaSV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating DCaaS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_direct_connect_virtual_gateway_v2" {
			continue
		}

		_, err := dcaas.Get(client, rs.Primary.ID)
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

resource "opentelekomcloud_dc_endpoint_group_v2" "dc_endpoint_group" {
  name        = "tf_acc_eg_1"
  type        = "cidr"
  endpoints   = ["10.2.0.0/24", "10.3.0.0/24"]
  description = "first"
  project_id  = data.opentelekomcloud_identity_project_v3.project.id
}

resource "opentelekomcloud_dc_virtual_gateway_v2" "vgw_1" {
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  name              = "%s"
  description       = "acc test"
  local_ep_group_id = opentelekomcloud_dc_endpoint_group_v2.dc_endpoint_group.id
}
`, common.DataSourceSubnet, common.DataSourceProject, name)
}

func testAccVirtualGatewayV2_update(name string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_dc_endpoint_group_v2" "dc_endpoint_group" {
  name        = "tf_acc_eg_1"
  type        = "cidr"
  endpoints   = ["10.2.0.0/24", "10.3.0.0/24"]
  description = "first"
  project_id  = data.opentelekomcloud_identity_project_v3.project.id
}

resource "opentelekomcloud_dc_endpoint_group_v2" "dc_endpoint_group_new" {
  name        = "tf_acc_eg_1"
  type        = "cidr"
  endpoints   = ["10.20.0.0/24", "10.30.0.0/24"]
  description = "first"
  project_id  = data.opentelekomcloud_identity_project_v3.project.id
}

resource "opentelekomcloud_dc_virtual_gateway_v2" "vgw_1" {
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  name              = "%s"
  description       = "acc test updated"
  local_ep_group_id = opentelekomcloud_dc_endpoint_group_v2.dc_endpoint_group_new.id
}
`, common.DataSourceSubnet, common.DataSourceProject, name)
}
