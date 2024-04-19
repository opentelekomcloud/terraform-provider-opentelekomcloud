package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"

	apis "github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/api"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceApigwApiName = "opentelekomcloud_apigw_api_v2.api"

func getApiFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.APIGWV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud APIG v2 client: %s", err)
	}
	return apis.Get(c, state.Primary.Attributes["gateway_id"], state.Primary.ID)
}

func TestAccApi_basic(t *testing.T) {
	var (
		api         apis.ApiResp
		rName       = resourceApigwApiName
		name        = fmt.Sprintf("apigw_acc_api%s", acctest.RandString(5))
		updateName  = fmt.Sprintf("apigw_acc_api_update%s", acctest.RandString(5))
		basicConfig = testAccApigwApi_base(name)
	)

	rc := common.InitResourceCheck(
		rName,
		&api,
		getApiFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApigwApi_basic(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "Public"),
					resource.TestCheckResourceAttr(rName, "description", "Created by script"),
					resource.TestCheckResourceAttr(rName, "request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(rName, "request_method", "GET"),
					resource.TestCheckResourceAttr(rName, "request_uri", "/user_info/{user_age}"),
					resource.TestCheckResourceAttr(rName, "security_authentication_type", "APP"),
					resource.TestCheckResourceAttr(rName, "match_mode", "EXACT"),
					resource.TestCheckResourceAttr(rName, "success_response", "Success response"),
					resource.TestCheckResourceAttr(rName, "failure_response", "Failed response"),
					resource.TestCheckResourceAttr(rName, "request_params.#", "2"),
					resource.TestCheckResourceAttr(rName, "backend_params.#", "1"),
					resource.TestCheckResourceAttr(rName, "http.0.request_uri", "/getUserAge/{userAge}"),
					resource.TestCheckResourceAttr(rName, "http.0.request_method", "GET"),
					resource.TestCheckResourceAttr(rName, "http.0.request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(rName, "http.0.timeout", "30000"),
					resource.TestCheckResourceAttr(rName, "http_policy.#", "1"),
					resource.TestCheckResourceAttr(rName, "http_policy.0.conditions.#", "1"),
					resource.TestCheckResourceAttr(rName, "mock.#", "0"),
					resource.TestCheckResourceAttr(rName, "func_graph.#", "0"),
					resource.TestCheckResourceAttr(rName, "mock_policy.#", "0"),
					resource.TestCheckResourceAttr(rName, "func_graph_policy.#", "0"),
					resource.TestCheckResourceAttr(rName, "http_policy.0.backend_params.#", "2"),
				),
			},
			{
				Config: testAccApigwApi_update(basicConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "type", "Public"),
					resource.TestCheckResourceAttr(rName, "description", "Updated by script"),
					resource.TestCheckResourceAttr(rName, "request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(rName, "request_method", "GET"),
					resource.TestCheckResourceAttr(rName, "request_uri", "/user_info/{user_name}"),
					resource.TestCheckResourceAttr(rName, "security_authentication_type", "APP"),
					resource.TestCheckResourceAttr(rName, "match_mode", "EXACT"),
					resource.TestCheckResourceAttr(rName, "success_response", "Updated Success response"),
					resource.TestCheckResourceAttr(rName, "failure_response", "Updated Failed response"),
					resource.TestCheckResourceAttr(rName, "request_params.#", "2"),
					resource.TestCheckResourceAttr(rName, "backend_params.#", "2"),
					resource.TestCheckResourceAttr(rName, "http.0.request_uri", "/getUserName/{userName}"),
					resource.TestCheckResourceAttr(rName, "http.0.request_method", "GET"),
					resource.TestCheckResourceAttr(rName, "http.0.request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(rName, "http.0.timeout", "60000"),
					resource.TestCheckResourceAttr(rName, "http_policy.#", "1"),
					resource.TestCheckResourceAttr(rName, "http_policy.0.conditions.#", "2"),
					resource.TestCheckResourceAttr(rName, "mock.#", "0"),
					resource.TestCheckResourceAttr(rName, "func_graph.#", "0"),
					resource.TestCheckResourceAttr(rName, "mock_policy.#", "0"),
					resource.TestCheckResourceAttr(rName, "func_graph_policy.#", "0"),
					resource.TestCheckResourceAttr(rName, "http_policy.0.backend_params.#", "2"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApiResourceImportStateFunc(),
			},
		},
	})
}

func testAccApiResourceImportStateFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rName := "opentelekomcloud_apigw_api_v2.api"
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["gateway_id"] == "" || rs.Primary.Attributes["name"] == "" {
			return "", fmt.Errorf("missing some attributes, want '{gateway_id}/{name}', but '%s/%s'",
				rs.Primary.Attributes["gateway_id"], rs.Primary.Attributes["name"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["gateway_id"], rs.Primary.Attributes["name"]), nil
	}
}

func testAccApigwApi_base(name string) string {
	return fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_apigw_gateway_v2" "gateway" {
  name                            = "%s"
  spec_id                         = "BASIC"
  vpc_id                          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id                       = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  security_group_id               = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zones              = ["eu-de-01", "eu-de-02"]
  description                     = "test gateway 2"
  ingress_bandwidth_size          = 5
  ingress_bandwidth_charging_mode = "bandwidth"
  maintain_begin                  = "02:00:00"
}

resource "opentelekomcloud_apigw_environment_v2" "env" {
  name        = "%[3]s"
  instance_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  description = "test description"
}

resource "opentelekomcloud_apigw_group_v2" "group" {
  name        = "%[3]s"
  instance_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  description = "test description"

  environment {
    variable {
      name  = "test-name"
      value = "test-value"
    }
    environment_id = opentelekomcloud_apigw_environment_v2.env.id
  }
}
`, common.DataSourceSubnet, common.DataSourceSecGroupDefault, name)
}

func testAccApigwApi_basic(relatedConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_api_v2" "api" {
  gateway_id                   = opentelekomcloud_apigw_gateway_v2.gateway.id
  group_id                     = opentelekomcloud_apigw_group_v2.group.id
  name                         = "%[2]s"
  type                         = "Public"
  request_protocol             = "HTTP"
  request_method               = "GET"
  request_uri                  = "/user_info/{user_age}"
  security_authentication_type = "APP"
  match_mode                   = "EXACT"
  success_response             = "Success response"
  failure_response             = "Failed response"
  description                  = "Created by script"

  request_params {
    name     = "user_age"
    type     = "NUMBER"
    location = "PATH"
    required = true
    maximum  = 200
    minimum  = 0
  }
  request_params {
    name        = "X-TEST-ENUM"
    type        = "STRING"
    location    = "HEADER"
    maximum     = 20
    minimum     = 10
    sample      = "ACC_TEST_XXX"
    passthrough = true
    enumeration = "ACC_TEST_A,ACC_TEST_B"
  }

  backend_params {
    type     = "REQUEST"
    name     = "userAge"
    location = "PATH"
    value    = "user_age"
  }

  http {
    url_domain       = "opentelekomcloud.my.com"
    request_uri      = "/getUserAge/{userAge}"
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 30000
    retry_count      = 1
  }

  http_policy {
    url_domain       = "opentelekomcloud.my.com"
    name             = "%[2]s_policy1"
    request_protocol = "HTTP"
    request_method   = "GET"
    effective_mode   = "ANY"
    request_uri      = "/getUserAge/{userAge}"
    timeout          = 30000
    retry_count      = 1

    backend_params {
      type     = "REQUEST"
      name     = "userAge"
      location = "PATH"
      value    = "user_age"
    }
    backend_params {
      type              = "SYSTEM"
      name              = "%[2]s"
      location          = "HEADER"
      value             = "serverName"
      system_param_type = "internal"
    }

    conditions {
      origin     = "param"
      param_name = "user_age"
      type       = "EXACT"
      value      = "28"
    }
  }
}
`, relatedConfig, name)
}

func testAccApigwApi_update(relatedConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_api_v2" "api" {
  gateway_id                   = opentelekomcloud_apigw_gateway_v2.gateway.id
  group_id                     = opentelekomcloud_apigw_group_v2.group.id
  name                         = "%[2]s"
  type                         = "Public"
  request_protocol             = "HTTP"
  request_method               = "GET"
  request_uri                  = "/user_info/{user_name}"
  security_authentication_type = "APP"
  match_mode                   = "EXACT"
  success_response             = "Updated Success response"
  failure_response             = "Updated Failed response"
  description                  = "Updated by script"

  request_params {
    name     = "user_name"
    type     = "STRING"
    location = "PATH"
    required = true
    maximum  = 64
    minimum  = 3
  }
  request_params {
    name        = "X-TEST-ENUM"
    type        = "STRING"
    location    = "HEADER"
    maximum     = 20
    minimum     = 10
    sample      = "ACC_TEST_XXXX"
    passthrough = false
    enumeration = "ACC_TEST_A,ACC_TEST_B,ACC_TEST_C"
  }

  backend_params {
    type     = "REQUEST"
    name     = "userName"
    location = "PATH"
    value    = "user_name"
  }
  backend_params {
    type              = "SYSTEM"
    name              = "%[2]s"
    location          = "HEADER"
    value             = "serverName"
    system_param_type = "internal"
  }

  http {
    url_domain       = "opentelekomcloud.my.com"
    request_uri      = "/getUserName/{userName}"
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 60000
  }

  http_policy {
    url_domain       = "opentelekomcloud.my.com"
    name             = "%[2]s_policy1"
    request_protocol = "HTTP"
    request_method   = "GET"
    effective_mode   = "ANY"
    request_uri      = "/getAdminName/{adminName}"
    timeout          = 60000

    backend_params {
      type     = "REQUEST"
      name     = "adminName"
      location = "PATH"
      value    = "user_name"
    }
    backend_params {
      type              = "SYSTEM"
      name              = "%[2]s"
      location          = "HEADER"
      value             = "serverName"
      system_param_type = "internal"
    }

    conditions {
      origin     = "param"
      param_name = "user_name"
      type       = "EXACT"
      value      = "Administrator"
    }
    conditions {
      origin     = "param"
      param_name = "user_name"
      type       = "EXACT"
      value      = "value_test"
    }
  }
}
`, relatedConfig, name)
}
