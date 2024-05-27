package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/group"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceGroupName = "opentelekomcloud_apigw_group_v2.group"

func getGroupFunc(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.APIGWV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return group.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccAPIGWv2Group_basic(t *testing.T) {
	var groupConfig group.GroupResp
	name := acctest.RandString(10)

	rc := common.InitResourceCheck(
		resourceGroupName,
		&groupConfig,
		getGroupFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGWv2GroupBasic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceGroupName, "name", "group_"+name),
					resource.TestCheckResourceAttr(resourceGroupName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceGroupName, "environment.0.variable.0.name", "test-name"),
					resource.TestCheckResourceAttr(resourceGroupName, "environment.0.variable.0.value", "test-value"),
				),
			},
			{
				Config: testAccAPIGWv2GroupUpdated(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceGroupName, "name", "group_"+name+"_updated"),
					resource.TestCheckResourceAttr(resourceGroupName, "description", "test description updated"),
					resource.TestCheckResourceAttr(resourceGroupName, "environment.0.variable.0.name", "test-name-2"),
					resource.TestCheckResourceAttr(resourceGroupName, "environment.0.variable.0.value", "test-value-2"),
					resource.TestCheckResourceAttr(resourceGroupName, "environment.0.variable.1.name", "test-name"),
					resource.TestCheckResourceAttr(resourceGroupName, "environment.0.variable.1.value", "test-value"),
				),
			},
			{
				ResourceName:      resourceGroupName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAPIGWv2GroupImportStateIdFunc(),
			},
		},
	})
}

func testAccAPIGWv2GroupImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var gatewayID string
		var groupID string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "opentelekomcloud_apigw_gateway_v2" {
				gatewayID = rs.Primary.ID
			} else if rs.Type == "opentelekomcloud_apigw_group_v2" && rs.Primary.ID != "" {
				groupID = rs.Primary.ID
			}
		}
		if gatewayID == "" || groupID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", gatewayID, groupID)
		}
		return fmt.Sprintf("%s/%s", gatewayID, groupID), nil
	}
}

func testAccAPIGWv2GroupBasic(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_apigw_group_v2" "group"{
  instance_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  name        = "%s"
  description = "test description"

  environment {
    variable {
      name  = "test-name"
      value = "test-value"
    }
    environment_id = opentelekomcloud_apigw_environment_v2.env.id
  }
}
`, testAccAPIGWv2EnvironmentBasic("gateway-"+name, "env_"+name), "group_"+name)
}

func testAccAPIGWv2GroupUpdated(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_apigw_group_v2" "group"{
  instance_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  name        = "%s"
  description = "test description updated"

  environment {
    variable {
      name  = "test-name"
      value = "test-value"
    }
    variable {
      name  = "test-name-2"
      value = "test-value-2"
    }
    environment_id = opentelekomcloud_apigw_environment_v2.env.id
  }
}
`, testAccAPIGWv2EnvironmentBasic("gateway-"+name, "env_"+name), "group_"+name+"_updated")
}
