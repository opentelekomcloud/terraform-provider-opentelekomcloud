package dcaas

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	virtualinterface "github.com/opentelekomcloud/gophertelekomcloud/openstack/dcaas/v2/virtual-interface"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const vi = "opentelekomcloud_dc_virtual_interface_v2.vi_1"

func TestDirectConnectVirtualInterfaceV2Resource_basic(t *testing.T) {
	dcId := os.Getenv("OS_DIRECT_CONNECT_ID")
	if dcId == "" {
		t.Skip("OS_DIRECT_CONNECT_ID should be set for acceptance tests")
	}
	intName := fmt.Sprintf("dc-%s", acctest.RandString(5))
	var vint virtualinterface.VirtualInterface
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDirectConnectVirtualInterfaceV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualInterfaceV2_basic(intName, dcId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDirectConnectVirtualInterfaceV2Exists(vi, &vint),
					resource.TestCheckResourceAttr(vi, "name", intName),
					resource.TestCheckResourceAttr(vi, "description", "description"),
					resource.TestCheckResourceAttr(vi, "bandwidth", "5"),
				),
			},
			{
				Config: testAccVirtualInterfaceV2_update(intName+"updated", dcId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDirectConnectVirtualInterfaceV2Exists(vi, &vint),
					resource.TestCheckResourceAttr(vi, "name", intName+"updated"),
					resource.TestCheckResourceAttr(vi, "description", "description updated"),
					resource.TestCheckResourceAttr(vi, "bandwidth", "10"),
				),
			},
			{
				ResourceName:      vi,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDirectConnectVirtualInterfaceV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.DCaaSV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating DCaaS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_dc_virtual_interface_v2" {
			continue
		}

		_, err := virtualinterface.Get(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("virtual interface still exists")
		}
		var errDefault404 golangsdk.ErrDefault404
		if !errors.As(err, &errDefault404) {
			return err
		}
	}
	return nil
}

func testAccCheckDirectConnectVirtualInterfaceV2Exists(n string, gateway *virtualinterface.VirtualInterface) resource.TestCheckFunc {
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

		found, err := virtualinterface.Get(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("DCaaS Virtual Interface not found")
		}

		*gateway = *found

		return nil
	}
}

func testAccVirtualInterfaceV2_basic(name, dcId string) string {
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
  name              = "tf_acc_vgw_1"
  description       = "acc test"
  local_ep_group_id = opentelekomcloud_dc_endpoint_group_v2.dc_endpoint_group.id
}

resource "opentelekomcloud_dc_virtual_interface_v2" "int_1" {
  direct_connect_id = "%s"
  vgw_id            = opentelekomcloud_dc_virtual_gateway_v2.vgw_1.id
  name              = "%s"
  description       = "description"
  type              = "private"
  route_mode        = "static"
  vlan              = 100
  bandwidth         = 5

  remote_ep_group_id   = opentelekomcloud_dc_endpoint_group_v2.dc_endpoint_group.id
  local_gateway_v4_ip  = "180.1.1.1/24"
  remote_gateway_v4_ip = "180.1.1.2/24"
}
`, common.DataSourceSubnet, common.DataSourceProject, dcId, name)
}

func testAccVirtualInterfaceV2_update(name, dcId string) string {
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
  name              = "tf_acc_vgw_1"
  description       = "acc test updated"
  local_ep_group_id = opentelekomcloud_dc_endpoint_group_v2.dc_endpoint_group.id
}

resource "opentelekomcloud_dc_virtual_interface_v2" "int_1" {
  direct_connect_id = "%s"
  vgw_id            = opentelekomcloud_dc_virtual_gateway_v2.vgw_1.id
  name              = "%s"
  description       = "description updated"
  type              = "private"
  route_mode        = "static"
  vlan              = 100
  bandwidth         = 10

  remote_ep_group_id   = opentelekomcloud_dc_endpoint_group_v2.dc_endpoint_group.id
  local_gateway_v4_ip  = "180.1.1.1/24"
  remote_gateway_v4_ip = "180.1.1.2/24"
}
`, common.DataSourceSubnet, common.DataSourceProject, dcId, name)
}
