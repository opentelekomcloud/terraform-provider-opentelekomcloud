package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v1/instances"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccDmsInstancesV1_basic(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dms_instance_%s", acctest.RandString(5))
	var instanceUpdate = fmt.Sprintf("dms_instance_update_%s", acctest.RandString(5))

	resourceName := "opentelekomcloud_dms_instance_v1.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1InstanceBasic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV1InstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
				),
			},
			{
				Config: testAccDmsV1InstanceUpdate(instanceUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV1InstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", "instance update description"),
				),
			},
		},
	})
}

func testAccCheckDmsV1InstanceDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.DmsV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DMSv1 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_dms_instance_v1" {
			continue
		}

		_, err := instances.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("DMS instance still exists")
		}
	}
	return nil
}

func testAccCheckDmsV1InstanceExists(n string, instance instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.DmsV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud DMSv1 client: %w", err)
		}

		v, err := instances.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return fmt.Errorf("error getting OpenTelekomCloud DMSv1 instance (%s): %w", rs.Primary.ID, err)
		}

		if v.InstanceID != rs.Primary.ID {
			return fmt.Errorf("DMS instance not found")
		}
		instance = *v
		return nil
	}
}

func testAccDmsV1InstanceBasic(instanceName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "secgroup_1"
}

data "opentelekomcloud_dms_az_v1" "az_1" {}

data "opentelekomcloud_dms_product_v1" "product_1" {
  engine        = "kafka"
  instance_type = "cluster"
  version       = "2.3.0"
}
resource "opentelekomcloud_dms_instance_v1" "instance_1" {
  name              = "%s"
  engine            = "kafka"
  storage_space     = data.opentelekomcloud_dms_product_v1.product_1.storage
  access_user       = "user"
  password          = "Dmstest@123"
  vpc_id            = "%s"
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_1.id
  subnet_id         = "%s"
  available_zones   = [data.opentelekomcloud_dms_az_v1.az_1.id]
  product_id        = data.opentelekomcloud_dms_product_v1.product_1.id
  engine_version    = data.opentelekomcloud_dms_product_v1.product_1.version
  storage_spec_code = data.opentelekomcloud_dms_product_v1.product_1.storage_spec_code
}
`, instanceName, env.OS_VPC_ID, env.OS_NETWORK_ID)
}

func testAccDmsV1InstanceUpdate(instanceUpdate string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "secgroup_1"
}

data "opentelekomcloud_dms_az_v1" "az_1" {}

data "opentelekomcloud_dms_product_v1" "product_1" {
  engine        = "kafka"
  instance_type = "cluster"
  version       = "2.3.0"
}

resource "opentelekomcloud_dms_instance_v1" "instance_1" {
  name              = "%s"
  description       = "instance update description"
  engine            = "kafka"
  storage_space     = data.opentelekomcloud_dms_product_v1.product_1.storage
  access_user       = "user"
  password          = "Dmstest@123"
  vpc_id            = "%s"
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup_1.id
  subnet_id         = "%s"
  available_zones   = [data.opentelekomcloud_dms_az_v1.az_1.id]
  product_id        = data.opentelekomcloud_dms_product_v1.product_1.id
  engine_version    = data.opentelekomcloud_dms_product_v1.product_1.version
  storage_spec_code = data.opentelekomcloud_dms_product_v1.product_1.storage_spec_code
}
`, instanceUpdate, env.OS_VPC_ID, env.OS_NETWORK_ID)
}
