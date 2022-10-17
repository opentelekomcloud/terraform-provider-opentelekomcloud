package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v1/configurations"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccASV1Configuration_basic(t *testing.T) {
	var asConfig configurations.Configuration
	resourceName := "opentelekomcloud_as_configuration_v1.as_config"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			quotas.BookOne(t, quotas.ASConfiguration)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASV1ConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccASV1ConfigurationBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1ConfigurationExists(resourceName, &asConfig),
				),
			},
		},
	})
}

func TestAccASV1Configuration_publicIP(t *testing.T) {
	var asConfig configurations.Configuration
	resourceName := "opentelekomcloud_as_configuration_v1.as_config"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			quotas.BookOne(t, quotas.ASConfiguration)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASV1ConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccASV1ConfigurationPublicIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1ConfigurationExists(resourceName, &asConfig),
					resource.TestCheckResourceAttr(resourceName, "scaling_configuration_name", "as_config"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_config.0.image"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.key_name", env.OS_KEYPAIR_NAME),
				),
			},
		},
	})
}

func TestAccASV1Configuration_invalidDiskSize(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			quotas.BookOne(t, quotas.ASConfiguration)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASV1ConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccASV1ConfigurationInvalidDiskSize,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`for system disk size should be.+`),
			},
		},
	})
}

func TestAccASV1Configuration_multipleSecurityGroups(t *testing.T) {
	var asConfig configurations.Configuration
	resourceName := "opentelekomcloud_as_configuration_v1.as_config"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
			qts := quotas.MultipleQuotas{
				{Q: quotas.ASConfiguration, Count: 1},
				{Q: quotas.SecurityGroup, Count: 3},
			}
			quotas.BookMany(t, qts)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASV1ConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccASV1ConfigurationMultipleSecurityGroups,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1ConfigurationExists(resourceName, &asConfig),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.security_groups.#", "3"),
				),
			},
		},
	})
}

func testAccCheckASV1ConfigurationDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.AutoscalingV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud AutoScalingV1 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_as_configuration_v1" {
			continue
		}

		_, err := configurations.Get(client, rs.Primary.ID)
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
		client, err := config.AutoscalingV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud AutoScalingV1 client: %w", err)
		}

		found, err := configurations.Get(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("autoscaling Configuration not found")
		}
		configuration = found

		return nil
	}
}

var testAccASV1ConfigurationBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_as_configuration_v1" "as_config"{
  scaling_configuration_name = "hth_as_config"
  instance_config {
    image = data.opentelekomcloud_images_image_v2.latest_image.id
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
    user_data = "#! /bin/bash"
  }
}
`, common.DataSourceImage, env.OS_KEYPAIR_NAME)

var testAccASV1ConfigurationPublicIP = fmt.Sprintf(`
%s

resource "opentelekomcloud_as_configuration_v1" "as_config"{
  scaling_configuration_name = "as_config"
  instance_config {
    image = data.opentelekomcloud_images_image_v2.latest_image.id
    disk {
      size        = 1000
      volume_type = "uh-l1"
      disk_type   = "SYS"
    }
    disk {
      size        = 1000
      volume_type = "co-p1"
      disk_type   = "DATA"
    }
    disk {
      size        = 1000
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
`, common.DataSourceImage, env.OS_KEYPAIR_NAME)

var testAccASV1ConfigurationInvalidDiskSize = fmt.Sprintf(`
%s

resource "opentelekomcloud_as_configuration_v1" "as_config"{
  scaling_configuration_name = "as_config"
  instance_config {
    image = data.opentelekomcloud_images_image_v2.latest_image.id
    disk {
      size        = 1
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
`, common.DataSourceImage, env.OS_KEYPAIR_NAME)

var testAccASV1ConfigurationMultipleSecurityGroups = fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_secgroup_v2" "secgroup_1" {
  name        = "acc-test-sg-1"
  description = "Security group for AS config tf test"
}

resource "opentelekomcloud_compute_secgroup_v2" "secgroup_2" {
  name        = "acc-test-sg-2"
  description = "Security group for AS config tf test"
}

resource "opentelekomcloud_compute_secgroup_v2" "secgroup_3" {
  name        = "acc-test-sg-3"
  description = "Security group for AS config tf test"
}

resource "opentelekomcloud_as_configuration_v1" "as_config" {
  scaling_configuration_name = "as_config"
  instance_config {
    image = data.opentelekomcloud_images_image_v2.latest_image.id
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
    security_groups = [
      opentelekomcloud_compute_secgroup_v2.secgroup_1.id,
      opentelekomcloud_compute_secgroup_v2.secgroup_2.id,
      opentelekomcloud_compute_secgroup_v2.secgroup_3.id
    ]
  }
}
`, common.DataSourceImage, env.OS_KEYPAIR_NAME)
