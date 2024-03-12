package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/gateway"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/group"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v1/configurations"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceGroupName = "opentelekomcloud_apigw_group_v2.gateway"

func TestAccAPIGWv2Group_basic(t *testing.T) {
	var groupConfig group.GroupResp
	name := fmt.Sprintf("group-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckAPIGWv2GroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGWv2GroupBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGWv2GroupExists(resourceGroupName, &groupConfig),
					resource.TestCheckResourceAttr(resourceGroupName, "name", name),
					resource.TestCheckResourceAttr(resourceGroupName, "spec_id", "BASIC"),
					resource.TestCheckResourceAttr(resourceGroupName, "description", "test gateway"),
					resource.TestCheckResourceAttr(resourceGroupName, "bandwidth_size", "5"),
					resource.TestCheckResourceAttr(resourceGroupName, "maintain_begin", "22:00:00"),
				),
			},
			{
				Config: testAccAPIGWv2GroupUpdated(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGWv2GroupExists(resourceGroupName, &groupConfig),
					resource.TestCheckResourceAttr(resourceGroupName, "name", name+"-updated"),
					resource.TestCheckResourceAttr(resourceGroupName, "description", "test gateway 2"),
					resource.TestCheckResourceAttr(resourceGroupName, "bandwidth_size", "0"),
					resource.TestCheckResourceAttr(resourceGroupName, "maintain_begin", "02:00:00"),
				),
			},
		},
	})
}

func testAccCheckAPIGWv2GroupDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.AutoscalingV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud AutoScalingV1 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_apigw_gateway_v2" {
			continue
		}

		_, err := configurations.Get(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("AS configuration still exists")
		}
	}

	return nil
}

func testAccCheckAPIGWv2GroupExists(n string, configuration *group.GroupResp) resource.TestCheckFunc {
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
			return fmt.Errorf("error creating OpenTelekomCloud AutoScalingV1 client: %w", err)
		}

		found, err := gateway.Get(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("autoscaling Configuration not found")
		}
		configuration = found

		return nil
	}
}

func TestAccAPIGWGroupV2ImportBasic(t *testing.T) {
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

func testAccAPIGWv2GroupBasic(gatewayName string) string {
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

func testAccAPIGWv2GroupUpdated(gatewayName string) string {
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
