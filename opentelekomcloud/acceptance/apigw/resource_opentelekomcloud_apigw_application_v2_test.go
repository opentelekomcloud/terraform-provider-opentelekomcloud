package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/app"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

const resourceApigwAppName = "opentelekomcloud_apigw_application_v2.app"

func getAppFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.APIGWV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud APIG v2 client: %s", err)
	}
	return app.Get(client, state.Primary.Attributes["gateway_id"], state.Primary.ID)
}

func TestAccApplication_basic(t *testing.T) {
	var appResp app.AppResp

	name := fmt.Sprintf("apigw_acc_app%s", acctest.RandString(5))
	updateName := fmt.Sprintf("apigw_acc_app_update%s", acctest.RandString(5))
	description := "Created by script"
	updateDescription := "Updated by script"

	rc := common.InitResourceCheck(
		resourceApigwAppName,
		&appResp,
		getAppFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApplication_basic(name, description),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceApigwAppName, "name", name),
					resource.TestCheckResourceAttr(resourceApigwAppName, "description", description),
					resource.TestCheckResourceAttrSet(resourceApigwAppName, "app_key"),
					resource.TestCheckResourceAttrSet(resourceApigwAppName, "app_secret"),
				),
			},
			{
				Config: testAccApplication_basic(updateName, updateDescription),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceApigwAppName, "name", updateName),
					resource.TestCheckResourceAttr(resourceApigwAppName, "description", updateDescription),
					resource.TestCheckResourceAttrSet(resourceApigwAppName, "app_key"),
					resource.TestCheckResourceAttrSet(resourceApigwAppName, "app_secret"),
				),
			},
			{
				ResourceName:      resourceApigwAppName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApplicationImportIdFunc(),
			},
		},
	})
}

func testAccApplicationImportIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceApigwAppName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceApigwAppName, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["gateway_id"] == "" {
			return "", fmt.Errorf("resource not found: %s/%s", rs.Primary.Attributes["gateway_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["gateway_id"], rs.Primary.ID), nil
	}
}

func testAccApplication_basic(name, description string) string {
	code := hashcode.TryBase64EncodeString(acctest.RandString(64))
	relatedConfig := testAccApigwApi_base(name)
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_application_v2" "app" {
  name        = "%[2]s"
  gateway_id  = opentelekomcloud_apigw_gateway_v2.gateway.id
  description = "%[3]s"

  app_codes = ["%[4]s"]
}
`, relatedConfig, name, description, code)
}
