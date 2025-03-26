package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/groups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getLtsGroupResourceFunc(config *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.LtsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS v2 client: %s", err)
	}

	requestResp, err := groups.List(client)
	if err != nil {
		return nil, err
	}
	if len(requestResp) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	var groupResult *groups.LogGroup
	for _, group := range requestResp {
		if group.LogGroupId == state.Primary.ID {
			groupResult = &group
		}
	}
	if groupResult == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return groupResult, nil
}

func TestAccLogTankGroupV2_basic(t *testing.T) {
	var (
		group        groups.LogGroup
		resourceName = "opentelekomcloud_logtank_group_v2.group"
		rName        = fmt.Sprintf("lts_group%s", acctest.RandString(3))
		rc           = common.InitResourceCheck(resourceName, &group, getLtsGroupResourceFunc)
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLogTankGroupV2_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "group_name", rName),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "7"),
				),
			},
			{
				Config: testAccLogTankGroupV2_updated(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "group_name", rName),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "6"),
				),
			},
		},
	})
}

func testAccLogTankGroupV2_basic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_logtank_group_v2" "group" {
  group_name  = "%s"
  ttl_in_days = 7
}
`, name)
}

func testAccLogTankGroupV2_updated(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_logtank_group_v2" "group" {
  group_name  = "%s"
  ttl_in_days = 6
}
`, name)
}
