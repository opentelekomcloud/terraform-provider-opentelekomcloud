package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	ac "github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/access-config"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCrossAccountAccess_basic(t *testing.T) {
	agencyProjectId := os.Getenv("OS_LTS_AGENCY_PROJECT_ID")
	agencyDomainName := os.Getenv("OS_LTS_AGENCY_DOMAIN_NAME")
	agencyName := os.Getenv("OS_LTS_AGENCY_NAME")
	agencyStreamName := os.Getenv("OS_LTS_AGENCY_STREAM_NAME")
	agencyStreamId := os.Getenv("OS_LTS_AGENCY_STREAM_ID")
	agencyGroupName := os.Getenv("OS_LTS_AGENCY_GROUP_NAME")
	agencyGroupId := os.Getenv("OS_LTS_AGENCY_GROUP_ID")

	if agencyProjectId == "" || agencyDomainName == "" || agencyName == "" || agencyStreamName == "" ||
		agencyStreamId == "" || agencyGroupName == "" || agencyGroupId == "" {
		t.Skip("The delegator account config of OS_LTS_AGENCY_STREAM_NAME, OS_LTS_AGENCY_STREAM_ID, " +
			"OS_LTS_AGENCY_GROUP_NAME, OS_LTS_AGENCY_GROUP_ID, OS_LTS_AGENCY_PROJECT_ID, " +
			"OS_LTS_AGENCY_DOMAIN_NAME and OS_LTS_AGENCY_NAME must be set for the acceptance test")
	}
	var (
		access ac.AccessConfigInfo
		rName  = "opentelekomcloud_lts_cross_account_access_v2.test"
		name   = fmt.Sprintf("lts_cross_access%s", acctest.RandString(3))
		rc     = common.InitResourceCheck(rName, &access, getHostAccessConfigResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCrossAccuntAccessBasic(name, agencyProjectId, agencyDomainName, agencyName, agencyStreamName, agencyStreamId, agencyGroupName, agencyGroupId),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "access_config_type"),
				),
			},
			{
				Config: testCrossAccuntAccessBasicUpdate(name, agencyProjectId, agencyDomainName, agencyName, agencyStreamName, agencyStreamId, agencyGroupName, agencyGroupId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar1"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1"),
				),
			},
		},
	})
}

func testCrossAccuntAccessBasic(name, agencyProjectId, agencyDomainName, agencyName, agencyStreamName, agencyStreamId, agencyGroupName, agencyGroupId string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_group_v2" "group" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "opentelekomcloud_lts_stream_v2" "stream" {
  group_id    = opentelekomcloud_lts_group_v2.group.id
  stream_name = "%[1]s"
}

resource "opentelekomcloud_lts_cross_account_access_v2" "test" {
  name               = "%[1]s"
  agency_project_id  = "%s"
  agency_domain_name = "%s"
  agency_name        = "%s"

  log_agency_stream_name = "%s"
  log_agency_stream_id   = "%s"
  log_agency_group_name  = "%s"
  log_agency_group_id    = "%s"

  log_stream_name = opentelekomcloud_lts_stream_v2.stream.stream_name
  log_stream_id   = opentelekomcloud_lts_stream_v2.stream.id
  log_group_name  = opentelekomcloud_lts_group_v2.group.group_name
  log_group_id    = opentelekomcloud_lts_group_v2.group.id

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, agencyProjectId, agencyDomainName, agencyName, agencyStreamName, agencyStreamId, agencyGroupName, agencyGroupId)
}

func testCrossAccuntAccessBasicUpdate(name, agencyProjectId, agencyDomainName, agencyName, agencyStreamName, agencyStreamId, agencyGroupName, agencyGroupId string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_group_v2" "group" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "opentelekomcloud_lts_stream_v2" "stream" {
  group_id    = opentelekomcloud_lts_group_v2.group.id
  stream_name = "%[1]s"
}

resource "opentelekomcloud_lts_cross_account_access_v2" "test" {
  name               = "%[1]s"
  agency_project_id  = "%s"
  agency_domain_name = "%s"
  agency_name        = "%s"

  log_agency_stream_name = "%s"
  log_agency_stream_id   = "%s"
  log_agency_group_name  = "%s"
  log_agency_group_id    = "%s"

  log_stream_name = opentelekomcloud_lts_stream_v2.stream.stream_name
  log_stream_id   = opentelekomcloud_lts_stream_v2.stream.id
  log_group_name  = opentelekomcloud_lts_group_v2.group.group_name
  log_group_id    = opentelekomcloud_lts_group_v2.group.id

  tags = {
    foo  = "bar1"
    key1 = "value1"
  }
}
`, name, agencyProjectId, agencyDomainName, agencyName, agencyStreamName, agencyStreamId, agencyGroupName, agencyGroupId)
}
