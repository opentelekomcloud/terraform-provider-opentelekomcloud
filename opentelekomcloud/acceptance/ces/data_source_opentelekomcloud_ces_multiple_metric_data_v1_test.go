package acceptance

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const multipleMetricDatadataSourceName = "data.opentelekomcloud_ces_multiple_metric_data_v1.metric_data_1"

func TestAccCESMultipleMetricDataV1DataSource_basic(t *testing.T) {
	if os.Getenv("RUN_CES_EVENTS") == "" {
		t.Skip("CES tests are not suitable for CI runner unless specifically flagged")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCESMultipleMetricDataBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCESMultipleMetricV1DataSourceID(multipleMetricDatadataSourceName),
					resource.TestCheckResourceAttr(multipleMetricDatadataSourceName, "metrics.0.unit", "%"),
				),
			},
		},
	})
}

func testAccCheckCESMultipleMetricV1DataSourceID(n string) resource.TestCheckFunc {
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

var testDataSourceCESMultipleMetricDataBasic = fmt.Sprintf(`
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

data "opentelekomcloud_ces_multiple_metric_data_v1" "metric_data_1" {
  depends_on = [opentelekomcloud_ces_metric_data_v1.metric_1]
  metrics {
    namespace   = "TEST.TF_ACC"
    metric_name = "cpu_util"
    dimensions {
      name  = "instance_id"
      value = "72d1377e-09e4-47bd-8ea4-71a815d4815d"
    }
  }
  from   = %d
  to     = %d
  period = "1"
  filter = "average"
}
`, time.Now().UnixMilli(), time.Now().UnixMilli()-300000, time.Now().UnixMilli()+10000)
