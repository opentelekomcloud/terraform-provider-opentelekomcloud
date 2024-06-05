package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/authorizer"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceApigwAuthorizerName = "opentelekomcloud_apigw_custom_authorizer_v2.authorizer"

func getCustomAuthorizerFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.APIGWV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud APIG v2 client: %s", err)
	}
	return authorizer.Get(client, state.Primary.Attributes["gateway_id"], state.Primary.ID)
}

func TestAccCustomAuthorizer_basic(t *testing.T) {
	var auth authorizer.AuthorizerResp

	name := fmt.Sprintf("apigw_acc_authorizer%s", acctest.RandString(5))
	updateName := fmt.Sprintf("apigw_acc_authorizer_up%s", acctest.RandString(5))

	rc := common.InitResourceCheck(
		resourceApigwAuthorizerName,
		&auth,
		getCustomAuthorizerFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomAuthorizer_front(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceApigwAuthorizerName, "name", name),
					resource.TestCheckResourceAttr(resourceApigwAuthorizerName, "type", "FRONTEND"),
					resource.TestCheckResourceAttr(resourceApigwAuthorizerName, "is_body_send", "true"),
					resource.TestCheckResourceAttr(resourceApigwAuthorizerName, "ttl", "60"),
					resource.TestCheckResourceAttr(resourceApigwAuthorizerName, "identity.#", "1"),
				),
			},
			{
				Config: testAccCustomAuthorizer_frontUpdate(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceApigwAuthorizerName, "name", updateName),
					resource.TestCheckResourceAttr(resourceApigwAuthorizerName, "type", "FRONTEND"),
					resource.TestCheckResourceAttr(resourceApigwAuthorizerName, "is_body_send", "false"),
					resource.TestCheckResourceAttr(resourceApigwAuthorizerName, "ttl", "0"),
					resource.TestCheckResourceAttr(resourceApigwAuthorizerName, "identity.#", "0"),
				),
			},
			{
				ResourceName:      resourceApigwAuthorizerName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCustomAuthorizerImportStateFunc(),
			},
		},
	})
}

func TestAccCustomAuthorizer_backend(t *testing.T) {
	var auth authorizer.AuthorizerResp

	name := fmt.Sprintf("apigw_acc_authorizer%s", acctest.RandString(5))
	updateName := fmt.Sprintf("apigw_acc_authorizer_up%s", acctest.RandString(5))

	rc := common.InitResourceCheck(
		resourceApigwAuthorizerName,
		&auth,
		getCustomAuthorizerFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomAuthorizer_backend(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceApigwAuthorizerName, "name", name),
					resource.TestCheckResourceAttr(resourceApigwAuthorizerName, "type", "BACKEND"),
					resource.TestCheckResourceAttr(resourceApigwAuthorizerName, "is_body_send", "false"),
					resource.TestCheckResourceAttr(resourceApigwAuthorizerName, "ttl", "60"),
				),
			},
			{
				Config: testAccCustomAuthorizer_backendUpdate(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceApigwAuthorizerName, "name", updateName),
					resource.TestCheckResourceAttr(resourceApigwAuthorizerName, "type", "BACKEND"),
					resource.TestCheckResourceAttr(resourceApigwAuthorizerName, "is_body_send", "false"),
					resource.TestCheckResourceAttr(resourceApigwAuthorizerName, "ttl", "45"),
				),
			},
			{
				ResourceName:      resourceApigwAuthorizerName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCustomAuthorizerImportStateFunc(),
			},
		},
	})
}

func testAccCustomAuthorizerImportStateFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceApigwAuthorizerName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceApigwAuthorizerName, rs)
		}
		if rs.Primary.Attributes["gateway_id"] == "" || rs.Primary.Attributes["name"] == "" {
			return "", fmt.Errorf("missing some attributes, want '{gateway_id}/{name}', but '%s/%s'",
				rs.Primary.Attributes["gateway_id"], rs.Primary.Attributes["name"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["gateway_id"], rs.Primary.Attributes["name"]), nil
	}
}

func testAccCustomAuthorizer_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_fgs_function_v2" "test" {
  name        = "%[2]s"
  app         = "default"
  description = "API custom authorization test"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"

  func_code = <<EOF
# -*- coding:utf-8 -*-
import json
def handler(event, context):
    if event["headers"]["authorization"]=='Basic dXNlcjE6cGFzc3dvcmQ=':
        return {
            'statusCode': 200,
            'body': json.dumps({
                "status":"allow",
                "context":{
                    "user_name":"user1"
                }
            })
        }
    else:
        return {
            'statusCode': 200,
            'body': json.dumps({
                "status":"deny",
                "context":{
                    "code":"1001",
                    "message":"incorrect username or password"
                }
            })
        }
EOF
}
`, testAccAPIGWv2GatewayBasic(name), name)
}

func testAccCustomAuthorizer_front(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_custom_authorizer_v2" "authorizer" {
  gateway_id   = opentelekomcloud_apigw_gateway_v2.gateway.id
  name         = "%[2]s"
  function_urn = opentelekomcloud_fgs_function_v2.test.urn
  type         = "FRONTEND"
  is_body_send = true
  ttl          = 60

  identity {
    name     = "user_name"
    location = "QUERY"
  }
}
`, testAccCustomAuthorizer_base(name), name)
}

func testAccCustomAuthorizer_frontUpdate(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_custom_authorizer_v2" "authorizer" {
  gateway_id   = opentelekomcloud_apigw_gateway_v2.gateway.id
  name         = "%[2]s"
  function_urn = opentelekomcloud_fgs_function_v2.test.urn
  type         = "FRONTEND"
}
`, testAccCustomAuthorizer_base(name), name)
}

func testAccCustomAuthorizer_backend(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_custom_authorizer_v2" "authorizer" {
  gateway_id   = opentelekomcloud_apigw_gateway_v2.gateway.id
  name         = "%[2]s"
  function_urn = opentelekomcloud_fgs_function_v2.test.urn
  type         = "BACKEND"
  ttl          = 60
}
`, testAccCustomAuthorizer_base(name), name)
}

func testAccCustomAuthorizer_backendUpdate(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_custom_authorizer_v2" "authorizer" {
  gateway_id   = opentelekomcloud_apigw_gateway_v2.gateway.id
  name         = "%[2]s"
  function_urn = opentelekomcloud_fgs_function_v2.test.urn
  type         = "BACKEND"
  ttl          = 45
}
`, testAccCustomAuthorizer_base(name), name)
}
