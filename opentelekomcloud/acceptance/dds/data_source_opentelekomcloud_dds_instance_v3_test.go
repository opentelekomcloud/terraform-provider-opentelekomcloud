package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

const dataInstanceName = "data.opentelekomcloud_dds_instance_v3.instances"

func TestAccDDSInstanceV3DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSInstanceV3DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSInstanceV3DataSourceID(dataInstanceName),
					resource.TestCheckResourceAttr(dataInstanceName, "name", "dds-instance"),
					resource.TestCheckResourceAttrSet(dataInstanceName, "vpc_id"),
					resource.TestCheckResourceAttr(dataInstanceName, "mode", "ReplicaSet"),
					resource.TestCheckResourceAttr(dataInstanceName, "datastore.0.type", "DDS-Community"),
				),
			},
		},
	})
}

func testAccCheckDDSInstanceV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find instances data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("node data source ID not set ")
		}

		return nil
	}
}

var testAccDDSInstanceV3DataSourceBasic = fmt.Sprintf(`
%s

%s

resource "opentelekomcloud_dds_instance_v3" "instance_1" {
  name              = "dds-instance"
  availability_zone = "%s"
  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  password          = "5ecuredPa55w0rd@"
  mode              = "ReplicaSet"
  flavor {
    type      = "replica"
    num       = 1
    size      = 20
    spec_code = "dds.mongodb.s2.medium.4.repset"
  }
}

data "opentelekomcloud_dds_instance_v3" "instances" {
  instance_id = opentelekomcloud_dds_instance_v3.instance_1.id
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE)
