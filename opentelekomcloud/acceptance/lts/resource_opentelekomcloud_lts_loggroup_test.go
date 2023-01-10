package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/groups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccLogTankGroupV2_basic(t *testing.T) {
	var group groups.LogGroup
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckLogTankGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLogTankGroupV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogTankGroupV2Exists(
						"opentelekomcloud_logtank_group_v2.testacc_group", &group),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_logtank_group_v2.testacc_group", "group_name", "testacc_group"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_logtank_group_v2.testacc_group", "ttl_in_days", "7"),
				),
			},
		},
	})
}

func testAccCheckLogTankGroupV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	ltsclient, err := config.LtsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud LTS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_logtank_group_v2" {
			continue
		}

		allGroups, err := groups.ListLogGroups(ltsclient)
		if err != nil {
			return fmt.Errorf("error listing lts groups: %s", err)
		}

		for _, group := range allGroups {
			if group.LogGroupId == rs.Primary.ID {
				return fmt.Errorf("log group (%s) still exists", rs.Primary.ID)
			}
		}

		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckLogTankGroupV2Exists(n string, group *groups.LogGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		ltsclient, err := config.LtsV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud LTS client: %s", err)
		}

		allGroups, err := groups.ListLogGroups(ltsclient)
		if err != nil {
			return err
		}

		for _, ltsGroup := range allGroups {
			if ltsGroup.LogGroupId == rs.Primary.ID {
				*group = ltsGroup
				break
			}
		}

		return nil
	}
}

const testAccLogTankGroupV2_basic = `
resource "opentelekomcloud_logtank_group_v2" "testacc_group" {
  group_name  = "testacc_group"
  ttl_in_days = 7
}
`
