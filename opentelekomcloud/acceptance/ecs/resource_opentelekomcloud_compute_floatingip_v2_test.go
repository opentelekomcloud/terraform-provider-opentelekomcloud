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

func TestAccComputeV2FloatingIP_basic(t *testing.T) {
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2FloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FloatingIP_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FloatingIPExists("opentelekomcloud_compute_floatingip_v2.fip_1", &fip),
				),
			},
		},
	})
}

func testAccCheckComputeV2FloatingIPDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	computeClient, err := config.ComputeV2Client(env.OsRegionName)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_compute_floatingip_v2" {
			continue
		}

		_, err := floatingips.Get(computeClient, rs.Primary.ID).Extract()
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
		computeClient, err := config.ComputeV2Client(env.OsRegionName)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud compute client: %s", err)
		}

		found, err := floatingips.Get(computeClient, rs.Primary.ID).Extract()
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

const testAccComputeV2FloatingIP_basic = `
resource "opentelekomcloud_compute_floatingip_v2" "fip_1" {}
`
