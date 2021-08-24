package acceptance

import (
	"testing"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func testAccPreCheckDms(t *testing.T) {
	common.TestAccPreCheckRequiredEnvVars(t)

	if env.OsDmsEnvironment == "" {
		t.Skip("This environment does not support DMS tests")
	}
}
