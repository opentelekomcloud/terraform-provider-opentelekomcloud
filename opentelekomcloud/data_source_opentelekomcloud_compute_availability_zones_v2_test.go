package opentelekomcloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOpenStackAvailabilityZonesV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackAvailabilityZonesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.opentelekomcloud_compute_availability_zones_v2.zones", "names.#", regexp.MustCompile("[1-9]\\d*")),
				),
			},
		},
	})
}

const testAccOpenStackAvailabilityZonesConfig = `
data "opentelekomcloud_compute_availability_zones_v2" "zones" {}
`
