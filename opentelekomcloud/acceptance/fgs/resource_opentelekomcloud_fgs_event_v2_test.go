package fgs

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/fgs/v2/events"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getFunctionEventResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.FuncGraphV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud FunctionGraph client: %s", err)
	}

	requestResp, err := events.Get(c, state.Primary.Attributes["function_urn"], state.Primary.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting FunctionGraph function event: %s", err)
	}
	return requestResp, nil
}

func TestAccFunctionEvent_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName    = "opentelekomcloud_fgs_event_v2.test"
		name            = fmt.Sprintf("fgs-events%s", acctest.RandString(5))
		eventContent    = base64.StdEncoding.EncodeToString([]byte("{\"foo\": \"bar\"}"))
		newEventContent = base64.StdEncoding.EncodeToString([]byte("{\"key\": \"value\"}"))

		rc = common.InitResourceCheck(resourceName, &obj, getFunctionEventResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionEvent_basic(name, eventContent),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "function_urn", "opentelekomcloud_fgs_function_v2.test", "urn"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "content", eventContent),
					resource.TestCheckResourceAttr(resourceName, "updated_at", "0"),
				),
			},
			{
				Config: testAccFunctionEvent_basic(name, newEventContent),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "function_urn", "opentelekomcloud_fgs_function_v2.test", "urn"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "content", newEventContent),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccFunctionEventImportStateFunc(resourceName),
			},
		},
	})
}

func testAccFunctionEventImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var functionUrn, eventId string
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of function event is not found in the tfstate", rsName)
		}
		functionUrn = rs.Primary.Attributes["function_urn"]
		eventId = rs.Primary.Attributes["id"]
		if functionUrn == "" || eventId == "" {
			return "", fmt.Errorf("the value of function URN or event name is empty")
		}
		return fmt.Sprintf("%s/%s", functionUrn, eventId), nil
	}
}

func testAccFunctionEvent_basic(name, funcCode string) string {
	return fmt.Sprintf(`
variable "js_script_content" {
  default = <<EOT
exports.handler = async (event, context) => {
    const result =
    {
        'repsonse_code': 200,
        'headers':
        {
            'Content-Type': 'application/json'
        },
        'isBase64Encoded': false,
        'body': JSON.stringify(event)
    }
    return result
}
EOT
}

resource "opentelekomcloud_fgs_function_v2" "test" {
  name        = "%[1]s"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  code_type   = "inline"
  runtime     = "Node.js12.13"
  func_code   = base64encode(jsonencode(var.js_script_content))
}

resource "opentelekomcloud_fgs_event_v2" "test" {
  function_urn = opentelekomcloud_fgs_function_v2.test.urn
  name         = "%[1]s"
  content      = "%[2]s"
}
`, name, funcCode)
}
