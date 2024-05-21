package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	appauth "github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/app_auth"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceApigwAppAuthName = "opentelekomcloud_apigw_application_authorization_v2.auth"

func getAppAuthFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.APIGWV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud APIG v2 client: %s", err)
	}

	opts := appauth.ListBoundOpts{
		GatewayID: state.Primary.Attributes["gateway_id"],
		AppID:     state.Primary.Attributes["application_id"],
	}
	resp, err := appauth.ListAPIBound(client, opts)
	if err != nil {
		return nil, err
	}
	if len(resp) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	return resp, nil
}

func TestAccAppAuth_basic(t *testing.T) {
	var authApis []appauth.ApiAuth

	rc := common.InitResourceCheck(resourceApigwAppAuthName, &authApis, getAppAuthFunc)
	name := fmt.Sprintf("apigw_acc_app_auth%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAppAuth_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config: testAccAppAuth_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				ResourceName:      resourceApigwAppAuthName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAppAuthImportIdFunc(resourceApigwAppAuthName),
			},
		},
	})
}

func testAccAppAuthImportIdFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rsName, rs)
		}

		gatewayId := rs.Primary.Attributes["gateway_id"]
		resourceId := rs.Primary.ID
		if gatewayId == "" || resourceId == "" {
			return "", fmt.Errorf("missing some attributes, want '<gateway_id>/<id>' (the format of resource ID is "+
				"'<env_id>/<application_id>'), but got '%s/%s'", gatewayId, resourceId)
		}
		return fmt.Sprintf("%s/%s", gatewayId, resourceId), nil
	}
}

func testAccAppAuth_basic(name string) string {
	relatedConfig := testAccApigwApi_base(name)
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_api_v2" "api" {
  count = 3

  gateway_id                   = opentelekomcloud_apigw_gateway_v2.gateway.id
  group_id                     = opentelekomcloud_apigw_group_v2.group.id
  name                         = "%[2]s_${count.index}"
  type                         = "Public"
  request_protocol             = "HTTP"
  request_method               = "GET"
  request_uri                  = "/user_info/{user_age}"
  security_authentication_type = "APP"
  match_mode                   = "EXACT"
  success_response             = "Success response"
  failure_response             = "Failed response"
  description                  = "Created by script"

  http {
    url_domain       = "opentelekomcloud.my.com"
    request_uri      = "/getUserAge/${count.index}"
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 30000
    retry_count      = 1
  }
}

resource "opentelekomcloud_apigw_api_publishment_v2" "pub" {
  gateway_id     = opentelekomcloud_apigw_gateway_v2.gateway.id
  environment_id = opentelekomcloud_apigw_environment_v2.env.id
  api_id         = opentelekomcloud_apigw_api_v2.api.id
}

resource "opentelekomcloud_apigw_application_v2" "app" {
  name        = "%[2]s"
  gateway_id  = opentelekomcloud_apigw_gateway_v2.gateway.id
  description = "acctest description"
}

resource "opentelekomcloud_apigw_application_authorization_v2" "auth" {
  depends_on = [opentelekomcloud_apigw_api_publishment_v2.pub]

  gateway_id    = opentelekomcloud_apigw_gateway_v2.gateway.id
  application_id = opentelekomcloud_apigw_application_v2.app.id
  env_id         = opentelekomcloud_apigw_environment_v2.env.id
  api_ids        = slice(opentelekomcloud_apigw_api_v2.api[*].id, 0, 2)
}
`, relatedConfig, name)
}

func testAccAppAuth_update(name string) string {
	relatedConfig := testAccApigwApi_base(name)
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_api_v2" "api" {
  count = 3

  gateway_id                   = opentelekomcloud_apigw_gateway_v2.gateway.id
  group_id                     = opentelekomcloud_apigw_group_v2.group.id
  name                         = "%[2]s_${count.index}"
  type                         = "Public"
  request_protocol             = "HTTP"
  request_method               = "GET"
  request_uri                  = "/user_info/{user_age}"
  security_authentication_type = "APP"
  match_mode                   = "EXACT"
  success_response             = "Success response"
  failure_response             = "Failed response"
  description                  = "Created by script"

  http {
    url_domain       = "opentelekomcloud.my.com"
    request_uri      = "/getUserAge/${count.index}"
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 30000
    retry_count      = 1
  }
}

resource "opentelekomcloud_apigw_api_publishment_v2" "pub" {
  gateway_id     = opentelekomcloud_apigw_gateway_v2.gateway.id
  environment_id = opentelekomcloud_apigw_environment_v2.env.id
  api_id         = opentelekomcloud_apigw_api_v2.api.id
}

resource "opentelekomcloud_apigw_application_v2" "app" {
  name        = "%[2]s"
  gateway_id  = opentelekomcloud_apigw_gateway_v2.gateway.id
  description = "acctest description"
}

resource "opentelekomcloud_apigw_application_authorization_v2" "auth" {
  depends_on = [opentelekomcloud_apigw_api_publishment_v2.pub]

  gateway_id    = opentelekomcloud_apigw_gateway_v2.gateway.id
  application_id = opentelekomcloud_apigw_application_v2.app.id
  env_id         = opentelekomcloud_apigw_environment_v2.env.id
  api_ids        = slice(opentelekomcloud_apigw_api_v2.api[*].id, 1, 3)
}
`, relatedConfig, name)
}
