package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/response"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceApigwResponse = "opentelekomcloud_apigw_response_v2.resp"

func getResponseFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.APIGWV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud APIG v2 client: %s", err)
	}
	return response.Get(client, state.Primary.Attributes["gateway_id"], state.Primary.Attributes["group_id"],
		state.Primary.ID)
}

func TestAccResponse_basic(t *testing.T) {
	var resp response.Response
	name := fmt.Sprintf("apigw_acc_resp%s", acctest.RandString(5))
	updateName := fmt.Sprintf("apigw_acc_resp_update%s", acctest.RandString(5))

	rc := common.InitResourceCheck(
		resourceApigwResponse,
		&resp,
		getResponseFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResponse_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceApigwResponse, "name", name),
				),
			},
			{
				Config: testAccResponse_basic(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceApigwResponse, "name", updateName),
				),
			},
			{
				ResourceName:      resourceApigwResponse,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResponseImportStateFunc(),
			},
		},
	})
}

func TestAccResponse_customRules(t *testing.T) {
	var resp response.Response
	name := fmt.Sprintf("apigw_acc_resp%s", acctest.RandString(5))

	rc := common.InitResourceCheck(
		resourceApigwResponse,
		&resp,
		getResponseFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResponse_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceApigwResponse, "name", name),
				),
			},
			{
				Config: testAccResponse_rules(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceApigwResponse, "rule.#", "2"),
				),
			},
			{
				Config: testAccResponse_rulesUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceApigwResponse, "rule.#", "1"),
				),
			},
			{
				ResourceName:      resourceApigwResponse,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResponseImportStateFunc(),
			},
		},
	})
}

func testAccResponseImportStateFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceApigwResponse]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceApigwResponse, rs)
		}
		if rs.Primary.Attributes["gateway_id"] == "" || rs.Primary.Attributes["group_id"] == "" ||
			rs.Primary.Attributes["name"] == "" {
			return "", fmt.Errorf("missing some attributes, want '{gateway_id}/{group_id}/{name}', but '%s/%s/%s'",
				rs.Primary.Attributes["gateway_id"], rs.Primary.Attributes["group_id"], rs.Primary.Attributes["name"])
		}
		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["gateway_id"], rs.Primary.Attributes["group_id"],
			rs.Primary.Attributes["name"]), nil
	}
}

func testAccResponse_basic(name string) string {
	relatedConfig := testAccApigwApi_base(name)
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_response_v2" "resp" {
  name       = "%[2]s"
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  group_id   = opentelekomcloud_apigw_group_v2.group.id

  rule {
    error_type  = "AUTHORIZER_FAILURE"
    body        = "{\"code\":\"$context.authorizer.frontend.code\",\"message\":\"$context.authorizer.frontend.message\"}"
    status_code = 401
  }
}
`, relatedConfig, name)
}

func testAccResponse_rules(name string) string {
	relatedConfig := testAccApigwApi_base(name)
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_response_v2" "resp" {
  name       = "%[2]s"
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  group_id   = opentelekomcloud_apigw_group_v2.group.id

  rule {
    error_type  = "ACCESS_DENIED"
    body        = "{\"error_code\":\"$context.error.code\",\"error_msg\":\"$context.error.message\"}"
    status_code = 400
  }
  rule {
    error_type  = "AUTHORIZER_FAILURE"
    body        = "{\"code\":\"$context.authorizer.frontend.code\",\"message\":\"$context.authorizer.frontend.message\"}"
    status_code = 401
  }
}
`, relatedConfig, name)
}

func testAccResponse_rulesUpdate(name string) string {
	relatedConfig := testAccApigwApi_base(name)
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_response_v2" "resp" {
  name       = "%[2]s"
  gateway_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  group_id   = opentelekomcloud_apigw_group_v2.group.id

  rule {
    error_type  = "AUTHORIZER_FAILURE"
    body        = "{\"code\":\"$context.authorizer.frontend.code\",\"message\":\"$context.authorizer.frontend.message\"}"
    status_code = 403
  }
}
`, relatedConfig, name)
}
