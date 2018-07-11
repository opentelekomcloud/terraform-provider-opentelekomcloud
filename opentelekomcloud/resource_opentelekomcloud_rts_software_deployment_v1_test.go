package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/rts/v1/softwaredeployment"
)

func TestAccOTCRTSSoftwareDeploymentV1_basic(t *testing.T) {
	var deployments softwaredeployment.Deployment

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCRTSSoftwareDeploymentV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRTSSoftwareDeploymentV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRTSSoftwareDeploymentV1Exists("opentelekomcloud_rts_software_deployment_v1.deployment_1", &deployments),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_software_deployment_v1.deployment_1", "status_reason", "Deploy data"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_software_deployment_v1.deployment_1", "status", "IN_PROGRESS"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_software_deployment_v1.deployment_1", "action", "COMPLETE"),
				),
			},
			resource.TestStep{
				Config: testAccRTSSoftwareDeploymentV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRTSSoftwareDeploymentV1Exists("opentelekomcloud_rts_software_deployment_v1.deployment_1", &deployments),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_software_deployment_v1.deployment_1", "output_values.#", "4"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_software_deployment_v1.deployment_1", "status_reason", "Outputs received"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_software_deployment_v1.deployment_1", "status", "COMPLETE"),
				),
			},
		},
	})
}

// PASS
func TestAccOTCRTSSoftwareDeploymentV1_timeout(t *testing.T) {
	var deployments softwaredeployment.Deployment

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCRTSSoftwareDeploymentV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRTSSoftwareDeploymentV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRTSSoftwareDeploymentV1Exists("opentelekomcloud_rts_software_deployment_v1.deployment_1", &deployments),
				),
			},
		},
	})
}

func testAccCheckOTCRTSSoftwareDeploymentV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	orchestrationClient, err := config.orchestrationV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating RTS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_rts_software_deployment_v1" {
			continue
		}

		stack, err := softwaredeployment.Get(orchestrationClient, rs.Primary.ID).Extract()

		if err == nil {
			if stack.Status != "DELETE_COMPLETE" {
				return fmt.Errorf("Deployment still exists")
			}
		}
	}

	return nil
}

func testAccCheckOTCRTSSoftwareDeploymentV1Exists(n string, stack *softwaredeployment.Deployment) resource.TestCheckFunc {
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
			return fmt.Errorf("Error creating RTS Client : %s", err)
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

var testAccRTSSoftwareDeploymentV1_basic  =  fmt.Sprintf(`
resource "opentelekomcloud_compute_instance_v2" "vm_1" {
  name = "instance_1"
  network {
    uuid = "%s"
  }
}
resource "opentelekomcloud_rts_software_config" "config_1" {
  name="terraform-provider_test"
}

resource "opentelekomcloud_rts_software_deployment_v1" "deployment_1" {
  config_id = "${opentelekomcloud_rts_software_config.config_1.id}"
  server_id = "${opentelekomcloud_compute_instance_v2.vm_1.id}"
  status= "IN_PROGRESS"
  action= "CREATE"
  status_reason= "Deploy data"
}
`,OS_NETWORK_ID)

var testAccRTSSoftwareDeploymentV1_update  = fmt.Sprintf( `
resource "opentelekomcloud_compute_instance_v2" "vm_1" {
  name = "instance_1"
  network {
    uuid = "%s"
  }
}
resource "opentelekomcloud_rts_software_config" "config_1" {
  name="terraform-provider_test"
}

resource "opentelekomcloud_rts_software_deployment_v1" "deployment_1" {
  config_id = "${opentelekomcloud_rts_software_config.config_1.id}"
  server_id = "${opentelekomcloud_compute_instance_v2.vm_1.id}"
  output_values{
    deploy_stdout= "Writing to /tmp/baaaaa\nWritten to /tmp/baaaaa\n"
    deploy_stderr= "+ echo Writing to /tmp/baaaaa\n+ echo fooooo\n+ cat /tmp/baaaaa\n+ echo -n The file /tmp/baaaaa contains fooooo for server ec14c864-096e-4e27-bb8a-2c2b4dc6f3f5 during CREATE\n+ echo Written to /tmp/baaaaa\n+ echo Output to stderr\nOutput to stderr\n"
    deploy_status_code= 0
    result= "The file /tmp/baaaaa contains fooooo for server ec14c864-096e-4e27-bb8a-2c2b4dc6f3f5 during CREATE"
  }
  status= "COMPLETE"
  action= "CREATE"
  status_reason= "Outputs received"
}
`,OS_NETWORK_ID)

var testAccRTSSoftwareDeploymentV1_timeout = fmt.Sprintf( `
resource "opentelekomcloud_compute_instance_v2" "vm_1" {
  name = "instance_1"
  network {
    uuid = "%s"
  }
}
resource "opentelekomcloud_rts_software_config" "config_1" {
  name="terraform-provider_test"
}

resource "opentelekomcloud_rts_software_deployment_v1" "deployment_1" {
  config_id = "${opentelekomcloud_rts_software_config.config_1.id}"
  server_id = "${opentelekomcloud_compute_instance_v2.vm_1.id}"
  status= "COMPLETE"
  action= "CREATE"
  status_reason= "Outputs received"

  timeouts {
    create = "10m"
    delete = "10m"
  }
}
`,OS_NETWORK_ID)
