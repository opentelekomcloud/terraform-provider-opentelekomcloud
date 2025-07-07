package acceptance

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const metricDatadataSourceName = "data.opentelekomcloud_ces_metric_data_v1.metric_data_1"

func TestAccCESMetricDataV1DataSource_basic(t *testing.T) {
	if os.Getenv("RUN_CES_EVENTS") == "" {
		t.Skip("CES tests are not suitable for CI runner unless specifically flagged")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCESMetricDataBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCESMetricV1DataSourceID(metricDatadataSourceName),
					resource.TestCheckResourceAttr(metricDatadataSourceName, "datapoints.0.unit", "%"),
				),
			},
		},
	})
}

func testAccCheckCESMetricV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find CES Metric Data data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CES metric data data source ID not set")
		}

		return nil
	}
}

var testDataSourceCESMetricDataBasic = fmt.Sprintf(`
resource "opentelekomcloud_ces_metric_data_v1" "metric_1" {
  metric {
    namespace   = "TEST.TF_ACC"
    metric_name = "cpu_util"
    dimensions {
      name  = "instance_id"
      value = "72d1377e-09e4-47bd-8ea4-71a815d4815d"
    }
  }
  ttl          = 172800
  collect_time = %d
  value        = 0.09
  unit         = "%%"
  type         = "float"
}

data "opentelekomcloud_ces_metric_data_v1" "metric_data_1" {
  depends_on  = [opentelekomcloud_ces_metric_data_v1.metric_1]
  namespace   = "TEST.TF_ACC"
  metric_name = "cpu_util"
  from        = "%s"
  to          = "%s"
  period      = 1
  filter      = "average"
  dim0        = "instance_id,72d1377e-09e4-47bd-8ea4-71a815d4815d"
}
`, time.Now().UnixMilli(), strconv.FormatInt(time.Now().UnixMilli()-300000, 10), strconv.FormatInt(time.Now().UnixMilli()+10000, 10))
