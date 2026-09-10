package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
)

func TestAccVpcSubnetPrivateIPV1DataSource_basic(t *testing.T) {
	const dataSourceName = "data.opentelekomcloud_vpc_subnet_private_ip_v1.test"

	t.Parallel()
	quotas.BookMany(t, vpcSubnetQuotas())

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetPrivateIPV1DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceVpcSubnetPrivateIPV1Name, "id"),
					resource.TestCheckResourceAttrPair(
						dataSourceName, "subnet_id", resourceVpcSubnetPrivateIPV1Name, "subnet_id",
					),
					resource.TestCheckResourceAttrPair(
						dataSourceName, "ip_address", resourceVpcSubnetPrivateIPV1Name, "ip_address",
					),
					resource.TestCheckResourceAttr(dataSourceName, "status", "DOWN"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tenant_id"),
					resource.TestCheckResourceAttr(dataSourceName, "device_owner", ""),
				),
			},
		},
	})
}

const testAccVpcSubnetPrivateIPV1DataSourceBasic = testAccVpcSubnetPrivateIPV1Basic + `
data "opentelekomcloud_vpc_subnet_private_ip_v1" "test" {
  id = opentelekomcloud_vpc_subnet_private_ip_v1.test.id
}
`
