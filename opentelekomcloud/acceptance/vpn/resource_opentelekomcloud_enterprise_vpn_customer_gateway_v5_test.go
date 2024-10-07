package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	cgw "github.com/opentelekomcloud/gophertelekomcloud/openstack/evpn/v5/customer-gateway"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceEvpnCustomerGatewayName = "opentelekomcloud_enterprise_vpn_customer_gateway_v5.cgw_1"

func getCustomerGatewayResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.EvpnV5Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud EVPN v5 client: %s", err)
	}
	return cgw.Get(client, state.Primary.ID)
}

func TestAccCustomerGateway_basic(t *testing.T) {
	var gw cgw.CustomerGateway
	name := fmt.Sprintf("evpn_acc_cgw_%s", acctest.RandString(5))
	updateName := fmt.Sprintf("evpn_acc_cgw_up_%s", acctest.RandString(5))

	rc := common.InitResourceCheck(
		resourceEvpnCustomerGatewayName,
		&gw,
		getCustomerGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCustomerGateway_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceEvpnCustomerGatewayName, "name", name),
					resource.TestCheckResourceAttr(resourceEvpnCustomerGatewayName, "id_value", "10.1.2.10"),
					resource.TestCheckResourceAttr(resourceEvpnCustomerGatewayName, "tags.key", "val"),
					resource.TestCheckResourceAttr(resourceEvpnCustomerGatewayName, "tags.foo", "bar"),
				),
			},
			{
				Config: testCustomerGateway_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceEvpnCustomerGatewayName, "name", updateName),
					resource.TestCheckResourceAttr(resourceEvpnCustomerGatewayName, "id_value", "10.1.2.10"),
					resource.TestCheckResourceAttr(resourceEvpnCustomerGatewayName, "tags.key", "val"),
					resource.TestCheckResourceAttr(resourceEvpnCustomerGatewayName, "tags.foo", "bar-update"),
				),
			},
			{
				ResourceName:      resourceEvpnCustomerGatewayName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCustomerGateway_basic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_enterprise_vpn_customer_gateway_v5" "cgw_1" {
  name     = "%s"
  id_value = "10.1.2.10"

  tags = {
    key = "val"
    foo = "bar"
  }
}`, name)
}

func testCustomerGateway_update(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_enterprise_vpn_customer_gateway_v5" "cgw_1" {
  name     = "%s"
  id_value = "10.1.2.10"

  tags = {
    key = "val"
    foo = "bar-update"
  }
}`, name)
}
