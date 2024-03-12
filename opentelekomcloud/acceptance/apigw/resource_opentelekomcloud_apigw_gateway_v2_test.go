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

const resourceName = "opentelekomcloud_apigw_gateway_v2.gateway"

func TestAccAPIGWv2Gateway_basic(t *testing.T) {
	var gatewayConfig gateway.Gateway
	name := fmt.Sprintf("gateway-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAPIGWv2GatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGWv2GatewayBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGWv2GatewayExists(resourceName, &gatewayConfig),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "spec_id", "BASIC"),
					resource.TestCheckResourceAttr(resourceName, "description", "test gateway"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_size", "5"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "22:00:00"),
				),
			},
			{
				Config: testAccAPIGWv2GatewayUpdated(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGWv2GatewayExists(resourceName, &gatewayConfig),
					resource.TestCheckResourceAttr(resourceName, "name", name+"-updated"),
					resource.TestCheckResourceAttr(resourceName, "description", "test gateway 2"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_size", "0"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "02:00:00"),
				),
			},
		},
	})
}

func testAccCheckAPIGWv2GatewayDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.APIGWV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud ApiGateway v2 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_apigw_gateway_v2" {
			continue
		}

		_, err := gateway.Get(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("API Gateway configuration still exists")
		}
	}

	return nil
}

func testAccCheckAPIGWv2GatewayExists(n string, configuration *gateway.Gateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.APIGWV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud APIGateway v2 client: %w", err)
		}

		found, err := gateway.Get(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("ApiGateway Configuration not found")
		}
		configuration = found

		return nil
	}
}

func TestAccAPIGWGatewayV2ImportBasic(t *testing.T) {
	name := fmt.Sprintf("gateway-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAPIGWv2GatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGWv2GatewayBasic(name),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"ingress_bandwidth_size",
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
  name                    = "%s"
  spec_id                 = "BASIC"
  vpc_id                  = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id               = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  security_group_id       = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zones      = ["eu-de-01", "eu-de-02"]
  description             = "test gateway"
  bandwidth_size          = 5
  maintain_begin          =  "22:00:00"
  ingress_bandwidth_size  = 5
}
`, common.DataSourceSubnet, common.DataSourceSecGroupDefault, gatewayName)
}

func testAccAPIGWv2GatewayUpdated(gatewayName string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_apigw_gateway_v2" "gateway"{
  name 					  = "%s-updated"
  spec_id 				  = "BASIC"
  vpc_id 				  = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id               = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  security_group_id       = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zones      = ["eu-de-01", "eu-de-02"]
  description             = "test gateway 2"
  bandwidth_size          = 0
  maintain_begin          = "02:00:00"
  ingress_bandwidth_size  = 5
}
`, common.DataSourceSubnet, common.DataSourceSecGroupDefault, gatewayName)
}
