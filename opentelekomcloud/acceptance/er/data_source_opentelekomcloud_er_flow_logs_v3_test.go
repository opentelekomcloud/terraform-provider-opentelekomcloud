package er

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDatasourceErFlowLogsV3_basic(t *testing.T) {
	var (
		name           = fmt.Sprintf("er-data-fl%s", acctest.RandString(5))
		bgpAsNum       = acctest.RandIntRange(64512, 65534)
		dataSourceName = "data.opentelekomcloud_er_flow_logs_v3.test"
		dc             = common.InitDataSourceCheck(dataSourceName)

		byResourceType   = "data.opentelekomcloud_er_flow_logs_v3.filter_by_resource_type"
		dcByResourceType = common.InitDataSourceCheck(byResourceType)

		byResourceId   = "data.opentelekomcloud_er_flow_logs_v3.filter_by_resource_id"
		dcByResourceId = common.InitDataSourceCheck(byResourceId)

		byFlowLogId   = "data.opentelekomcloud_er_flow_logs_v3.filter_by_flow_log_id"
		dcByFlowLogId = common.InitDataSourceCheck(byFlowLogId)

		byName   = "data.opentelekomcloud_er_flow_logs_v3.filter_by_name"
		dcByName = common.InitDataSourceCheck(byName)

		byStatus   = "data.opentelekomcloud_er_flow_logs_v3.filter_by_status"
		dcByStatus = common.InitDataSourceCheck(byStatus)

		byEnabled   = "data.opentelekomcloud_er_flow_logs_v3.filter_by_enabled"
		dcByEnabled = common.InitDataSourceCheck(byEnabled)

		byLogGroupId   = "data.opentelekomcloud_er_flow_logs_v3.filter_by_log_group_id"
		dcByLogGroupId = common.InitDataSourceCheck(byLogGroupId)

		byLogStreamId   = "data.opentelekomcloud_er_flow_logs_v3.filter_by_log_stream_id"
		dcByLogStreamId = common.InitDataSourceCheck(byLogStreamId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceErFlowLogsV3_basic(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					dcByResourceType.CheckResourceExists(),
					resource.TestCheckOutput("resource_type_filter_is_useful", "true"),

					dcByResourceId.CheckResourceExists(),
					resource.TestCheckOutput("resource_id_filter_is_useful", "true"),

					dcByFlowLogId.CheckResourceExists(),
					resource.TestCheckOutput("flow_log_id_filter_is_useful", "true"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),

					dcByEnabled.CheckResourceExists(),
					resource.TestCheckOutput("enabled_filter_is_useful", "true"),

					dcByLogGroupId.CheckResourceExists(),
					resource.TestCheckOutput("log_group_id_filter_is_useful", "true"),

					dcByLogStreamId.CheckResourceExists(),
					resource.TestCheckOutput("log_stream_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceErFlowLogsV3_basic(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%s

data "opentelekomcloud_er_flow_logs_v3" "test" {
  depends_on  = [opentelekomcloud_er_flow_log_v3.test]
  instance_id = opentelekomcloud_er_instance_v3.test.id
}

locals {
  resource_type = data.opentelekomcloud_er_flow_logs_v3.test.flow_logs[0].resource_type
}

data "opentelekomcloud_er_flow_logs_v3" "filter_by_resource_type" {
  instance_id   = opentelekomcloud_er_instance_v3.test.id
  resource_type = local.resource_type
}

locals {
  resource_type_filter_result = [
    for v in data.opentelekomcloud_er_flow_logs_v3.filter_by_resource_type.flow_logs[*].resource_type :
    v == local.resource_type
  ]
}

output "resource_type_filter_is_useful" {
  value = alltrue(local.resource_type_filter_result) && length(local.resource_type_filter_result) > 0
}

locals {
  resource_id = data.opentelekomcloud_er_flow_logs_v3.test.flow_logs[0].resource_id
}

data "opentelekomcloud_er_flow_logs_v3" "filter_by_resource_id" {
  instance_id = opentelekomcloud_er_instance_v3.test.id
  resource_id = local.resource_id
}

locals {
  resource_id_filter_result = [
    for v in data.opentelekomcloud_er_flow_logs_v3.filter_by_resource_id.flow_logs[*].resource_id : v == local.resource_id
  ]
}

output "resource_id_filter_is_useful" {
  value = alltrue(local.resource_id_filter_result) && length(local.resource_id_filter_result) > 0
}

locals {
  flow_log_id = data.opentelekomcloud_er_flow_logs_v3.test.flow_logs[0].id
}

data "opentelekomcloud_er_flow_logs_v3" "filter_by_flow_log_id" {
  instance_id = opentelekomcloud_er_instance_v3.test.id
  flow_log_id = local.flow_log_id
}

locals {
  flow_log_id_filter_result = [
    for v in data.opentelekomcloud_er_flow_logs_v3.filter_by_flow_log_id.flow_logs[*].id : v == local.flow_log_id
  ]
}

output "flow_log_id_filter_is_useful" {
  value = alltrue(local.flow_log_id_filter_result) && length(local.flow_log_id_filter_result) > 0
}

locals {
  name = data.opentelekomcloud_er_flow_logs_v3.test.flow_logs[0].name
}

data "opentelekomcloud_er_flow_logs_v3" "filter_by_name" {
  instance_id = opentelekomcloud_er_instance_v3.test.id
  name        = local.name
}

locals {
  name_filter_result = [
    for v in data.opentelekomcloud_er_flow_logs_v3.filter_by_name.flow_logs[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

locals {
  status = data.opentelekomcloud_er_flow_logs_v3.test.flow_logs[0].status
}

data "opentelekomcloud_er_flow_logs_v3" "filter_by_status" {
  instance_id = opentelekomcloud_er_instance_v3.test.id
  status      = local.status
}

locals {
  status_filter_result = [
    for v in data.opentelekomcloud_er_flow_logs_v3.filter_by_status.flow_logs[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

locals {
  log_group_id = data.opentelekomcloud_er_flow_logs_v3.test.flow_logs[0].log_group_id
}

data "opentelekomcloud_er_flow_logs_v3" "filter_by_log_group_id" {
  instance_id  = opentelekomcloud_er_instance_v3.test.id
  log_group_id = local.log_group_id
}

locals {
  log_group_id_filter_result = [
    for v in data.opentelekomcloud_er_flow_logs_v3.filter_by_log_group_id.flow_logs[*].log_group_id : v == local.log_group_id
  ]
}

output "log_group_id_filter_is_useful" {
  value = alltrue(local.log_group_id_filter_result) && length(local.log_group_id_filter_result) > 0
}

locals {
  log_stream_id = data.opentelekomcloud_er_flow_logs_v3.test.flow_logs[0].log_stream_id
}

data "opentelekomcloud_er_flow_logs_v3" "filter_by_log_stream_id" {
  instance_id   = opentelekomcloud_er_instance_v3.test.id
  log_stream_id = local.log_stream_id
}

locals {
  log_stream_id_filter_result = [
    for v in data.opentelekomcloud_er_flow_logs_v3.filter_by_log_stream_id.flow_logs[*].log_stream_id :
    v == local.log_stream_id
  ]
}

output "log_stream_id_filter_is_useful" {
  value = alltrue(local.log_stream_id_filter_result) && length(local.log_stream_id_filter_result) > 0
}

locals {
  enabled = data.opentelekomcloud_er_flow_logs_v3.test.flow_logs[0].enabled
}

data "opentelekomcloud_er_flow_logs_v3" "filter_by_enabled" {
  instance_id = opentelekomcloud_er_instance_v3.test.id
  enabled     = local.enabled
}

locals {
  enabled_filter_result = [
    for v in data.opentelekomcloud_er_flow_logs_v3.filter_by_enabled.flow_logs[*].enabled : v == local.enabled
  ]
}

output "enabled_filter_is_useful" {
  value = alltrue(local.enabled_filter_result) && length(local.enabled_filter_result) > 0
}
`, testFlowLog_basic(testaccFlowLog_base(name), name))
}
