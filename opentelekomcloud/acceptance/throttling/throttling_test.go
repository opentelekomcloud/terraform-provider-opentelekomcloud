package throttling

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/throttling/src"
)

var mergedConfigs = fmt.Sprintf("%s\n%s\n%s", src.Vars, src.Main, src.Outputs)

func TestThrottlingConfiguration(t *testing.T) {
	if os.Getenv("OS_THROTTLING") == "" {
		t.Skip("OS_THROTTLING is not set; skipping OpenTelekomCloud THROTTLING test.")
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             mergedConfigs,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: mergedConfigs,
			},
		},
	})
}
