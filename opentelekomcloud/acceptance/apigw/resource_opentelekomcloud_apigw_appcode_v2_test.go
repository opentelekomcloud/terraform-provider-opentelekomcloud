package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	appcode "github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/app_code"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

const resourceApigwAppCodeName = "opentelekomcloud_apigw_appcode_v2.code"

func getAppcodeFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.APIGWV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud APIG v2 client: %s", err)
	}
	return appcode.Get(client, state.Primary.Attributes["gateway_id"],
		state.Primary.Attributes["application_id"], state.Primary.ID)
}

func TestAccAppcode_auto(t *testing.T) {
	var appCode appcode.CodeResp

	rc := common.InitResourceCheck(resourceApigwAppCodeName, &appCode, getAppcodeFunc)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAppcode_autoConfig(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				ResourceName:      resourceApigwAppCodeName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAppcodeImportIdFunc(),
			},
		},
	})
}

func TestAccAppcode_manualConfig(t *testing.T) {
	var appCode appcode.CodeResp

	rc := common.InitResourceCheck(resourceApigwAppCodeName, &appCode, getAppcodeFunc)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAppcode_manualConfig(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				ResourceName:      resourceApigwAppCodeName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAppcodeImportIdFunc(),
			},
		},
	})
}

func testAccAppcodeImportIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceApigwAppCodeName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceApigwAppCodeName, rs)
		}
		gatewayId := rs.Primary.Attributes["gateway_id"]
		appId := rs.Primary.Attributes["application_id"]
		appCodeId := rs.Primary.ID
		if gatewayId == "" || appId == "" || appCodeId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<gateway_id>/<application_id>/<id>', but got '%s/%s/%s'",
				gatewayId, appId, appCodeId)
		}
		return fmt.Sprintf("%s/%s/%s", gatewayId, appId, appCodeId), nil
	}
}

func testAccAppcode_autoConfig() string {
	appName := fmt.Sprintf("apigw_acc_app%s", acctest.RandString(5))
	gwName := fmt.Sprintf("apigw_acc_gw%s", acctest.RandString(5))
	relatedConfig := testAccApigwApi_base(gwName)
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_apigw_application_v2" "app" {
  name        = "%[2]s"
  gateway_id  = opentelekomcloud_apigw_gateway_v2.gateway.id
  description = "acctest description"
}

resource "opentelekomcloud_apigw_appcode_v2" "code" {
  gateway_id     = opentelekomcloud_apigw_gateway_v2.gateway.id
  application_id = opentelekomcloud_apigw_application_v2.app.id
}
`, relatedConfig, appName)
}

func testAccAppcode_manualConfig() string {
	code := hashcode.TryBase64EncodeString(acctest.RandString(64))
	appName := fmt.Sprintf("apigw_acc_app%s", acctest.RandString(5))
	gwName := fmt.Sprintf("apigw_acc_gw%s", acctest.RandString(5))
	relatedConfig := testAccApigwApi_base(gwName)

	return fmt.Sprintf(`
%s

resource "opentelekomcloud_apigw_application_v2" "app" {
  name        = "%[2]s"
  gateway_id  = opentelekomcloud_apigw_gateway_v2.gateway.id
  description = "acctest description"
}

resource "opentelekomcloud_apigw_appcode_v2" "code" {
  gateway_id     = opentelekomcloud_apigw_gateway_v2.gateway.id
  application_id = opentelekomcloud_apigw_application_v2.app.id
  value          = "%[3]s"
}
`, relatedConfig, appName, code)
}
