package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/rts/v1/softwareconfig"
)

func TestAccOTCRtsSoftwareConfigV1_basic(t *testing.T) {
	var config softwareconfig.SoftwareConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRtsSoftwareConfigV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRtsSoftwareConfigV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRtsSoftwareConfigV1Exists("opentelekomcloud_rts_software_config_v1.config_1", &config),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_software_config_v1.config_1", "name", "opentelekomcloud-config"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_software_config_v1.config_1", "group", "script"),
				),
			},
		},
	})
}

func TestAccOTCRtsSoftwareConfigV1_timeout(t *testing.T) {
	var config softwareconfig.SoftwareConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRtsSoftwareConfigV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRtsSoftwareConfigV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRtsSoftwareConfigV1Exists("opentelekomcloud_rts_software_config_v1.config_1", &config),
				),
			},
		},
	})
}

func testAccCheckRtsSoftwareConfigV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	orchestrationClient, err := config.orchestrationV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud orchestration client: %s", err)
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
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		orchestrationClient, err := config.orchestrationV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud orchestration client: %s", err)
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

const testAccRtsSoftwareConfigV1_basic = `
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

const testAccRtsSoftwareConfigV1_timeout = `
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
