package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rts/v1/softwaredeployment"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccOTCRtsSoftwareDeploymentV1_basic(t *testing.T) {
	var deployments softwaredeployment.Deployment

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckOTCRtsSoftwareDeploymentV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRtsSoftwareDeploymentV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRtsSoftwareDeploymentV1Exists("opentelekomcloud_rts_software_deployment_v1.deployment_1", &deployments),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_software_deployment_v1.deployment_1", "status_reason", "Deploy data"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_software_deployment_v1.deployment_1", "status", "IN_PROGRESS"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_software_deployment_v1.deployment_1", "action", "CREATE"),
				),
			},
			{
				Config: testAccRtsSoftwareDeploymentV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRtsSoftwareDeploymentV1Exists("opentelekomcloud_rts_software_deployment_v1.deployment_1", &deployments),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_software_deployment_v1.deployment_1", "output_values.%", "1"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_software_deployment_v1.deployment_1", "output_values.deploy_stdout", "Writing to /tmp/baaaaa\nWritten to /tmp/baaaaa\n"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_software_deployment_v1.deployment_1", "status_reason", "Outputs received"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_software_deployment_v1.deployment_1", "status", "COMPLETE"),
				),
			},
		},
	})
}

func TestAccOTCRtsSoftwareDeploymentV1_timeout(t *testing.T) {
	var deployments softwaredeployment.Deployment

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckOTCRtsSoftwareDeploymentV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRtsSoftwareDeploymentV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRtsSoftwareDeploymentV1Exists("opentelekomcloud_rts_software_deployment_v1.deployment_1", &deployments),
				),
			},
		},
	})
}

func testAccCheckOTCRtsSoftwareDeploymentV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	orchestrationClient, err := config.OrchestrationV1Client(env.OsRegionName)
	if err != nil {
		return fmt.Errorf("error creating RTS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_rts_software_deployment_v1" {
			continue
		}

		stack, err := softwaredeployment.Get(orchestrationClient, rs.Primary.ID).Extract()

		if err == nil {
			if stack.Status != "DELETE_COMPLETE" {
				return fmt.Errorf("deployment still exists")
			}
		}
	}

	return nil
}

func testAccCheckOTCRtsSoftwareDeploymentV1Exists(n string, stack *softwaredeployment.Deployment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		orchestrationClient, err := config.OrchestrationV1Client(env.OsRegionName)
		if err != nil {
			return fmt.Errorf("error creating RTS Client : %s", err)
		}

		found, err := softwaredeployment.Get(orchestrationClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("deployment not found")
		}

		*stack = *found

		return nil
	}
}

var testAccRtsSoftwareDeploymentV1_basic = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "vm_1" {
  name = "instance_1"
  image_id = "%s"
  flavor_id = "%s"
  network {
    uuid = "%s"
  }
}
resource "opentelekomcloud_rts_software_config_v1" "config_1" {
  name="terraform-provider_test"
}

resource "opentelekomcloud_rts_software_deployment_v1" "deployment_1" {
  config_id = opentelekomcloud_rts_software_config_v1.config_1.id
  server_id = opentelekomcloud_compute_instance_v2.vm_1.id
  status= "IN_PROGRESS"
  action= "CREATE"
  status_reason= "Deploy data"
}
`, env.OsImageID, env.OsFlavorID, env.OsNetworkID)

var testAccRtsSoftwareDeploymentV1_update = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "vm_1" {
  name = "instance_1"
  image_id = "%s"
  flavor_id = "%s"
  network {
    uuid = "%s"
  }
}
resource "opentelekomcloud_rts_software_config_v1" "config_1" {
  name="terraform-provider_test"
}

resource "opentelekomcloud_rts_software_deployment_v1" "deployment_1" {
  config_id = opentelekomcloud_rts_software_config_v1.config_1.id
  server_id = opentelekomcloud_compute_instance_v2.vm_1.id
  output_values = {
    deploy_stdout= "Writing to /tmp/baaaaa\nWritten to /tmp/baaaaa\n"
  }
  status= "COMPLETE"
  action= "CREATE"
  status_reason= "Outputs received"
}
`, env.OsImageID, env.OsFlavorID, env.OsNetworkID)

var testAccRtsSoftwareDeploymentV1_timeout = fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "vm_1" {
  name = "instance_1"
  image_id = "%s"
  flavor_id = "%s"
  network {
    uuid = "%s"
  }
}
resource "opentelekomcloud_rts_software_config_v1" "config_1" {
  name="terraform-provider_test"
}

resource "opentelekomcloud_rts_software_deployment_v1" "deployment_1" {
  config_id = opentelekomcloud_rts_software_config_v1.config_1.id
  server_id = opentelekomcloud_compute_instance_v2.vm_1.id
  status= "COMPLETE"
  action= "CREATE"
  status_reason= "Outputs received"

  timeouts {
    create = "10m"
    delete = "10m"
  }
}
`, env.OsImageID, env.OsFlavorID, env.OsNetworkID)
