package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/users"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceUserV2Name = "opentelekomcloud_dms_user_v2.user_1"

func TestAccDmsUsersV2_basic(t *testing.T) {
	var user users.Users
	var instanceName = fmt.Sprintf("dms_instance_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV2UserBasic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV2UserExists(resourceUserV2Name, user),
					resource.TestCheckResourceAttr(resourceUserV2Name, "username", "Test-user"),
					resource.TestCheckResourceAttr(resourceUserV2Name, "default_app", "false"),
					resource.TestCheckResourceAttr(resourceUserV2Name, "role", "guest"),
				),
			},
			{
				Config: testAccDmsV2UserPasswordUpdate(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV2UserExists(resourceUserV2Name, user),
					resource.TestCheckResourceAttr(resourceUserV2Name, "username", "Test-user"),
					resource.TestCheckResourceAttr(resourceUserV2Name, "default_app", "false"),
					resource.TestCheckResourceAttr(resourceUserV2Name, "role", "guest"),
				),
			},
		},
	})
}

func testAccCheckDmsV2UserExists(n string, user users.Users) resource.TestCheckFunc {
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

		v, err := users.List(client, rs.Primary.Attributes["instance_id"])
		if err != nil {
			return fmt.Errorf("error getting OpenTelekomCloud DMSv2 instance (%s): %w", rs.Primary.ID, err)
		}

		for _, userList := range v {
			if userList.UserName == rs.Primary.ID {
				user = userList
				break
			}
		}

		if user.UserName == "" {
			return fmt.Errorf("DMSv2 user '%s' doesn't exist", rs.Primary.ID)
		}

		return nil
	}
}

func testAccDmsV2UserBasic(instanceName string) string {
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

resource "opentelekomcloud_dms_user_v2" "user_1" {
  instance_id = opentelekomcloud_dms_instance_v2.instance_1.id
  username    = "Test-user"
  password    = "Dmstest@123"
}

`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, instanceName)
}

func testAccDmsV2UserPasswordUpdate(instanceUpdate string) string {
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
resource "opentelekomcloud_dms_user_v2" "user_1" {
  instance_id = opentelekomcloud_dms_instance_v2.instance_1.id
  username    = "Test-user"
  password    = "Dmstest@123@"
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, instanceUpdate)
}
