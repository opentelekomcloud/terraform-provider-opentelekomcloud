package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const resourceDataName = "data.opentelekomcloud_rts_stack_resource_v1.resource_1"

func TestAccRTSStackResourcesV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRTSStackResourcesV1Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRTSStackResourcesV1DataSourceID(resourceDataName),
					resource.TestCheckResourceAttr(resourceDataName, "resource_name", "random"),
					resource.TestCheckResourceAttr(resourceDataName, "resource_type", "OS::Heat::RandomString"),
					resource.TestCheckResourceAttr(resourceDataName, "resource_status", "CREATE_COMPLETE"),
				),
			},
		},
	})
}

func testAccCheckRTSStackResourcesV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find stack resource data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("stack resource data source ID not set")
		}

		return nil
	}
}

const testAccDataSourceRTSStackResourcesV1Config = `
resource "opentelekomcloud_rts_stack_v1" "stack_1" {
  name             = "opentelekomcloud_rts_stack"
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

data "opentelekomcloud_rts_stack_resource_v1" "resource_1" {
  stack_name    = opentelekomcloud_rts_stack_v1.stack_1.name
  resource_name = "random"
}
`
