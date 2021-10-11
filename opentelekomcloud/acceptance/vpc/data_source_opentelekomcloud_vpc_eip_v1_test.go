package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
)

func TestAccVpcEipV1DataSource_basic(t *testing.T) {
	dataSourceNameByID := "data.opentelekomcloud_vpc_eip_v1.by_id"
	dataSourceNameByTags := "data.opentelekomcloud_vpc_eip_v1.by_tags"
	t.Parallel()
	th.AssertNoErr(t, quotas.FloatingIP.Acquire())
	defer quotas.FloatingIP.Release()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcEipV1Init,
			},
			{
				Config: testAccDataSourceVpcEipV1Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNameByID, "type", "5_bgp"),
					resource.TestCheckResourceAttr(dataSourceNameByID, "bandwidth_share_type", "PER"),
					resource.TestCheckResourceAttr(dataSourceNameByID, "status", "DOWN"),
					resource.TestCheckResourceAttr(dataSourceNameByTags, "type", "5_bgp"),
					resource.TestCheckResourceAttr(dataSourceNameByTags, "bandwidth_share_type", "PER"),
					resource.TestCheckResourceAttr(dataSourceNameByTags, "status", "DOWN"),
				),
			},
			{
				Config: testAccDataSourceVpcEipV1Init,
			},
		},
	})
}

const testAccDataSourceVpcEipV1Init = `
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

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}
`

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

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}

data "opentelekomcloud_vpc_eip_v1" "by_id" {
  id = opentelekomcloud_vpc_eip_v1.eip.id
}

data "opentelekomcloud_vpc_eip_v1" "by_tags" {
  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}
`
