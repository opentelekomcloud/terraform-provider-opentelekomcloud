package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rts/v1/softwareconfig"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const configResourceName = "opentelekomcloud_rts_software_config_v1.config_1"

func TestAccRTSSoftwareConfigV1_basic(t *testing.T) {
	var config softwareconfig.SoftwareConfig

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRtsSoftwareConfigV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRtsSoftwareConfigV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRtsSoftwareConfigV1Exists(configResourceName, &config),
					resource.TestCheckResourceAttr(configResourceName, "name", "opentelekomcloud-config"),
					resource.TestCheckResourceAttr(configResourceName, "group", "script"),
				),
			},
		},
	})
}

func TestAccRTSSoftwareConfigV1_timeout(t *testing.T) {
	var config softwareconfig.SoftwareConfig

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRtsSoftwareConfigV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRtsSoftwareConfigV1Timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRtsSoftwareConfigV1Exists(configResourceName, &config),
				),
			},
		},
	})
}

func testAccCheckRtsSoftwareConfigV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	orchestrationClient, err := config.OrchestrationV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud orchestration client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_rts_software_config_v1" {
			continue
		}

		_, err := softwareconfig.Get(orchestrationClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("RTS Software Config still exists")
		}
	}

	return nil
}

func testAccCheckRtsSoftwareConfigV1Exists(n string, configs *softwareconfig.SoftwareConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		orchestrationClient, err := config.OrchestrationV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud orchestration client: %s", err)
		}

		found, err := softwareconfig.Get(orchestrationClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("RTS Software Config not found")
		}

		*configs = *found

		return nil
	}
}

const (
	testAccRtsSoftwareConfigV1Basic = `
resource "opentelekomcloud_rts_software_config_v1" "config_1" {
  name = "opentelekomcloud-config"
  output_values = [{
    type = "String"
    name = "result"
    error_output = "false"
    description = "value1"
  }]
  input_values=[{
    default = "0"
    type = "String"
    name = "foo"
    description = "value2"
  }]
  group = "script"
}
`
	testAccRtsSoftwareConfigV1Timeout = `
resource "opentelekomcloud_rts_software_config_v1" "config_1" {
  name = "opentelekomcloud-config"
  output_values = [{
    type = "String"
    name = "result"
    error_output = "false"
    description = "value1"
  }]
  input_values=[{
    default = "0"
    type = "String"
    name = "foo"
    description = "value2"
  }]
  group = "script"
  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
)
