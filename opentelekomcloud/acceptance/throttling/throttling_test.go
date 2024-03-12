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

var mergedWafdConfigs = fmt.Sprintf("%s\n%s\n", src.WafdVars, src.WafdMain)

func TestThrottlingWafDedicatedConfiguration(t *testing.T) {
	if os.Getenv("OS_THROTTLING_WAFD") == "" {
		t.Skip("OS_THROTTLING_WAFD is not set; skipping OpenTelekomCloud WAFD THROTTLING test.")
	}
	err := os.Setenv("OS_MAX_BACKOFF_RETRIES", "30")
	if err != nil {
		return
	}
	err = os.Setenv("OS_BACKOFF_RETRY_TIMEOUT", "60")
	if err != nil {
		return
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             mergedWafdConfigs,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: mergedWafdConfigs,
			},
		},
	})
}
