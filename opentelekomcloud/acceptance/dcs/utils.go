package acceptance

import (
	"testing"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func testAccPreCheckDcs(t *testing.T) {
	common.TestAccPreCheckRequiredEnvVars(t)

	if env.OS_DCS_ENVIRONMENT == "" {
		t.Skip("This environment does not support DCS tests")
	}
}
