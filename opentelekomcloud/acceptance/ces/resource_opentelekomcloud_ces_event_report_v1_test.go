package acceptance

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const resourceCesEventReportName = "opentelekomcloud_ces_event_report_v1.event_report_1"

func TestResourceCesEventReportV1(t *testing.T) {
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
				Config: testCESEventReportBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceCesEventReportName, "event_id"),
				),
			},
		},
	})
}

var testCESEventReportBasic = fmt.Sprintf(`
resource "opentelekomcloud_ces_event_report_v1" "event_report_1" {
  event_name   = "Test_Acc_tf_event"
  event_source = "SYS.ECS"
  time         = %d
  detail {
    content     = "This is a test event run by TF tests"
    event_state = "normal"
    event_level = "Info"
  }
}
`, time.Now().Unix()*1000)
