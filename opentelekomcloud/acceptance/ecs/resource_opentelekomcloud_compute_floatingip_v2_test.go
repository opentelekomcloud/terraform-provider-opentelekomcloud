package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/floatingips"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceFloatingIpName = "opentelekomcloud_compute_floatingip_v2.fip_1"

func TestAccComputeV2FloatingIP_basic(t *testing.T) {
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2FloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FloatingIPBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FloatingIPExists(resourceFloatingIpName, &fip),
				),
			},
		},
	})
}

func testAccCheckComputeV2FloatingIPDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ComputeV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud ComputeV2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_compute_floatingip_v2" {
			continue
		}

		_, err := floatingips.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("floatingIP still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2FloatingIPExists(n string, kp *floatingips.FloatingIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.ComputeV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud ComputeV2 client: %s", err)
		}

		found, err := floatingips.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("floatingIP not found")
		}

		*kp = *found

		return nil
	}
}

const testAccComputeV2FloatingIPBasic = `
resource "opentelekomcloud_compute_floatingip_v2" "fip_1" {}
`
