package fgs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/fgs/v2/async_config"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getAsyncInvokeConfigFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.FuncGraphV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud FunctionGraph client: %s", err)
	}
	return async_config.Get(c, state.Primary.ID)
}

func TestAccAsyncInvokeConfig_basic(t *testing.T) {
	var config async_config.AsyncInvokeResp
	name := fmt.Sprintf("fgs-async-config-%s", acctest.RandString(5))
	rName := "opentelekomcloud_fgs_async_invoke_config_v2.test"

	rc := common.InitResourceCheck(rName, &config, getAsyncInvokeConfigFunc)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			common.TestAccPreCheckFgsAgency(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testaccasyncinvokeconfigBasicStep1(name, common.OS_FGS_AGENCY_NAME),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "function_urn",
						"opentelekomcloud_fgs_function_v2.test", "urn"),
					resource.TestCheckResourceAttr(rName, "max_async_event_age_in_seconds", "3500"),
					resource.TestCheckResourceAttr(rName, "max_async_retry_attempts", "2"),
					resource.TestCheckResourceAttr(rName, "on_success.0.destination", "OBS"),
					resource.TestCheckResourceAttrSet(rName, "on_success.0.param"),
					resource.TestCheckResourceAttr(rName, "on_failure.0.destination", "SMN"),
					resource.TestCheckResourceAttrSet(rName, "on_failure.0.param"),
					// resource.TestCheckResourceAttr(rName, "enable_async_status_log", "true"),
				),
			},
			{
				Config: testaccasyncinvokeconfigBasicStep2(name, common.OS_FGS_AGENCY_NAME),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "function_urn",
						"opentelekomcloud_fgs_function_v2.test", "urn"),
					resource.TestCheckResourceAttr(rName, "max_async_event_age_in_seconds", "4000"),
					resource.TestCheckResourceAttr(rName, "max_async_retry_attempts", "3"),
					resource.TestCheckResourceAttr(rName, "on_success.0.destination", "DIS"),
					resource.TestCheckResourceAttrSet(rName, "on_success.0.param"),
					resource.TestCheckResourceAttr(rName, "on_failure.0.destination", "FunctionGraph"),
					resource.TestCheckResourceAttrSet(rName, "on_failure.0.param"),
					// resource.TestCheckResourceAttr(rName, "enable_async_status_log", "false"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testaccasyncinvokeconfigBasicStep1(name, agency string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "test" {
  bucket        = "%[1]s"
  acl           = "private"
  force_destroy = true
}

resource "opentelekomcloud_smn_topic_v2" "test" {
  name = "%[1]s"
}

resource "opentelekomcloud_fgs_function_v2" "test" {
  name        = "%[1]s"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "e42a37a22f4988ba7a681e3042e5c7d13c04e6c1"
  agency      = "%[2]s"
}

resource "opentelekomcloud_fgs_async_invoke_config_v2" "test" {
  function_urn                   = opentelekomcloud_fgs_function_v2.test.urn
  max_async_event_age_in_seconds = 3500
  max_async_retry_attempts       = 2

  on_success {
    destination = "OBS"
    param = jsonencode({
      bucket  = opentelekomcloud_obs_bucket.test.bucket
      prefix  = "/success"
      expires = 5
    })
  }

  on_failure {
    destination = "SMN"
    param = jsonencode({
      topic_urn = opentelekomcloud_smn_topic_v2.test.topic_urn
    })
  }
}
`, name, agency)
}

func testaccasyncinvokeconfigBasicStep2(name, agency string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_dis_stream_v2" "test" {
  name            = "%[2]s"
  partition_count = 1
}

resource "opentelekomcloud_fgs_function_v2" "failure_transport" {
  name        = "%[2]s-failure-transport"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "e42a37a22f4988ba7a681e3042e5c7d13c04e6c1"
}

resource "opentelekomcloud_fgs_function_v2" "test" {
  name        = "%[2]s"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "e42a37a22f4988ba7a681e3042e5c7d13c04e6c1"
  agency      = "%[3]s"
}

resource "opentelekomcloud_fgs_async_invoke_config_v2" "test" {
  function_urn                   = opentelekomcloud_fgs_function_v2.test.urn
  max_async_event_age_in_seconds = 4000
  max_async_retry_attempts       = 3

  on_success {
    destination = "DIS"
    param = jsonencode({
      stream_name = opentelekomcloud_dis_stream_v2.test.name
    })
  }

  on_failure {
    destination = "FunctionGraph"
    param = jsonencode({
      func_urn = opentelekomcloud_fgs_function_v2.failure_transport.id
    })
  }
}
`, common.DataSourceSubnet, name, agency)
}
