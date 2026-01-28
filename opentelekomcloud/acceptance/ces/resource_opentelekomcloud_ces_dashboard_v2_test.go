package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v2/dashboards"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceDashboardV2Name = "opentelekomcloud_ces_dashboard_v2.test"

func getDashboardV2Func(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CesV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud CES v2 client: %s", err)
	}
	listResp, err := dashboards.List(c, dashboards.ListOpts{
		DashboardId: state.Primary.ID,
	})
	if err != nil {
		return nil, err
	}
	if listResp == nil || len(listResp) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}
	return listResp[0], nil
}

func TestAccCESDashboardV2_basic(t *testing.T) {
	var (
		dashboard dashboards.Dashboard
		rName     = resourceDashboardV2Name
		name      = acctest.RandomWithPrefix("ces-dash")
	)

	rc := common.InitResourceCheck(
		rName,
		&dashboard,
		getDashboardV2Func,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDashboardV2_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "row_widget_num", "1"),
					resource.TestCheckResourceAttr(rName, "is_favorite", "true"),
					resource.TestCheckResourceAttrSet(rName, "creator_name"),
					resource.TestMatchResourceAttr(rName,
						"created_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(Z|([+-]\d{2}:\d{2}))`)),
				),
			},
			{
				Config: testDashboardV2_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "row_widget_num", "2"),
					resource.TestCheckResourceAttr(rName, "is_favorite", "false"),
				),
			},
			{
				Config: testDashboardV2_zeroWidgetNum(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "row_widget_num", "0"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"dashboard_id",
				},
			},
		},
	})
}

func TestAccCESDashboardV2_copy(t *testing.T) {
	var (
		dashboard dashboards.Dashboard
		rName     = resourceDashboardV2Name
		name      = acctest.RandomWithPrefix("ces-dash")
	)

	rc := common.InitResourceCheck(
		rName,
		&dashboard,
		getDashboardV2Func,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDashboardV2_copy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "row_widget_num", "1"),
					resource.TestCheckResourceAttr(rName, "is_favorite", "false"),
					resource.TestCheckResourceAttrSet(rName, "creator_name"),
					resource.TestCheckResourceAttrPair(rName, "dashboard_id", "opentelekomcloud_ces_dashboard_v2.base", "id"),
					resource.TestMatchResourceAttr(rName,
						"created_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(Z|([+-]\d{2}:\d{2}))`)),
				),
			},
			{
				Config: testDashboardV2_copy_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "row_widget_num", "3"),
					resource.TestCheckResourceAttr(rName, "is_favorite", "true"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"dashboard_id",
				},
			},
		},
	})
}

func testDashboardV2_basic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_ces_dashboard_v2" "test" {
  name           = "%s"
  row_widget_num = 1
  is_favorite    = true
}
`, name)
}

func testDashboardV2_update(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_ces_dashboard_v2" "test" {
  name           = "%s-update"
  row_widget_num = 2
  is_favorite    = false
}
`, name)
}

func testDashboardV2_zeroWidgetNum(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_ces_dashboard_v2" "test" {
  name           = "%s-update"
  row_widget_num = 0
}
`, name)
}

func testDashboardV2_copy_basic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_ces_dashboard_v2" "base" {
  name           = "%s-base"
  row_widget_num = 2
}

resource "opentelekomcloud_ces_dashboard_v2" "test" {
  name           = "%s"
  dashboard_id   = opentelekomcloud_ces_dashboard_v2.base.id
  row_widget_num = 1
}
`, name, name)
}

func testDashboardV2_copy_update(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_ces_dashboard_v2" "base" {
  name           = "%s-base"
  row_widget_num = 2
}

resource "opentelekomcloud_ces_dashboard_v2" "test" {
  name           = "%s"
  dashboard_id   = opentelekomcloud_ces_dashboard_v2.base.id
  is_favorite    = true
  row_widget_num = 3
}
`, name, name)
}
