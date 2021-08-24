package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccOTCRtsSoftwareDeploymentV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOTCRtsSoftwareDeploymentV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRtsSoftwareDeploymentV1DataSourceID("data.opentelekomcloud_rts_software_deployment_v1.deployment_1"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_rts_software_deployment_v1.deployment_1", "status_reason", "Deploy data"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_rts_software_deployment_v1.deployment_1", "action", "CREATE"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_rts_software_deployment_v1.deployment_1", "status", "COMPLETE"),
				),
			},
		},
	})
}

func testAccCheckOTCRtsSoftwareDeploymentV1DataSourceID(n string) resource.TestCheckFunc {
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

var testAccOTCRtsSoftwareDeploymentV1DataSource_basic = fmt.Sprintf(`
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
  status_reason= "Deploy data"
}

data "opentelekomcloud_rts_software_deployment_v1" "deployment_1" {
  id = opentelekomcloud_rts_software_deployment_v1.deployment_1.id
 }
`, env.OsImageID, env.OsFlavorID, env.OsNetworkID)
