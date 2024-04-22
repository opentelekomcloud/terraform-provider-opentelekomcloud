package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	accenv "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceNameEnvironment = "opentelekomcloud_apigw_environment_v2.env"

func getEnvironmentFunc(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.APIGWV2Client(accenv.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	found, err := env.List(client, env.ListOpts{
		GatewayID: state.Primary.Attributes["instance_id"],
		Name:      state.Primary.Attributes["name"],
	})
	return found[0], err
}

func TestAccAPIGWv2Environment_basic(t *testing.T) {
	var envConfig env.EnvResp
	nameEnv := fmt.Sprintf("environment_%s", acctest.RandString(10))
	nameGateway := fmt.Sprintf("gateway-%s", acctest.RandString(10))

	rc := common.InitResourceCheck(
		resourceNameEnvironment,
		&envConfig,
		getEnvironmentFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGWv2EnvironmentBasic(nameGateway, nameEnv),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceNameEnvironment, "name", nameEnv),
					resource.TestCheckResourceAttr(resourceNameEnvironment, "description", "test description"),
				),
			},
			{
				Config: testAccAPIGWv2EnvironmentUpdated(nameGateway, nameEnv),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceNameEnvironment, "name", nameEnv+"_updated"),
					resource.TestCheckResourceAttr(resourceNameEnvironment, "description", "test description updated"),
				),
			},
			{
				ResourceName:      resourceNameEnvironment,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAPIGWv2EnvironmentImportStateIdFunc(),
			},
		},
	})
}

func testAccAPIGWv2EnvironmentImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var gatewayID string
		var envName string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "opentelekomcloud_apigw_gateway_v2" {
				gatewayID = rs.Primary.ID
			} else if rs.Type == "opentelekomcloud_apigw_environment_v2" && rs.Primary.ID != "" {
				envName = rs.Primary.Attributes["name"]
			}
		}
		if gatewayID == "" || envName == "" {
			return "", fmt.Errorf("resource not found: %s/%s", gatewayID, envName)
		}
		return fmt.Sprintf("%s/%s", gatewayID, envName), nil
	}
}

func testAccAPIGWv2EnvironmentBasic(gatewayName, envName string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_apigw_environment_v2" "env"{
  name        = "%s"
  instance_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  description = "test description"
}
`, testAccAPIGWv2GatewayBasic(gatewayName), envName)
}

func testAccAPIGWv2EnvironmentUpdated(gatewayName, envName string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_apigw_environment_v2" "env"{
  name        = "%s_updated"
  instance_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  description = "test description updated"
}
`, testAccAPIGWv2GatewayBasic(gatewayName), envName)
}
