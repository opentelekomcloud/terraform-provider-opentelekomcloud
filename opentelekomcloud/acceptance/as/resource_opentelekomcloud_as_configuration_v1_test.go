package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v1/configurations"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccASV1Configuration_basic(t *testing.T) {
	var asConfig configurations.Configuration

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccFlavorPreCheck(t) },
		Providers:    common.TestAccProviders,
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

func TestAccASV1Configuration_publicIP(t *testing.T) {
	var asConfig configurations.Configuration

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccFlavorPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckASV1ConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccASV1Configuration_publicIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1ConfigurationExists("opentelekomcloud_as_configuration_v1.as_config", &asConfig),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_as_configuration_v1.as_config", "scaling_configuration_name", "as_config"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_as_configuration_v1.as_config", "instance_config.0.image", env.OS_IMAGE_ID),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_as_configuration_v1.as_config", "instance_config.0.key_name", env.OS_KEYPAIR_NAME),
				),
			},
		},
	})
}

func TestAccASV1Configuration_invalidDiskSize(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccFlavorPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckASV1ConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccASV1Configuration_invalidDiskSize,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`for system disk size should be.+`),
			},
		},
	})
}

func testAccCheckASV1ConfigurationDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	asClient, err := config.AutoscalingV1Client(env.OS_REGION_NAME)
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

		config := common.TestAccProvider.Meta().(*cfg.Config)
		asClient, err := config.AutoscalingV1Client(env.OS_REGION_NAME)
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
`, env.OS_IMAGE_ID, env.OS_KEYPAIR_NAME)

var testAccASV1Configuration_publicIP = fmt.Sprintf(`
resource "opentelekomcloud_as_configuration_v1" "as_config"{
  scaling_configuration_name = "as_config"
  instance_config {
    image = "%s"
    disk {
      size        = 32768
      volume_type = "uh-l1"
      disk_type   = "SYS"
    }
    disk {
      size        = 32768
      volume_type = "co-p1"
      disk_type   = "DATA"
    }
    disk {
      size        = 32768
      volume_type = "uh-l1"
      disk_type   = "DATA"
    }
    key_name = "%s"
    public_ip {
      eip {
        ip_type = "5_mailbgp"
        bandwidth {
          charging_mode = "traffic"
          share_type    = "PER"
          size          = 125
        }
      }
    }
  }
}
`, env.OS_IMAGE_ID, env.OS_KEYPAIR_NAME)

var testAccASV1Configuration_invalidDiskSize = fmt.Sprintf(`
resource "opentelekomcloud_as_configuration_v1" "as_config"{
  scaling_configuration_name = "as_config"
  instance_config {
    image = "%s"
    disk {
      size        = 10
      volume_type = "uh-l1"
      disk_type   = "SYS"
    }
    disk {
      size        = 5
      volume_type = "co-p1"
      disk_type   = "DATA"
    }
    key_name = "%s"
  }
}
`, env.OS_IMAGE_ID, env.OS_KEYPAIR_NAME)
