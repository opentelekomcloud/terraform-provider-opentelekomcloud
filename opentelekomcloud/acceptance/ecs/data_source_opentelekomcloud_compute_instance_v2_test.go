package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccComputeV2InstanceDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceDataSourceBasic(),
			},
			{
				Config: testAccComputeV2InstanceDataSourceSource(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceV2DataSourceID("data.opentelekomcloud_compute_instance_v2.source_1"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_compute_instance_v2.source_1", "name", "instance_1"),
					resource.TestCheckResourceAttrPair("data.opentelekomcloud_compute_instance_v2.source_1", "metadata", "opentelekomcloud_compute_instance_v2.instance_1", "metadata"),
					resource.TestCheckResourceAttrSet("data.opentelekomcloud_compute_instance_v2.source_1", "network.0.name"),
				),
			},
		},
	})
}

func testAccCheckComputeInstanceV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find compute instance data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Compute instance data source ID not set")
		}

		return nil
	}
}

func testAccComputeV2InstanceDataSourceBasic() string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  availability_zone = "%s"
  image_name        = "Standard_Debian_10_latest"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)
}

func testAccComputeV2InstanceDataSourceSource() string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_compute_instance_v2" "source_1" {
  id = "${opentelekomcloud_compute_instance_v2.instance_1.id}"
}
`, testAccComputeV2InstanceDataSourceBasic())
}
