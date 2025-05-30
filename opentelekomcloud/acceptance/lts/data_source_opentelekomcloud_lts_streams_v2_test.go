package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDataSourceStreams_basic(t *testing.T) {
	var (
		dataSource = "data.opentelekomcloud_lts_streams_v2.test"
		rName      = fmt.Sprintf("lts_streams%s", acctest.RandString(5))
		dc         = common.InitDataSourceCheck(dataSource)

		byName   = "data.opentelekomcloud_lts_streams_v2.filter_by_name"
		dcByName = common.InitDataSourceCheck(byName)

		byLogGroupName   = "data.opentelekomcloud_lts_streams_v2.filter_by_log_group_name"
		dcByLogGroupName = common.InitDataSourceCheck(byLogGroupName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceStreams_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "streams.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(dataSource, "streams.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcByName.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(byName, "streams.0.id", "opentelekomcloud_lts_stream_v2.test", "id"),
					resource.TestCheckResourceAttr(byName, "streams.0.name", rName),
					resource.TestCheckResourceAttr(byName, "streams.0.ttl_in_days", "60"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByLogGroupName.CheckResourceExists(),
					resource.TestCheckOutput("is_log_group_name_filter_useful", "true"),
				),
			},
			{
				Config:      testDataSourceStreams_logStreamNotFoundError(),
				ExpectError: regexp.MustCompile("The log stream does not exist"),
			},
			{
				Config:      testDataSourceStreams_logGroupNotFoundError(),
				ExpectError: regexp.MustCompile("The log group does not exist"),
			},
		},
	})
}

func testDataSourceStreams_basic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_group_v2" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "opentelekomcloud_lts_stream_v2" "test" {
  group_id    = opentelekomcloud_lts_group_v2.test.id
  stream_name = opentelekomcloud_lts_group_v2.test.group_name
  ttl_in_days = 60
}

data "opentelekomcloud_lts_streams_v2" "test" {
  depends_on = [
    opentelekomcloud_lts_stream_v2.test
  ]
}

locals {
  stream_name = opentelekomcloud_lts_stream_v2.test.stream_name
}

data "opentelekomcloud_lts_streams_v2" "filter_by_name" {
  depends_on = [
    opentelekomcloud_lts_stream_v2.test
  ]

  name = local.stream_name
}

locals {
  name_filter_result = [
    for v in data.opentelekomcloud_lts_streams_v2.filter_by_name.streams[*].name : v == local.stream_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

locals {
  log_group_name = opentelekomcloud_lts_group_v2.test.group_name
}

data "opentelekomcloud_lts_streams_v2" "filter_by_log_group_name" {
  depends_on = [
    opentelekomcloud_lts_stream_v2.test
  ]

  log_group_name = local.log_group_name
}

locals {
  stream_ids = data.opentelekomcloud_lts_streams_v2.filter_by_log_group_name.streams[*].id
}

output "is_log_group_name_filter_useful" {
  value = length(local.stream_ids) == 1 && local.stream_ids[0] == opentelekomcloud_lts_stream_v2.test.id
}
`, name)
}

func testDataSourceStreams_logStreamNotFoundError() string {
	return `
data "opentelekomcloud_lts_streams_v2" "not_found_log_stream" {
  name = "not_found_log_stream_name"
}
`
}

func testDataSourceStreams_logGroupNotFoundError() string {
	return `
data "opentelekomcloud_lts_streams_v2" "not_found_log_group" {
  log_group_name = "not_found_log_group_name"
}
`
}
