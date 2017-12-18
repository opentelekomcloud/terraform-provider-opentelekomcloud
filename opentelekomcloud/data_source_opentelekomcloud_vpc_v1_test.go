package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// PASS
func TestAccOpenTelekomCloudNetworkingVpcV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOpenTelekomCloudNetworkingVpcV1DataSource_group,
			},
			resource.TestStep{
				Config: testAccOpenTelekomCloudNetworkingVpcV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingVpcV1DataSourceID("data.opentelekomcloud_networking_vpc_v1.network_data_vpc_v1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_networking_vpc_v1.network_data_vpc_v1", "name", "testvpc"),
				),
			},
		},
	})
}

// PASS
func TestAccOpenTelekomCloudNetworkingVpcV1DataSource_vpcID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOpenTelekomCloudNetworkingVpcV1DataSource_group,
			},
			resource.TestStep{
				Config: testAccOpenTelekomCloudNetworkingVpcV1DataSource_vpcID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingVpcV1DataSourceID("data.opentelekomcloud_networking_vpc_v1.network_data_vpc_v1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_networking_vpc_v1.network_data_vpc_v1", "name", "testvpc"),
				),
			},
		},
	})
}

func testAccCheckNetworkingVpcV1DataSourceID(n string) resource.TestCheckFunc {
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

const testAccOpenTelekomCloudNetworkingVpcV1DataSource_group = `
resource "opentelekomcloud_networking_vpc_v1" "networking_vpc_v1" {
        name        = "testvpc"
	description = "My test vpc"
}
`

var testAccOpenTelekomCloudNetworkingVpcV1DataSource_basic = fmt.Sprintf(`
%s

data "opentelekomcloud_networking_vpc_v1" "network_data_vpc_v1" {
	name = "${opentelekomcloud_networking_vpc_v1.network_data_vpc_v1.name}"
}
`, testAccOpenTelekomCloudNetworkingVpcV1DataSource_group)

var testAccOpenTelekomCloudNetworkingVpcV1DataSource_vpcID = fmt.Sprintf(`
%s

data "opentelekomcloud_networking_vpc_v1" "network_data_vpc_v1" {
	vpc_id = "${opentelekomcloud_networking_vpc_v1.network_data_vpc_v1.id}"
}
`, testAccOpenTelekomCloudNetworkingVpcV1DataSource_group)

