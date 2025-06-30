package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const metricsdataSourceName = "data.opentelekomcloud_ces_metrics_v1.metrics_1"

func TestAccCESMetricsV1DataSource_basic(t *testing.T) {
	if os.Getenv("RUN_CES_EVENTS") == "" {
		t.Skip("CES tests are not suitable for CI runner unless specifically flagged")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testCESMetricsBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCESMetricsV1DataSourceID(metricsdataSourceName),
					resource.TestCheckResourceAttr(metricsdataSourceName, "metrics.0.namespace", "SYS.ECS"),
					resource.TestCheckResourceAttr(metricsdataSourceName, "metrics.0.metric_name", "cpu_util"),
				),
			},
		},
	})
}

func testAccCheckCESMetricsV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find CES Metrics data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CES metrics data source ID not set")
		}

		return nil
	}
}

var testCESMetricsBasic = `
data "opentelekomcloud_ces_metrics_v1" "metrics_1" {
  namespace   = "SYS.ECS"
  metric_name = "cpu_util"
}
`
