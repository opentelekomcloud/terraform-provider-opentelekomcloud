package acceptance

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const resourceCesMetricDataName = "opentelekomcloud_ces_metric_data_v1.metric_1"

func TestResourceCesMetricDataV1(t *testing.T) {
	if os.Getenv("RUN_CES_EVENTS") == "" {
		t.Skip("CES tests are not suitable for CI runner unless specifically flagged")
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testCESMetricDataBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceCesMetricDataName, "ttl"),
				),
			},
		},
	})
}

var testCESMetricDataBasic = fmt.Sprintf(`
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
`, time.Now().UnixMilli())
