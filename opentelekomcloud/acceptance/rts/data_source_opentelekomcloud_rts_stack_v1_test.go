package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const stackDataName = "data.opentelekomcloud_rts_stack_v1.stacks"

func TestAccRTSStackV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRTSStackV1DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRTSStackV1DataSourceID(stackDataName),
					resource.TestCheckResourceAttr(stackDataName, "name", "terraform_provider_stack"),
					resource.TestCheckResourceAttr(stackDataName, "disable_rollback", "true"),
					resource.TestCheckResourceAttr(stackDataName, "parameters.%", "4"),
					resource.TestCheckResourceAttr(stackDataName, "status", "CREATE_COMPLETE"),
				),
			},
		},
	})
}

func testAccCheckRTSStackV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find rts data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("RTS data source ID not set ")
		}

		return nil
	}
}

var testAccRTSStackV1DataSourceBasic = `
resource "opentelekomcloud_rts_stack_v1" "stack_1" {
  name             = "terraform_provider_stack"
  disable_rollback = true
  timeout_mins     = 60
  template_body    = <<JSON
{
  "outputs": {
    "str1": {
      "description": "The description of the nat server.",
      "value": {
        "get_resource": "random"
      }
    }
  },
  "heat_template_version": "2013-05-23",
  "description": "A HOT template that create a single server and boot from volume.",
  "parameters": {
    "key_name": {
      "type": "string",
      "description": "Name of existing key pair for the instance to be created.",
      "default": "KeyPair-click2cloud"
    }
  },
  "resources": {
    "random": {
      "type": "OS::Heat::RandomString",
      "properties": {
        "length": "6"
      }
    }
  }
}
JSON
}

data "opentelekomcloud_rts_stack_v1" "stacks" {
  name = opentelekomcloud_rts_stack_v1.stack_1.name
}
`
