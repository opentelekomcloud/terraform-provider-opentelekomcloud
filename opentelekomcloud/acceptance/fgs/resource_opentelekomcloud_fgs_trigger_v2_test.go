package fgs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/fgs/v2/function"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/fgs/v2/trigger"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/fgs"
)

func getFunctionTriggerFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.FuncGraphV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating FunctionGraph client: %s", err)
	}

	return fgs.GetTriggerById(client, state.Primary.Attributes["function_urn"], state.Primary.Attributes["type"],
		state.Primary.ID)
}

func TestAccFunctionTrigger_basic(t *testing.T) {
	var (
		relatedFunc      function.FuncGraph
		timeTrigger      trigger.TriggerFuncResp
		randName         = fmt.Sprintf("fgs-acc-api%s", acctest.RandString(5))
		resNameFunc      = "opentelekomcloud_fgs_function_v2.test"
		resNameTimerRate = "opentelekomcloud_fgs_trigger_v2.timer_rate"
		resNameTimerCron = "opentelekomcloud_fgs_trigger_v2.timer_cron"

		rcFunc      = common.InitResourceCheck(resNameFunc, &relatedFunc, getResourceObj)
		rcTimerRate = common.InitResourceCheck(resNameTimerRate, &timeTrigger, getFunctionTriggerFunc)
		rcTimerCron = common.InitResourceCheck(resNameTimerCron, &timeTrigger, getFunctionTriggerFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rcFunc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionTimingTrigger_basic_step1(randName),
				Check: resource.ComposeTestCheckFunc(
					rcTimerRate.CheckResourceExists(),
					resource.TestCheckResourceAttr(resNameTimerRate, "type", "TIMER"),
					resource.TestCheckResourceAttr(resNameTimerRate, "status", "ACTIVE"),
					rcTimerCron.CheckResourceExists(),
					resource.TestCheckResourceAttr(resNameTimerCron, "type", "TIMER"),
					resource.TestCheckResourceAttr(resNameTimerCron, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccFunctionTimingTrigger_basic_step2(randName),
				Check: resource.ComposeTestCheckFunc(
					rcTimerRate.CheckResourceExists(),
					resource.TestCheckResourceAttr(resNameTimerRate, "status", "DISABLED"),
					rcTimerCron.CheckResourceExists(),
					resource.TestCheckResourceAttr(resNameTimerCron, "status", "DISABLED"),
				),
			},
			{
				ResourceName:      resNameTimerRate,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccFunctionTriggerImportStateFunc(resNameTimerRate),
			},
			{
				ResourceName:      resNameTimerCron,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccFunctionTriggerImportStateFunc(resNameTimerCron),
			},
		},
	})
}

func testAccFunctionTriggerImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var functionUrn, triggerType, triggerId string
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of function trigger is not found in the tfstate", rsName)
		}
		functionUrn = rs.Primary.Attributes["function_urn"]
		triggerType = rs.Primary.Attributes["type"]
		triggerId = rs.Primary.ID
		if functionUrn == "" || triggerType == "" || triggerId == "" {
			return "", fmt.Errorf("the function trigger is not exist or related function URN is missing")
		}
		return fmt.Sprintf("%s/%s/%s", functionUrn, triggerType, triggerId), nil
	}
}

func testAccFunctionTrigger_base(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_fgs_function_v2" "test" {
  name        = "%[1]s"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 10
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgW1wcyhldybiBvdXRwdXQ="
}`, name)
}

func testAccFunctionTimingTrigger_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_fgs_trigger_v2" "timer_rate" {
  function_urn = opentelekomcloud_fgs_function_v2.test.urn
  type         = "TIMER"
  event_data   = jsonencode({
    "name": "%[2]s_rate",
    "schedule_type": "Rate",
    "user_event": "Created by acc test",
    "schedule": "3m"
  })
}

resource "opentelekomcloud_fgs_trigger_v2" "timer_cron" {
  function_urn = opentelekomcloud_fgs_function_v2.test.urn
  type         = "TIMER"
  event_data   = jsonencode({
    "name": "%[2]s_cron",
    "schedule_type": "Cron",
    "user_event": "Created by acc test",
    "schedule": "@every 1h30m"
  })
}
`, testAccFunctionTrigger_base(name), name)
}

func testAccFunctionTimingTrigger_basic_step2(name string) string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_fgs_trigger_v2" "timer_rate" {
  function_urn = opentelekomcloud_fgs_function_v2.test.urn
  type         = "TIMER"
  status       = "DISABLED"
  event_data   = jsonencode({
    "name": "%[2]s_rate",
    "schedule_type": "Rate",
    "user_event": "Created by acc test",
    "schedule": "3m"
  })
}

resource "opentelekomcloud_fgs_trigger_v2" "timer_cron" {
  function_urn = opentelekomcloud_fgs_function_v2.test.urn
  type         = "TIMER"
  status       = "DISABLED"
  event_data   = jsonencode({
    "name": "%[2]s_cron",
    "schedule_type": "Cron",
    "user_event": "Created by acc test",
    "schedule": "@every 1h30m"
  })
}
`, testAccFunctionTrigger_base(name), name)
}
