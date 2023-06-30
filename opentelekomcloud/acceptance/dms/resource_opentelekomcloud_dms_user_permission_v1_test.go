package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const resourceUserPermissionsV1Name = "opentelekomcloud_dms_user_permission_v1.perm_1"

func TestAccDmsUsersPermissionsV1_basic(t *testing.T) {
	var instanceName = fmt.Sprintf("dms_instance_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1UserPermissionsBasic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceUserPermissionsV1Name, "topic_name", "test-topic"),
					resource.TestCheckResourceAttr(resourceUserPermissionsV1Name, "policies.0.username", "Test-user"),
					resource.TestCheckResourceAttr(resourceUserPermissionsV1Name, "policies.0.access_policy", "all"),
					resource.TestCheckResourceAttr(resourceUserPermissionsV1Name, "policies.1.username", "Test-user2"),
					resource.TestCheckResourceAttr(resourceUserPermissionsV1Name, "policies.1.access_policy", "sub"),
				),
			},
			{
				Config: testAccDmsV1UserPermissionsUpdate(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceUserPermissionsV1Name, "topic_name", "test-topic"),
					resource.TestCheckResourceAttr(resourceUserPermissionsV1Name, "policies.0.username", "Test-user"),
					resource.TestCheckResourceAttr(resourceUserPermissionsV1Name, "policies.0.access_policy", "sub"),
				),
			},
		},
	})
}

func testAccDmsV1UserPermissionsBasic(instanceName string) string {
	return fmt.Sprintf(`
%s

%s

data "opentelekomcloud_dms_az_v1" "az_1" {}

data "opentelekomcloud_dms_product_v1" "product_1" {
  engine        = "kafka"
  instance_type = "cluster"
  version       = "2.3.0"
}

resource "opentelekomcloud_dms_instance_v2" "instance_1" {
  name              = "%s"
  engine            = "kafka"
  storage_space     = data.opentelekomcloud_dms_product_v1.product_1.storage
  access_user       = "user"
  password          = "Dmstest@123"
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  available_zones   = [data.opentelekomcloud_dms_az_v1.az_1.id]
  product_id        = data.opentelekomcloud_dms_product_v1.product_1.id
  engine_version    = data.opentelekomcloud_dms_product_v1.product_1.version
  storage_spec_code = data.opentelekomcloud_dms_product_v1.product_1.storage_spec_code
}

resource "opentelekomcloud_dms_topic_v1" "topic_1" {
  instance_id      = opentelekomcloud_dms_instance_v2.instance_1.id
  name             = "test-topic"
  partition        = 10
  replication      = 2
  sync_replication = true
  retention_time   = 720
}

resource "opentelekomcloud_dms_user_v2" "user_1" {
  instance_id = opentelekomcloud_dms_instance_v2.instance_1.id
  username    = "Test-user"
  password    = "Dmstest@123"
}

resource "opentelekomcloud_dms_user_v2" "user_2" {
  instance_id = opentelekomcloud_dms_instance_v2.instance_1.id
  username    = "Test-user2"
  password    = "Dmstest@123"
}

resource "opentelekomcloud_dms_user_permission_v1" "perm_1" {
  instance_id = opentelekomcloud_dms_instance_v2.instance_1.id
  topic_name  = "test-topic"
  policies {
    username      = opentelekomcloud_dms_user_v2.user_1.id
    access_policy = "all"
  }

  policies {
    username      = opentelekomcloud_dms_user_v2.user_2.id
    access_policy = "sub"
  }

}


`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, instanceName)
}

func testAccDmsV1UserPermissionsUpdate(instanceUpdate string) string {
	return fmt.Sprintf(`
%s

%s

data "opentelekomcloud_dms_az_v1" "az_1" {}

data "opentelekomcloud_dms_product_v1" "product_1" {
  engine        = "kafka"
  instance_type = "cluster"
  version       = "2.3.0"
}

resource "opentelekomcloud_dms_instance_v2" "instance_1" {
  name              = "%s"
  engine            = "kafka"
  storage_space     = data.opentelekomcloud_dms_product_v1.product_1.storage
  access_user       = "user"
  password          = "Dmstest@123"
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  available_zones   = [data.opentelekomcloud_dms_az_v1.az_1.id]
  product_id        = data.opentelekomcloud_dms_product_v1.product_1.id
  engine_version    = data.opentelekomcloud_dms_product_v1.product_1.version
  storage_spec_code = data.opentelekomcloud_dms_product_v1.product_1.storage_spec_code
}

resource "opentelekomcloud_dms_topic_v1" "topic_1" {
  instance_id      = opentelekomcloud_dms_instance_v2.instance_1.id
  name             = "test-topic"
  partition        = 10
  replication      = 2
  sync_replication = true
  retention_time   = 720
}

resource "opentelekomcloud_dms_user_v2" "user_1" {
  instance_id = opentelekomcloud_dms_instance_v2.instance_1.id
  username    = "Test-user"
  password    = "Dmstest@123"
}

resource "opentelekomcloud_dms_user_permission_v1" "perm_1" {
  instance_id = opentelekomcloud_dms_instance_v2.instance_1.id
  topic_name  = "test-topic"
  policies {
    username      = opentelekomcloud_dms_user_v2.user_1.id
    access_policy = "sub"
  }
}

`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, instanceUpdate)
}
