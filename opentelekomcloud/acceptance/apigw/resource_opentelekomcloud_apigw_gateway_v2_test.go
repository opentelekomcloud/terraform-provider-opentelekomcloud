package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/gateway"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceGwName = "opentelekomcloud_apigw_gateway_v2.gateway"

func getGwFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.APIGWV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud APIG v2 client: %s", err)
	}
	return gateway.Get(c, state.Primary.ID)
}

func TestAccAPIGWv2Gateway_basic(t *testing.T) {
	var gatewayConfig gateway.Gateway
	name := fmt.Sprintf("gateway-%s", acctest.RandString(10))

	rc := common.InitResourceCheck(
		resourceGwName,
		&gatewayConfig,
		getGwFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGWv2GatewayBasic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceGwName, "name", name),
					resource.TestCheckResourceAttr(resourceGwName, "spec_id", "BASIC"),
					resource.TestCheckResourceAttr(resourceGwName, "description", "test gateway"),
					resource.TestCheckResourceAttr(resourceGwName, "maintain_begin", "22:00:00"),
				),
			},
			{
				Config: testAccAPIGWv2GatewayUpdated(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceGwName, "name", name+"-updated"),
					resource.TestCheckResourceAttr(resourceGwName, "description", "test gateway 2"),
					resource.TestCheckResourceAttr(resourceGwName, "bandwidth_size", "0"),
					resource.TestCheckResourceAttr(resourceGwName, "maintain_begin", "02:00:00"),
				),
			},
			{
				ResourceName:      resourceGwName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"ingress_bandwidth_size",
					"ingress_bandwidth_charging_mode",
				},
			},
		},
	})
}

func testAccAPIGWv2GatewayBasic(gatewayName string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_apigw_gateway_v2" "gateway"{
  name                    		  = "%s"
  spec_id                 		  = "BASIC"
  vpc_id                  		  = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id               		  = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  security_group_id       		  = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zones      		  = ["eu-de-01", "eu-de-02"]
  description             		  = "test gateway"
  ingress_bandwidth_size          = 5
  ingress_bandwidth_charging_mode = "bandwidth"
  maintain_begin                  = "22:00:00"
}
`, common.DataSourceSubnet, common.DataSourceSecGroupDefault, gatewayName)
}

func testAccAPIGWv2GatewayUpdated(gatewayName string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_apigw_gateway_v2" "gateway"{
  name 					          = "%s-updated"
  spec_id 				          = "BASIC"
  vpc_id 				          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id                       = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  security_group_id               = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zones              = ["eu-de-01", "eu-de-02"]
  description                     = "test gateway 2"
  ingress_bandwidth_size          = 5
  ingress_bandwidth_charging_mode = "bandwidth"
  maintain_begin                  = "02:00:00"
}
`, common.DataSourceSubnet, common.DataSourceSecGroupDefault, gatewayName)
}
