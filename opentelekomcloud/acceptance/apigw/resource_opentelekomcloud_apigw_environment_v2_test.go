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

func TestAccAPIGWv2Environment_basic(t *testing.T) {
	var envConfig env.EnvResp
	nameEnv := fmt.Sprintf("environment_%s", acctest.RandString(10))
	nameGateway := fmt.Sprintf("gateway-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGWv2EnvironmentBasic(nameGateway, nameEnv),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGWv2EnvironmentExists(resourceNameEnvironment, &envConfig),
					resource.TestCheckResourceAttr(resourceNameEnvironment, "name", nameEnv),
					resource.TestCheckResourceAttr(resourceNameEnvironment, "description", "test description"),
				),
			},
			{
				Config: testAccAPIGWv2EnvironmentUpdated(nameGateway, nameEnv),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGWv2EnvironmentExists(resourceNameEnvironment, &envConfig),
					resource.TestCheckResourceAttr(resourceNameEnvironment, "name", nameEnv+"_updated"),
					resource.TestCheckResourceAttr(resourceNameEnvironment, "description", "test description updated"),
				),
			},
		},
	})
}

func testAccCheckAPIGWv2EnvironmentExists(n string, configuration *env.EnvResp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		rsgw, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.APIGWV2Client(accenv.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud APIGW v2 client: %w", err)
		}

		found, err := env.List(client, env.ListOpts{
			GatewayID: rsgw.Primary.ID,
			Name:      rs.Primary.Attributes["name"],
		})
		if err != nil {
			return err
		}

		if found[0].ID != rs.Primary.ID {
			return fmt.Errorf("APIGW environment not found")
		}
		configuration = &found[0]

		return nil
	}
}

func TestAccAPIGWEnvironmentV2ImportBasic(t *testing.T) {
	nameGateway := fmt.Sprintf("gateway-%s", acctest.RandString(10))
	nameEnv := fmt.Sprintf("environment_%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGWv2EnvironmentBasic(nameGateway, nameEnv),
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
