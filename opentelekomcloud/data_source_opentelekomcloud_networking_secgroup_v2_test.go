package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOpenTelekomCloudNetworkingSecGroupV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenTelekomCloudNetworkingSecGroupV2DataSource_group,
			},
			{
				Config: testAccOpenTelekomCloudNetworkingSecGroupV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.opentelekomcloud_networking_secgroup_v2.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_networking_secgroup_v2.secgroup_1", "name", "secgroup_1"),
				),
			},
		},
	})
}

func TestAccOpenTelekomCloudNetworkingSecGroupV2DataSource_secGroupID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenTelekomCloudNetworkingSecGroupV2DataSource_group,
			},
			{
				Config: testAccOpenTelekomCloudNetworkingSecGroupV2DataSource_secGroupID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.opentelekomcloud_networking_secgroup_v2.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_networking_secgroup_v2.secgroup_1", "name", "secgroup_1"),
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

const testAccOpenTelekomCloudNetworkingSecGroupV2DataSource_group = `
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
        name        = "secgroup_1"
	description = "My neutron security group"
}
`

var testAccOpenTelekomCloudNetworkingSecGroupV2DataSource_basic = fmt.Sprintf(`
%s

data "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
	name = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.name}"
}
`, testAccOpenTelekomCloudNetworkingSecGroupV2DataSource_group)

var testAccOpenTelekomCloudNetworkingSecGroupV2DataSource_secGroupID = fmt.Sprintf(`
%s

data "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
	secgroup_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_1.id}"
}
`, testAccOpenTelekomCloudNetworkingSecGroupV2DataSource_group)
