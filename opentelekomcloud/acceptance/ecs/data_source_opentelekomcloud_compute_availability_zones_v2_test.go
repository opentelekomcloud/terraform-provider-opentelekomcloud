package acceptance

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataAzName = "data.opentelekomcloud_compute_availability_zones_v2.zones"

func TestAccOpenStackAvailabilityZonesV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackAvailabilityZonesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataAzName, "names.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

const testAccOpenStackAvailabilityZonesConfig = `
data "opentelekomcloud_compute_availability_zones_v2" "zones" {}
`
