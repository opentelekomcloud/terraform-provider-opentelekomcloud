package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccVpcEipV1DataSource_basic(t *testing.T) {
	dataSourceNameByID := "data.opentelekomcloud_vpc_eip_v1.by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { common.TestAccPreCheck(t) },
		Providers: common.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcEipV1Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNameByID, "type", "5_bgp"),
					resource.TestCheckResourceAttr(dataSourceNameByID, "bandwidth_share_type", "PER"),
					resource.TestCheckResourceAttr(dataSourceNameByID, "status", "DOWN"),
				),
			},
		},
	})
}

const testAccDataSourceVpcEipV1Config = `
resource "opentelekomcloud_vpc_eip_v1" "eip" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "acc-band"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

data "opentelekomcloud_vpc_eip_v1" "by_id" {
  id = opentelekomcloud_vpc_eip_v1.eip.id
}
`
