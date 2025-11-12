package gemini

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/gemini/v3/template"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const geminiParameterTemplateResourceName = "opentelekomcloud_gemini_template_v3.test"

func getGeminiParameterTemplateResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.GeminiDBV3Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating GeminiDB client: %s", err)
	}

	resp, err := template.Get(client, state.Primary.ID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving GeminiDB parameter template: %s", err)
	}

	return resp, nil
}

func TestAccGeminiParameterTemplate_basic(t *testing.T) {
	var obj interface{}

	name := "tf_gemini_" + acctest.RandString(3)
	nameUpdate := "tf_gemini_" + acctest.RandString(3)

	rc := common.InitResourceCheck(
		geminiParameterTemplateResourceName,
		&obj,
		getGeminiParameterTemplateResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiParameterTemplateBasic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(geminiParameterTemplateResourceName, "name", name),
					resource.TestCheckResourceAttr(geminiParameterTemplateResourceName, "description", "configuration test"),
					resource.TestCheckResourceAttr(geminiParameterTemplateResourceName, "instance_type", "cassandra"),
					resource.TestCheckResourceAttr(geminiParameterTemplateResourceName, "engine_version", "3.11"),
					resource.TestCheckResourceAttr(geminiParameterTemplateResourceName, "parameters.#", "2"),
				),
			},
			{
				Config: testAccGeminiParameterTemplateUpdate(nameUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(geminiParameterTemplateResourceName, "name", nameUpdate),
					resource.TestCheckResourceAttr(geminiParameterTemplateResourceName, "description", "updated configuration"),
					resource.TestCheckResourceAttr(geminiParameterTemplateResourceName, "parameters.#", "3"),
				),
			},
		},
	})
}

func testAccGeminiParameterTemplateBasic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_gemini_template_v3" "test" {
  name           = "%s"
  description    = "configuration test"
  instance_type  = "cassandra"
  engine_version = "3.11"

  parameters {
    name  = "write_request_timeout_in_ms"
    value = "5000"
  }

  parameters {
    name  = "slow_query_log_timeout_in_ms"
    value = "10000"
  }
}
`, name)
}

func testAccGeminiParameterTemplateUpdate(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_gemini_template_v3" "test" {
  name           = "%s"
  description    = "updated configuration"
  instance_type  = "cassandra"
  engine_version = "3.11"

  parameters {
    name  = "write_request_timeout_in_ms"
    value = "6000"
  }

  parameters {
    name  = "slow_query_log_timeout_in_ms"
    value = "15000"
  }

  parameters {
    name  = "read_request_timeout_in_ms"
    value = "8000"
  }
}
`, name)
}
