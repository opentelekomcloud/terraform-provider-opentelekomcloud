package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccComputeV2InstancesDataSource_basic(t *testing.T) {
	resourceName := "data.opentelekomcloud_compute_instances_v2.source_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstancesDataSourceBasic(),
			},
			{
				Config: testAccComputeV2InstancesDataSourceName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstancesV2DataSourceName(resourceName),
					resource.TestCheckResourceAttr(resourceName, "instances.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "instances.0.name", "instance_1"),
				),
			},
			{
				Config: testAccComputeV2InstancesDataSourceProject(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstancesV2DataSourceName(resourceName),
					resource.TestCheckResourceAttr(resourceName, "instances.#", "1"),
				),
			},
			{
				Config: testAccComputeV2InstancesDataSourceStatus(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstancesV2DataSourceName(resourceName),
					resource.TestCheckResourceAttr(resourceName, "instances.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "instances.0.status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckComputeInstancesV2DataSourceName(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find compute instance data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("compute instance data source name not set")
		}

		return nil
	}
}

func testAccComputeV2InstancesDataSourceBasic() string {
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

func testAccComputeV2InstancesDataSourceName() string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_compute_instances_v2" "source_1" {
  name = opentelekomcloud_compute_instance_v2.instance_1.name
}
`, testAccComputeV2InstancesDataSourceBasic())
}

func testAccComputeV2InstancesDataSourceProject() string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_identity_project_v3" "project" {
  name = "%s"
}

data "opentelekomcloud_compute_instances_v2" "source_1" {
  project_id = data.opentelekomcloud_identity_project_v3.project.id
}
`, testAccComputeV2InstancesDataSourceBasic(), env.OS_TENANT_NAME)
}

func testAccComputeV2InstancesDataSourceStatus() string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_compute_instances_v2" "source_1" {
  status = "ACTIVE"
}
`, testAccComputeV2InstancesDataSourceBasic())
}
