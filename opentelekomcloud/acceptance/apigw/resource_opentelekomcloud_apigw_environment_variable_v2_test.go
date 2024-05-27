package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	vars "github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/env_vars"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	accenv "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceNameEnvironmentVars = "opentelekomcloud_apigw_environment_variable_v2.var"

func getEnvironmentVariableFunc(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.APIGWV2Client(accenv.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return vars.Get(client, state.Primary.Attributes["gateway_id"], state.Primary.ID)
}

func TestAccEnvironmentVariable_basic(t *testing.T) {
	var variable vars.EnvVarsResp
	name := fmt.Sprintf("env_var_%s", acctest.RandString(5))

	rc := common.InitResourceCheck(
		resourceNameEnvironmentVars,
		&variable,
		getEnvironmentVariableFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentVariable_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceNameEnvironmentVars, "name", name),
					resource.TestCheckResourceAttr(resourceNameEnvironmentVars, "value", "/stage/demo"),
					resource.TestCheckResourceAttrPair(resourceNameEnvironmentVars, "gateway_id", "opentelekomcloud_apigw_gateway_v2.gateway", "id"),
					resource.TestCheckResourceAttrPair(resourceNameEnvironmentVars, "group_id", "opentelekomcloud_apigw_group_v2.group", "id"),
					resource.TestCheckResourceAttrPair(resourceNameEnvironmentVars, "environment_id", "opentelekomcloud_apigw_environment_v2.env", "id"),
				),
			},
			{
				ResourceName:      resourceNameEnvironmentVars,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccEnvironmentVariableImportStateFunc(),
			},
		},
	})
}

func testAccEnvironmentVariableImportStateFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceNameEnvironmentVars]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceNameEnvironmentVars, rs)
		}

		gatewayId := rs.Primary.Attributes["gateway_id"]
		groupId := rs.Primary.Attributes["group_id"]
		variableName := rs.Primary.Attributes["name"]
		if gatewayId == "" || groupId == "" || variableName == "" {
			return "", fmt.Errorf("missing some attributes, want '<gateway_id>/<group_id>/<name>', but '%s/%s/%s'",
				gatewayId, groupId, variableName)
		}
		return fmt.Sprintf("%s/%s/%s", gatewayId, groupId, variableName), nil
	}
}

func testAccEnvironmentVariable_base(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_apigw_environment_v2" "env"{
  name        = "%[2]s"
  instance_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  description = "test description"
}

resource "opentelekomcloud_apigw_group_v2" "group"{
  instance_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  name        = "%[2]s"
  description = "test description"
}

`, testAccAPIGWv2GatewayBasic(name), name)
}

func testAccEnvironmentVariable_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_apigw_environment_variable_v2" "var" {
  gateway_id     = opentelekomcloud_apigw_gateway_v2.gateway.id
  group_id       = opentelekomcloud_apigw_group_v2.group.id
  environment_id = opentelekomcloud_apigw_environment_v2.env.id
  name           = "%[2]s"
  value          = "/stage/demo"
}
`, testAccEnvironmentVariable_base(name), name)
}
