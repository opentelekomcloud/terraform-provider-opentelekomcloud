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

const deploymentResourceName = "opentelekomcloud_rts_software_deployment_v1.deployment_1"

func TestAccRTSSoftwareDeploymentV1_basic(t *testing.T) {
	var deployments softwaredeployment.Deployment

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccImagePreCheck(t)
			common.TestAccSubnetPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRTSSoftwareDeploymentV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRtsSoftwareDeploymentV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRTSSoftwareDeploymentV1Exists(deploymentResourceName, &deployments),
					resource.TestCheckResourceAttr(deploymentResourceName, "status_reason", "Deploy data"),
					resource.TestCheckResourceAttr(deploymentResourceName, "status", "IN_PROGRESS"),
					resource.TestCheckResourceAttr(deploymentResourceName, "action", "CREATE"),
				),
			},
			{
				Config: testAccRtsSoftwareDeploymentV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRTSSoftwareDeploymentV1Exists(deploymentResourceName, &deployments),
					resource.TestCheckResourceAttr(deploymentResourceName, "output_values.%", "1"),
					resource.TestCheckResourceAttr(deploymentResourceName, "output_values.deploy_stdout",
						"Writing to /tmp/baaaaa\nWritten to /tmp/baaaaa\n"),
					resource.TestCheckResourceAttr(deploymentResourceName, "status_reason", "Outputs received"),
					resource.TestCheckResourceAttr(deploymentResourceName, "status", "COMPLETE"),
				),
			},
		},
	})
}

func TestAccRTSSoftwareDeploymentV1_timeout(t *testing.T) {
	var deployments softwaredeployment.Deployment

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccImagePreCheck(t)
			common.TestAccSubnetPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRTSSoftwareDeploymentV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRtsSoftwareDeploymentV1Timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRTSSoftwareDeploymentV1Exists(deploymentResourceName, &deployments),
				),
			},
		},
	})
}

func testAccCheckRTSSoftwareDeploymentV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	orchestrationClient, err := config.OrchestrationV1Client(env.OS_REGION_NAME)
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

func testAccCheckRTSSoftwareDeploymentV1Exists(n string, stack *softwaredeployment.Deployment) resource.TestCheckFunc {
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

var testAccRtsSoftwareDeploymentV1Basic = fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_compute_instance_v2" "vm_1" {
  name = "instance_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor_id = "%s"
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
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
`, common.DataSourceImage, common.DataSourceSubnet, env.OS_FLAVOR_ID)

var testAccRtsSoftwareDeploymentV1Update = fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_compute_instance_v2" "vm_1" {
  name = "instance_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor_id = "%s"
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
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
`, common.DataSourceImage, common.DataSourceSubnet, env.OsFlavorID)

var testAccRtsSoftwareDeploymentV1Timeout = fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_compute_instance_v2" "vm_1" {
  name = "instance_1"
  image_id = data.opentelekomcloud_images_image_v2.latest_image.id
  flavor_id = "%s"
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
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
`, common.DataSourceImage, common.DataSourceSubnet, env.OS_FLAVOR_ID)
