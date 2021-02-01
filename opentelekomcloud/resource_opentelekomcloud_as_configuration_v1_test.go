package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v1/configurations"
)

func TestAccASV1Configuration_basic(t *testing.T) {
	var asConfig configurations.Configuration

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccAsConfigPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckASV1ConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccASV1Configuration_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1ConfigurationExists(
						"opentelekomcloud_as_configuration_v1.hth_as_config", &asConfig),
				),
			},
		},
	})
}

func testAccCheckASV1ConfigurationDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	asClient, err := config.autoscalingV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud AutoScaling client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_as_configuration_v1" {
			continue
		}

		_, err := configurations.Get(asClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("AS configuration still exists")
		}
	}

	return nil
}

func testAccCheckASV1ConfigurationExists(n string, configuration *configurations.Configuration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		asClient, err := config.autoscalingV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud AutoScaling client: %s", err)
		}

		found, err := configurations.Get(asClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("autoscaling Configuration not found")
		}
		configuration = &found

		return nil
	}
}

var testAccASV1Configuration_basic = fmt.Sprintf(`
resource "opentelekomcloud_as_configuration_v1" "hth_as_config"{
  scaling_configuration_name = "hth_as_config"
  instance_config {
    image = "%s"
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "DATA"
    }
    key_name = "%s"
  }
}
`, OS_IMAGE_ID, OS_KEYPAIR_NAME)
