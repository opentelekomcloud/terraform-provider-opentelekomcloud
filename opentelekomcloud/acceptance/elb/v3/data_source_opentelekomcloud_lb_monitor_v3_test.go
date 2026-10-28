package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
)

func TestLBMonitorV3DataSource_basic(t *testing.T) {
	t.Parallel()
	quotas.BookMany(t, []*quotas.ExpectedQuota{
		{Q: quotas.LbPool, Count: 1},
		{Q: quotas.LoadBalancer, Count: 1},
	})
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      checkMonitorDestroy,
		Steps: []resource.TestStep{{
			Config: testLBMonitorV3DataSourceBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrPair("data.opentelekomcloud_lb_monitor_v3.by_id", "id", resourceMonitorName, "id"),
				resource.TestCheckResourceAttrPair("data.opentelekomcloud_lb_monitor_v3.by_name", "id", resourceMonitorName, "id"),
				resource.TestCheckResourceAttrPair("data.opentelekomcloud_lb_monitor_v3.by_name", "pool_id", resourceMonitorName, "pool_id"),
				resource.TestCheckResourceAttrPair("data.opentelekomcloud_lb_monitor_v3.by_name", "created_at", resourceMonitorName, "created_at"),
				resource.TestCheckResourceAttrPair("data.opentelekomcloud_lb_monitor_v3.by_id", "updated_at", resourceMonitorName, "updated_at"),
				resource.TestCheckResourceAttr("data.opentelekomcloud_lb_monitor_v3.by_name", "type", "HTTP"),
				resource.TestCheckResourceAttr("data.opentelekomcloud_lb_monitor_v3.by_name", "admin_state_up", "true"),
			),
		}},
	})
}

var testLBMonitorV3DataSourceBasic = testResourceMonitorBasic + `
data "opentelekomcloud_lb_monitor_v3" "by_id" {
  id = opentelekomcloud_lb_monitor_v3.monitor.id
}

data "opentelekomcloud_lb_monitor_v3" "by_name" {
  name           = opentelekomcloud_lb_monitor_v3.monitor.name
  admin_state_up = true
}
`
