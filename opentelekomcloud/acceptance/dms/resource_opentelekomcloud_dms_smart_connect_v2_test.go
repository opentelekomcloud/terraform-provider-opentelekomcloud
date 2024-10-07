package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const resourceDMSSmartConnectV2Name = "opentelekomcloud_dms_smart_connect_v2.test"

func TestAccDmsSmartConnectV2_basic(t *testing.T) {
	var instanceName = fmt.Sprintf("dms_instance_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV2SmartConnectBasic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceDMSSmartConnectV2Name, "instance_id"),
					resource.TestCheckResourceAttrSet(resourceDMSSmartConnectV2Name, "id"),
				),
			},
		},
	})
}

func testAccDmsV2SmartConnectBasic(instanceName string) string {
	return fmt.Sprintf(`
%s

%s

data "opentelekomcloud_dms_az_v1" "az_1" {}

data "opentelekomcloud_dms_flavor_v2" "test" {
  type      = "cluster"
  flavor_id = "c6.2u4g.cluster"
}

locals {
  flavor = data.opentelekomcloud_dms_flavor_v2.test.flavors[0]
}

resource "opentelekomcloud_dms_dedicated_instance_v2" "instance_1" {
  name        = "%s"
  description = "kafka test"

  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_id        = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id

  flavor_id         = local.flavor.id
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  available_zones   = [data.opentelekomcloud_dms_az_v1.az_1.id]
  engine_version    = "2.7"
  storage_space     = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  broker_num        = 3

  ssl_enable         = true
  access_user        = "user"
  password           = "Dmstest@123"
  security_protocol  = "SASL_PLAINTEXT"
  enabled_mechanisms = ["SCRAM-SHA-512"]

  cross_vpc_accesses {
    advertised_ip = ""
  }
  cross_vpc_accesses {
    advertised_ip = "www.terraform-test.com"
  }
  cross_vpc_accesses {
    advertised_ip = "192.168.0.53"
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "opentelekomcloud_dms_smart_connect_v2" "test" {
  instance_id = opentelekomcloud_dms_dedicated_instance_v2.instance_1.id
  node_count  = 2
}


`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, instanceName)
}
