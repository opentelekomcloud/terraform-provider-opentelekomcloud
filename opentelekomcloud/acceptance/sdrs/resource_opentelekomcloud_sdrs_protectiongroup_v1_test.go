package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/sdrs/v1/protectiongroups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccSdrsProtectiongroupV1_basic(t *testing.T) {
	var group protectiongroups.Group

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckSdrsProtectiongroupV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSdrsProtectiongroupV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSdrsProtectiongroupV1Exists("opentelekomcloud_sdrs_protectiongroup_v1.group_1", &group),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_sdrs_protectiongroup_v1.group_1", "name", "group_1"),
				),
			},
			{
				Config: testAccSdrsProtectiongroupV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSdrsProtectiongroupV1Exists("opentelekomcloud_sdrs_protectiongroup_v1.group_1", &group),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_sdrs_protectiongroup_v1.group_1", "name", "group_updated"),
				),
			},
		},
	})
}

func testAccCheckSdrsProtectiongroupV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	sdrsClient, err := config.SdrsV1Client(env.OsRegionName)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud SDRS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_sdrs_protectiongroup_v1" {
			continue
		}

		_, err := protectiongroups.Get(sdrsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("SDRS protectiongroup still exists")
		}
	}

	return nil
}

func testAccCheckSdrsProtectiongroupV1Exists(n string, group *protectiongroups.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		sdrsClient, err := config.SdrsV1Client(env.OsRegionName)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud SDRS client: %s", err)
		}

		found, err := protectiongroups.Get(sdrsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("SDRS protectiongroup not found")
		}

		*group = *found

		return nil
	}
}

var testAccSdrsProtectiongroupV1_basic = fmt.Sprintf(`
resource "opentelekomcloud_sdrs_protectiongroup_v1" "group_1" {
	name = "group_1"
	description = "test description"
	source_availability_zone = "eu-de-02"
	target_availability_zone = "eu-de-01"
	domain_id = "cdba26b2-cc35-4988-a904-82b7abf20094"
	source_vpc_id = "%s"
	dr_type = "migration"
}`, env.OsRouterID)

var testAccSdrsProtectiongroupV1_update = fmt.Sprintf(`
resource "opentelekomcloud_sdrs_protectiongroup_v1" "group_1" {
	name = "group_updated"
	description = "test description"
	source_availability_zone = "eu-de-02"
	target_availability_zone = "eu-de-01"
	domain_id = "cdba26b2-cc35-4988-a904-82b7abf20094"
	source_vpc_id = "%s"
	dr_type = "migration"
}`, env.OsRouterID)
