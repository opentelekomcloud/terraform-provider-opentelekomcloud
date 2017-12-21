package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// PASS
func TestAccOpenTelekomCloudVpcV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOpenTelekomCloudVpcV1DataSource_resource,
			},
			resource.TestStep{
				Config: testAccOpenTelekomCloudVpcV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1DataSourceID("data.opentelekomcloud_vpc_v1.vpc_data"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_v1.vpc_data", "name", "terraform_provider_test"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_v1.vpc_data", "cidr", "192.168.0.0/16"),
				),
			},
		},
	})
}

// PASS
func TestAccOpenTelekomCloudVpcV1DataSource_vpcID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOpenTelekomCloudVpcV1DataSource_resource,
			},
			resource.TestStep{
				Config: testAccOpenTelekomCloudVpcV1DataSource_vpcID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1DataSourceID("data.opentelekomcloud_vpc_v1.vpc_data"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_v1.vpc_data", "name", "terraform_provider_test"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_vpc_v1.vpc_data", "cidr", "192.168.0.0/16"),
				),
			},
		},
	})
}

func testAccCheckVpcV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find vpc data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Vpc data source ID not set")
		}

		return nil
	}
}

const testAccOpenTelekomCloudVpcV1DataSource_resource = `
resource "opentelekomcloud_vpc_v1" "vpc_1" {
	name = "terraform_provider_test"
	cidr="192.168.0.0/16"
}
`

var testAccOpenTelekomCloudVpcV1DataSource_basic = fmt.Sprintf(`
%s

data "opentelekomcloud_vpc_v1" "vpc_data" {
	name = "${opentelekomcloud_vpc_v1.vpc_1.name}"
}
`, testAccOpenTelekomCloudVpcV1DataSource_resource)


var testAccOpenTelekomCloudVpcV1DataSource_vpcID = fmt.Sprintf(`
%s

data "opentelekomcloud_vpc_v1" "vpc_data" {
	id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
}
`, testAccOpenTelekomCloudVpcV1DataSource_resource)
