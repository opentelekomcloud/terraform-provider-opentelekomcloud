package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/instances/lifecycle"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceInstanceV2Name = "opentelekomcloud_dms_instance_v2.instance_1"

func TestAccDmsInstancesV2_basic(t *testing.T) {
	var instance lifecycle.Instance
	var instanceName = fmt.Sprintf("dms_instance_%s", acctest.RandString(5))
	var instanceUpdate = fmt.Sprintf("dms_instance_update_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV2InstanceBasic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV2InstanceExists(resourceInstanceV2Name, instance),
					resource.TestCheckResourceAttr(resourceInstanceV2Name, "name", instanceName),
					resource.TestCheckResourceAttr(resourceInstanceV2Name, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceInstanceV2Name, "tags.foo", "bar"),
				),
			},
			{
				Config: testAccDmsV2InstanceUpdate(instanceUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV2InstanceExists(resourceInstanceV2Name, instance),
					resource.TestCheckResourceAttr(resourceInstanceV2Name, "name", instanceUpdate),
					resource.TestCheckResourceAttr(resourceInstanceV2Name, "description", "instance update description"),
					resource.TestCheckResourceAttr(resourceInstanceV2Name, "tags.new_test", "new_test2"),
					resource.TestCheckResourceAttr(resourceInstanceV2Name, "tags.john", "doe"),
				),
			},
		},
	})
}

// Not supported on current version
func TestAccDmsInstancesV2_Encrypted(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Disk encryption is not supported in current version")
	}
	var instance lifecycle.Instance
	var instanceName = fmt.Sprintf("dms_instance_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV2InstanceEncrypted(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV2InstanceExists(resourceInstanceV2Name, instance),
					resource.TestCheckResourceAttr(resourceInstanceV2Name, "name", instanceName),
					resource.TestCheckResourceAttr(resourceInstanceV2Name, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceInstanceV2Name, "disk_encrypted_enable", "true"),
				),
			},
		},
	})
}

func TestAccDmsInstancesV2_EIP(t *testing.T) {
	var instance lifecycle.Instance
	var instanceName = fmt.Sprintf("dms_instance_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV2InstanceEIP(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV2InstanceExists(resourceInstanceV2Name, instance),
					resource.TestCheckResourceAttr(resourceInstanceV2Name, "name", instanceName),
					resource.TestCheckResourceAttr(resourceInstanceV2Name, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceInstanceV2Name, "enable_publicip", "true"),
				),
			},
		},
	})
}

func testAccCheckDmsV2InstanceDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.DmsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DMSv2 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_dms_instance_v2" {
			continue
		}

		_, err := lifecycle.Get(client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("DMS instance still exists")
		}
	}
	return nil
}

func testAccCheckDmsV2InstanceExists(n string, instance lifecycle.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.DmsV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud DMSv2 client: %w", err)
		}

		v, err := lifecycle.Get(client, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error getting OpenTelekomCloud DMSv2 instance (%s): %w", rs.Primary.ID, err)
		}

		if v.InstanceID != rs.Primary.ID {
			return fmt.Errorf("DMS instance not found")
		}
		instance = *v
		return nil
	}
}

func testAccDmsV2InstanceBasic(instanceName string) string {
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

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, instanceName)
}

func testAccDmsV2InstanceUpdate(instanceUpdate string) string {
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
  description       = "instance update description"
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

  tags = {
    foo      = "bar"
    john     = "doe"
    new_test = "new_test2"
  }
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, instanceUpdate)
}

func testAccDmsV2InstanceEncrypted(instanceName string) string {
	return fmt.Sprintf(`
%s

%s

data "opentelekomcloud_dms_az_v1" "az_1" {}

data "opentelekomcloud_dms_product_v1" "product_1" {
  engine        = "kafka"
  instance_type = "cluster"
  version       = "1.1.0"
}

resource "opentelekomcloud_dms_instance_v2" "instance_1" {
  name                  = "%s"
  engine                = "kafka"
  storage_space         = data.opentelekomcloud_dms_product_v1.product_1.storage
  vpc_id                = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  security_group_id     = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id             = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  available_zones       = [data.opentelekomcloud_dms_az_v1.az_1.id]
  product_id            = data.opentelekomcloud_dms_product_v1.product_1.id
  engine_version        = data.opentelekomcloud_dms_product_v1.product_1.version
  storage_spec_code     = data.opentelekomcloud_dms_product_v1.product_1.storage_spec_code
  disk_encrypted_enable = true
  disk_encrypted_key    = "%s"
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, instanceName, env.OS_KMS_ID)
}

func testAccDmsV2InstanceEIP(instanceName string) string {
	return fmt.Sprintf(`
%s

%s

data "opentelekomcloud_dms_az_v1" "az_1" {}

data "opentelekomcloud_dms_product_v1" "product_1" {
  engine        = "kafka"
  instance_type = "cluster"
  version       = "1.1.0"
}

resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {
}

resource "opentelekomcloud_networking_floatingip_v2" "fip_2" {
}

resource "opentelekomcloud_networking_floatingip_v2" "fip_3" {
}

resource "opentelekomcloud_dms_instance_v2" "instance_1" {
  name              = "%s"
  engine            = "kafka"
  storage_space     = data.opentelekomcloud_dms_product_v1.product_1.storage
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  available_zones   = [data.opentelekomcloud_dms_az_v1.az_1.id]
  product_id        = data.opentelekomcloud_dms_product_v1.product_1.id
  engine_version    = data.opentelekomcloud_dms_product_v1.product_1.version
  storage_spec_code = data.opentelekomcloud_dms_product_v1.product_1.storage_spec_code
  enable_publicip   = true
  publicip_id = [opentelekomcloud_networking_floatingip_v2.fip_1.id,
    opentelekomcloud_networking_floatingip_v2.fip_2.id,
  opentelekomcloud_networking_floatingip_v2.fip_3.id]
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, instanceName)
}
