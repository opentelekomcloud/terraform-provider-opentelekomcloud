package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	quick_search "github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/quick-search"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getQuickSearchCriteriaResourceFunc(config *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.LtsV10Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS v1.0 client: %s", err)
	}

	groupID := state.Primary.Attributes["log_group_id"]
	streamID := state.Primary.Attributes["log_stream_id"]
	requestResp, err := quick_search.ListCriterias(
		client,
		groupID,
		streamID,
		quick_search.ListOpts{
			SearchType: state.Primary.Attributes["type"],
		},
	)
	if err != nil {
		return nil, err
	}
	if len(requestResp) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	var searchResult quick_search.SearchCriteria
	for _, sq := range requestResp {
		if sq.ID == state.Primary.ID {
			searchResult = sq
			break
		}
	}
	if searchResult.ID == "" {
		return nil, golangsdk.ErrDefault404{}
	}
	return searchResult, nil
}

func TestAccSearchCriteria_basic(t *testing.T) {
	var (
		sq    quick_search.SearchCriteria
		rName = "opentelekomcloud_lts_quick_search_criteria_v1.log"
		name  = fmt.Sprintf("lts_qs%s", acctest.RandString(3))
		rc    = common.InitResourceCheck(rName, &sq, getQuickSearchCriteriaResourceFunc)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rc.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testquicksearchcriteriaBasic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "criteria", "context:test"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "ORIGINALLOG"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: resourceQuickSearchCriteriaImportState(rName),
			},
		},
	})
}

func resourceQuickSearchCriteriaImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		qsID := rs.Primary.ID
		groupID := rs.Primary.Attributes["log_group_id"]
		streamID := rs.Primary.Attributes["log_stream_id"]

		return fmt.Sprintf("%s/%s/%s", groupID, streamID, qsID), nil
	}
}

func testquicksearchcriteriaBasic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_group_v2" "group" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "opentelekomcloud_lts_stream_v2" "stream" {
  group_id    = opentelekomcloud_lts_group_v2.group.id
  stream_name = "%[1]s"
}

resource "opentelekomcloud_lts_quick_search_criteria_v1" "log" {
  log_group_id  = opentelekomcloud_lts_group_v2.group.id
  log_stream_id = opentelekomcloud_lts_stream_v2.stream.id

  criteria = "context:test"
  name     = "%[1]s"
  type     = "ORIGINALLOG"
}
`, name)
}
