package hwcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccHWCloudNetworkingSecGroupV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccHWCloudNetworkingSecGroupV2DataSource_group,
			},
			resource.TestStep{
				Config: testAccHWCloudNetworkingSecGroupV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.hwcloud_networking_secgroup_v2.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.hwcloud_networking_secgroup_v2.secgroup_1", "name", "secgroup_1"),
				),
			},
		},
	})
}

func TestAccHWCloudNetworkingSecGroupV2DataSource_secGroupID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccHWCloudNetworkingSecGroupV2DataSource_group,
			},
			resource.TestStep{
				Config: testAccHWCloudNetworkingSecGroupV2DataSource_secGroupID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.hwcloud_networking_secgroup_v2.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.hwcloud_networking_secgroup_v2.secgroup_1", "name", "secgroup_1"),
				),
			},
		},
	})
}

func testAccCheckNetworkingSecGroupV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find security group data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Security group data source ID not set")
		}

		return nil
	}
}

const testAccHWCloudNetworkingSecGroupV2DataSource_group = `
resource "hwcloud_networking_secgroup_v2" "secgroup_1" {
        name        = "secgroup_1"
	description = "My neutron security group"
}
`

var testAccHWCloudNetworkingSecGroupV2DataSource_basic = fmt.Sprintf(`
%s

data "hwcloud_networking_secgroup_v2" "secgroup_1" {
	name = "${hwcloud_networking_secgroup_v2.secgroup_1.name}"
}
`, testAccHWCloudNetworkingSecGroupV2DataSource_group)

var testAccHWCloudNetworkingSecGroupV2DataSource_secGroupID = fmt.Sprintf(`
%s

data "hwcloud_networking_secgroup_v2" "secgroup_1" {
	secgroup_id = "${hwcloud_networking_secgroup_v2.secgroup_1.id}"
}
`, testAccHWCloudNetworkingSecGroupV2DataSource_group)
