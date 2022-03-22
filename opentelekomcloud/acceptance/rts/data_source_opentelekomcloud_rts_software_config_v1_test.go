package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const configDataName = "data.opentelekomcloud_rts_software_config_v1.configs"

func TestAccRTSSoftwareConfigV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRtsSoftwareConfigV1DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRtsSoftwareConfigV1DataSourceID(configDataName),
					resource.TestCheckResourceAttr(configDataName, "name", "opentelekomcloud-config"),
					resource.TestCheckResourceAttr(configDataName, "group", "script"),
				),
			},
		},
	})
}

func testAccCheckRtsSoftwareConfigV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find software config data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("RTS software config data source ID not set ")
		}

		return nil
	}
}

const testAccRtsSoftwareConfigV1DataSourceBasic = `
resource "opentelekomcloud_rts_software_config_v1" "config_1" {
  name = "opentelekomcloud-config"
  output_values = [{
    type         = "String"
    name         = "result"
    error_output = "false"
    description  = "value1"
  }]
  input_values = [{
    default     = "0"
    type        = "String"
    name        = "foo"
    description = "value2"
  }]
  group = "script"
}

data "opentelekomcloud_rts_software_config_v1" "configs" {
  id = opentelekomcloud_rts_software_config_v1.config_1.id
}
`
