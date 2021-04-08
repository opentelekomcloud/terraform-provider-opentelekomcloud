package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/sdrs/v1/protectedinstances"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccSdrsProtectedInstanceV1_basic(t *testing.T) {
	var intance protectedinstances.Instance
	resourceName := "opentelekomcloud_sdrs_protected_instance_v1.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccSdrsProtectedInstanceV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSdrsProtectedInstanceV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccSdrsProtectedInstanceV1Exists(resourceName, &intance),
					resource.TestCheckResourceAttr(resourceName, "name", "instance_create"),
					resource.TestCheckResourceAttr(resourceName, "description", "some interesting description"),
					resource.TestCheckResourceAttr(resourceName, "tags.muh", "value-create"),
				),
			},
			{
				Config: testAccSdrsProtectedInstanceV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccSdrsProtectedInstanceV1Exists(resourceName, &intance),
					resource.TestCheckResourceAttr(resourceName, "name", "instance_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.muh", "value-update"),
				),
			},
		},
	})
}

func testAccSdrsProtectedInstanceV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.SdrsV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud SDRS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_sdrs_protected_instance_v1" {
			continue
		}

		_, err := protectedinstances.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("SDRS protected instance still exists")
		}
	}

	return nil
}

func testAccSdrsProtectedInstanceV1Exists(n string, instance *protectedinstances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.SdrsV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud SDRS client: %s", err)
		}

		found, err := protectedinstances.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("SDRS protectiongroup not found")
		}

		*instance = *found

		return nil
	}
}

var testAccSdrsProtectedInstanceV1_basic = fmt.Sprintf(`
resource "opentelekomcloud_sdrs_protectiongroup_v1" "group_1" {
  name                     = "group_1"
  source_availability_zone = "%[1]s"
  target_availability_zone = "eu-de-01"
  domain_id                = "cdba26b2-cc35-4988-a904-82b7abf20094"
  source_vpc_id            = "%[2]s"
  dr_type                  = "migration"
}

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = "%[3]s"
  flavor   = "s2.medium.1"
  vpc_id   = "%[2]s"

  nics {
    network_id = "%[4]s"
  }

  availability_zone = "%[1]s"
}

resource "opentelekomcloud_sdrs_protected_instance_v1" "instance_1" {
  name                 = "instance_create"
  description          = "some interesting description"
  group_id             = opentelekomcloud_sdrs_protectiongroup_v1.group_1.id
  server_id            = opentelekomcloud_ecs_instance_v1.instance_1.id
  delete_target_server = true

  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
}
`, env.OS_AVAILABILITY_ZONE, env.OS_VPC_ID, env.OS_IMAGE_ID, env.OS_NETWORK_ID)

var testAccSdrsProtectedInstanceV1_update = fmt.Sprintf(`
resource "opentelekomcloud_sdrs_protectiongroup_v1" "group_1" {
  name                     = "group_1"
  source_availability_zone = "%[1]s"
  target_availability_zone = "eu-de-01"
  domain_id                = "cdba26b2-cc35-4988-a904-82b7abf20094"
  source_vpc_id            = "%[2]s"
  dr_type                  = "migration"
}

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = "%[3]s"
  flavor   = "s2.medium.1"
  vpc_id   = "%[2]s"

  nics {
    network_id = "%[4]s"
  }

  availability_zone = "%[1]s"
}

resource "opentelekomcloud_sdrs_protected_instance_v1" "instance_1" {
  name                 = "instance_update"
  description          = "some interesting description"
  group_id             = opentelekomcloud_sdrs_protectiongroup_v1.group_1.id
  server_id            = opentelekomcloud_ecs_instance_v1.instance_1.id
  delete_target_server = true

  tags = {
    muh = "value-update"
  }
}
`, env.OS_AVAILABILITY_ZONE, env.OS_VPC_ID, env.OS_IMAGE_ID, env.OS_NETWORK_ID)
