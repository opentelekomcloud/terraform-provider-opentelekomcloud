package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const vpcIPAddressGroupDataSourceName = "data.opentelekomcloud_vpc_ip_address_group_v3.group_1"

func TestAccVpcIPAddressGroupV3DS_basic(t *testing.T) {
	dc := common.InitDataSourceCheck(vpcIPAddressGroupDataSourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcIPAddressGroupV3DSBasic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(vpcIPAddressGroupDataSourceName, "name", "test-acc-ip-address-group-v3"),
					resource.TestCheckResourceAttr(vpcIPAddressGroupDataSourceName, "description", "created"),
					resource.TestCheckResourceAttr(vpcIPAddressGroupDataSourceName, "ip_version", "4"),
				),
			},
		},
	})
}

var testAccVpcIPAddressGroupV3DSBasic = `
resource "opentelekomcloud_vpc_ip_address_group_v3" "group_acc" {
  name        = "test-acc-ip-address-group-v3"
  description = "created"
  ip_version  = 4
}

data "opentelekomcloud_vpc_ip_address_group_v3" "group_1" {
  id = opentelekomcloud_vpc_ip_address_group_v3.group_acc.id
}
`
