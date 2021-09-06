package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccRTSSoftwareDeploymentV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccSubnetPreCheck(t)
			common.TestAccImagePreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRTSSoftwareDeploymentV1DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRTSSoftwareDeploymentV1DataSourceID("data.opentelekomcloud_rts_software_deployment_v1.deployment_1"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_rts_software_deployment_v1.deployment_1", "status_reason", "Deploy data"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_rts_software_deployment_v1.deployment_1", "action", "CREATE"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_rts_software_deployment_v1.deployment_1", "status", "COMPLETE"),
				),
			},
		},
	})
}

func testAccCheckRTSSoftwareDeploymentV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find Software Deployment data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("RTS Software Deployment data source ID not set ")
		}

		return nil
	}
}

var testAccRTSSoftwareDeploymentV1DataSourceBasic = fmt.Sprintf(`
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
  status_reason= "Deploy data"
}

data "opentelekomcloud_rts_software_deployment_v1" "deployment_1" {
  id = opentelekomcloud_rts_software_deployment_v1.deployment_1.id
 }
`, common.DataSourceSubnet, common.DataSourceImage, env.OsFlavorID)
