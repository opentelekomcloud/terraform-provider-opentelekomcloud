package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/instances/lifecycle"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const resourceDedicatedInstanceV2Name = "opentelekomcloud_dms_dedicated_instance_v2.instance_1"

func TestAccDmsDedicatedInstancesV2_basic(t *testing.T) {
	var instance lifecycle.Instance
	var instanceName = fmt.Sprintf("dms_instance_%s", acctest.RandString(5))
	var instanceUpdate = fmt.Sprintf("dms_instance_update_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV2DedicatedInstanceBasic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV2InstanceExists(resourceDedicatedInstanceV2Name, instance),
					resource.TestCheckResourceAttr(resourceDedicatedInstanceV2Name, "name", instanceName),
					resource.TestCheckResourceAttr(resourceDedicatedInstanceV2Name, "description", "kafka test"),
					resource.TestCheckResourceAttr(resourceDedicatedInstanceV2Name, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceDedicatedInstanceV2Name, "tags.foo", "bar"),
				),
			},
			{
				Config: testAccDmsV2DedicatedInstanceUpdate(instanceUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV2InstanceExists(resourceDedicatedInstanceV2Name, instance),
					resource.TestCheckResourceAttr(resourceDedicatedInstanceV2Name, "name", instanceUpdate),
					resource.TestCheckResourceAttr(resourceDedicatedInstanceV2Name, "description", "kafka test updated"),
					resource.TestCheckResourceAttr(resourceDedicatedInstanceV2Name, "tags.new", "test"),
					resource.TestCheckResourceAttr(resourceDedicatedInstanceV2Name, "tags.foo", "bar"),
				),
			},
		},
	})
}

func TestAccDmsDedicatedInstancesV2_Advanced(t *testing.T) {
	var instance lifecycle.Instance
	var instanceName = fmt.Sprintf("dms_instance_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV2DedicatedInstanceAdvanced(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV2InstanceExists(resourceDedicatedInstanceV2Name, instance),
					resource.TestCheckResourceAttr(resourceDedicatedInstanceV2Name, "name", instanceName),
					resource.TestCheckResourceAttr(resourceDedicatedInstanceV2Name, "maintain_begin", "02:00"),
					resource.TestCheckResourceAttr(resourceDedicatedInstanceV2Name, "maintain_end", "06:00"),
					resource.TestCheckResourceAttr(resourceDedicatedInstanceV2Name, "retention_policy", "time_base"),
					resource.TestCheckResourceAttr(resourceDedicatedInstanceV2Name, "node_num", "4"),
				),
			},
			{
				ResourceName:      resourceDedicatedInstanceV2Name,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"used_storage_space",
					"cross_vpc_accesses",
				},
			},
		},
	})
}

func testAccDmsV2DedicatedInstanceBasic(instanceName string) string {
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
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, instanceName)
}

func testAccDmsV2DedicatedInstanceUpdate(instanceUpdate string) string {
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
  description = "kafka test updated"

  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  network_id        = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id

  flavor_id         = local.flavor.id
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  available_zones   = [data.opentelekomcloud_dms_az_v1.az_1.id]
  engine_version    = "2.7"
  storage_space     = 600
  broker_num        = 4

  ssl_enable         = true
  access_user        = "user"
  password           = "Dmstest@123"
  security_protocol  = "SASL_PLAINTEXT"
  enabled_mechanisms = ["SCRAM-SHA-512"]

  cross_vpc_accesses {
    advertised_ip = "192.168.0.61"
  }
  cross_vpc_accesses {
    advertised_ip = "test.terraform.com"
  }
  cross_vpc_accesses {
    advertised_ip = "192.168.0.62"
  }

  cross_vpc_accesses {
    advertised_ip = "192.168.0.63"
  }

  tags = {
    foo = "bar"
    key = "value"
    new = "test"
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, instanceUpdate)
}

func testAccDmsV2DedicatedInstanceAdvanced(instanceName string) string {
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
  storage_space     = 600
  broker_num        = 4

  ssl_enable       = true
  access_user      = "user"
  password         = "Dmstest@123"
  maintain_begin   = "02:00"
  maintain_end     = "06:00"
  retention_policy = "time_base"
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, instanceName)
}
