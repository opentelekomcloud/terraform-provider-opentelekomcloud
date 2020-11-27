package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDDSInstanceV3DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSInstanceV3DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSInstanceV3DataSourceID("data.opentelekomcloud_dds_instance_v3.instances"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_dds_instance_v3.instances", "name", "dds-instance"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_dds_instance_v3.instances", "vpc_id", OS_VPC_ID),
					resource.TestCheckResourceAttr("data.opentelekomcloud_dds_instance_v3.instances", "mode", "ReplicaSet"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_dds_instance_v3.instances", "datastore.0.type", "DDS-Community"),
				),
			},
		},
	})
}

func testAccCheckDDSInstanceV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find instances data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Node data source ID not set ")
		}

		return nil
	}
}

var testAccDDSInstanceV3DataSource_basic = fmt.Sprintf(`
resource "opentelekomcloud_networking_secgroup_v2" "sg_acc" {
  name = "secgroup_acc"
}
resource "opentelekomcloud_dds_instance_v3" "instance_1" {
  name              = "dds-instance"
  availability_zone = "%s"
  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }
  vpc_id            = "%s"
  subnet_id         = "%s"
  security_group_id = opentelekomcloud_networking_secgroup_v2.sg_acc.id
  password          = "5ecuredPa55w0rd@"
  mode              = "ReplicaSet"
  flavor {
    type = "replica"
    num = 1
    size = 20
    spec_code = "dds.mongodb.s2.medium.4.repset"
  }
}

data "opentelekomcloud_dds_instance_v3" "instances" {
  instance_id = opentelekomcloud_dds_instance_v3.instance_1.id
}
`, OS_AVAILABILITY_ZONE, OS_VPC_ID, OS_NETWORK_ID)
