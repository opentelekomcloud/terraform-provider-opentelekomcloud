package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataFlavorName = "data.opentelekomcloud_lb_flavor_v3.l7_s2_small"

func TestLBFlavorV3_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testLBFlavorV3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataFlavorName, "id"),
					resource.TestCheckResourceAttr(dataFlavorName, "qps", "8000"),
					resource.TestCheckResourceAttr(dataFlavorName, "cps", "4000"),
					resource.TestCheckResourceAttr(dataFlavorName, "bandwidth", "100000"),
				),
			},
		},
	})
}

func TestLBFlavorV3_byID(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testLBFlavorV3ID,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataFlavorName, "qps", "8000"),
					resource.TestCheckResourceAttr(dataFlavorName, "cps", "4000"),
					resource.TestCheckResourceAttr(dataFlavorName, "bandwidth", "100000"),
				),
			},
		},
	})
}

const testLBFlavorV3 = `
data opentelekomcloud_lb_flavor_v3 l7_s2_small {
  name = "L7_flavor.elb.s2.small"
}
`

const testLBFlavorV3ID = `
data opentelekomcloud_lb_flavor_v3 l7_s2_small {
  id = "7628a037-c229-4c0c-820b-4ea862743aef"
}
`
