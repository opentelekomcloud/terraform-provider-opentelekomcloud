package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
)

const dataSourceVPCFlowLogsName = "data.opentelekomcloud_vpc_flow_logs_v1.flow_logs"

func TestAccVpcFlowLogsV1DataSource_basic(t *testing.T) {
	t.Parallel()
	quotas.BookOne(t, quotas.Router)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcFlowLogsV1DataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceVPCFlowLogsName, "flow_logs.#", "1"),
					resource.TestCheckResourceAttrPair(
						dataSourceVPCFlowLogsName, "flow_logs.0.id",
						resourceVPCFlowLogName, "id",
					),
					resource.TestCheckResourceAttr(dataSourceVPCFlowLogsName, "flow_logs.0.name", flName),
					resource.TestCheckResourceAttr(dataSourceVPCFlowLogsName, "flow_logs.0.enabled", "false"),
					resource.TestCheckResourceAttrSet(dataSourceVPCFlowLogsName, "flow_logs.0.created_at"),
				),
			},
		},
	})
}

var testAccVpcFlowLogsV1DataSource = fmt.Sprintf(`
%s

data "opentelekomcloud_vpc_flow_logs_v1" "flow_logs" {
  id            = opentelekomcloud_vpc_flow_log_v1.flow_logl.id
  name          = opentelekomcloud_vpc_flow_log_v1.flow_logl.name
  tenant_id     = opentelekomcloud_vpc_flow_log_v1.flow_logl.tenant_id
  description   = opentelekomcloud_vpc_flow_log_v1.flow_logl.description
  resource_type = opentelekomcloud_vpc_flow_log_v1.flow_logl.resource_type
  resource_id   = opentelekomcloud_vpc_flow_log_v1.flow_logl.resource_id
  traffic_type  = opentelekomcloud_vpc_flow_log_v1.flow_logl.traffic_type
  log_group_id  = opentelekomcloud_vpc_flow_log_v1.flow_logl.log_group_id
  log_topic_id  = opentelekomcloud_vpc_flow_log_v1.flow_logl.log_topic_id
  status        = opentelekomcloud_vpc_flow_log_v1.flow_logl.status
  limit         = 1
}
`, testAccVpcFlowLogV1Basic)
