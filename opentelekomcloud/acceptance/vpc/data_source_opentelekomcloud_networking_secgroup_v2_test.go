package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccOpenTelekomCloudNetworkingSecGroupV2DataSource_basic(t *testing.T) {
	t.Parallel()
	quotas.BookOne(t, quotas.SecurityGroup)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupV2DataSourceGroup,
			},
			{
				Config: testAccOpenTelekomCloudNetworkingSecGroupV2DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.opentelekomcloud_networking_secgroup_v2.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_networking_secgroup_v2.secgroup_1", "name", "secgroup_1_ds"),
				),
			},
			{
				Config: testAccOpenTelekomCloudNetworkingSecGroupV2DataSourceSecGroupID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.opentelekomcloud_networking_secgroup_v2.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_networking_secgroup_v2.secgroup_1", "name", "secgroup_1_ds"),
				),
			},
		},
	})
}

func TestAccOpenTelekomCloudNetworkingSecGroupV2DataSource_regex(t *testing.T) {
	dataSourceName := "data.opentelekomcloud_networking_secgroup_v2.secgroup_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupV2DataSourceGroup,
			},
			{
				Config: testAccOpenTelekomCloudNetworkingSecGroupV2DataSourceRegex,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", "secgroup_1_ds"),
				),
			},
		},
	})
}

func testAccCheckNetworkingSecGroupV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find security group data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("security group data source ID not set")
		}

		return nil
	}
}

const testAccNetworkingSecGroupV2DataSourceGroup = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1_ds"
  description = "My neutron security group"
}
`

var testAccOpenTelekomCloudNetworkingSecGroupV2DataSourceBasic = fmt.Sprintf(`
%s

data "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name = opentelekomcloud_networking_secgroup_v2.secgroup_1.name
}
`, testAccNetworkingSecGroupV2DataSourceGroup)

var testAccOpenTelekomCloudNetworkingSecGroupV2DataSourceSecGroupID = fmt.Sprintf(`
%s

data "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  secgroup_id = opentelekomcloud_networking_secgroup_v2.secgroup_1.id
}
`, testAccNetworkingSecGroupV2DataSourceGroup)

var testAccOpenTelekomCloudNetworkingSecGroupV2DataSourceRegex = fmt.Sprintf(`
%s

data "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name_regex  = "^secgroup_1.+"
}
`, testAccNetworkingSecGroupV2DataSourceGroup)
