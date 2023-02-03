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
	resourceName := "data.opentelekomcloud_compute_instance_v2.source_1"
	instanceName := "opentelekomcloud_compute_instance_v2.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceDataSourceBasic(),
			},
			{
				Config: testAccComputeV2InstanceDataSourceID(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceV2DataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "instance_1"),
					resource.TestCheckResourceAttrPair(resourceName, "metadata", instanceName, "metadata"),
					resource.TestCheckResourceAttrSet(resourceName, "network.0.name"),
				),
			},
			{
				Config: testAccComputeV2InstanceDataSourceWindowsPassword(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceV2DataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "instance_1"),
					resource.TestCheckResourceAttrPair(resourceName, "metadata", instanceName, "metadata"),
					resource.TestCheckResourceAttrSet(resourceName, "network.0.name"),
				),
			},
			{
				Config: testAccComputeV2InstanceDataSourceName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceV2DataSourceName(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "instance_1"),
					resource.TestCheckResourceAttrPair(resourceName, "metadata", instanceName, "metadata"),
					resource.TestCheckResourceAttrSet(resourceName, "network.0.name"),
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

func testAccCheckComputeInstanceV2DataSourceName(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find compute instance data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Compute instance data source name not set")
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
  flavor_name		= "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, getFlavorName())
}

func testAccComputeV2InstanceDataSourceWindows() string {
	return fmt.Sprintf(`
%s

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "instance_1"
  availability_zone = "%s"
  image_name        = "Enterprise_Windows_STD_2019_CORE_KVM"
  flavor_name		= "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  }
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, getFlavorName())
}

func testAccComputeV2InstanceDataSourceID() string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_compute_instance_v2" "source_1" {
  id = "${opentelekomcloud_compute_instance_v2.instance_1.id}"
}
`, testAccComputeV2InstanceDataSourceBasic())
}

func testAccComputeV2InstanceDataSourceName() string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_compute_instance_v2" "source_1" {
  name = "${opentelekomcloud_compute_instance_v2.instance_1.name}"
}
`, testAccComputeV2InstanceDataSourceBasic())
}

func testAccComputeV2InstanceDataSourceWindowsPassword() string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_compute_instance_v2" "source_1" {
  name = "${opentelekomcloud_compute_instance_v2.instance_1.name}"
}
`, testAccComputeV2InstanceDataSourceWindows())
}
