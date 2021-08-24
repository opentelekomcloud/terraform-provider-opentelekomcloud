package acceptance

import (
	"testing"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func testAccPreCheckMrs(t *testing.T) {
	common.TestAccPreCheckRequiredEnvVars(t)

	if env.OsMrsEnvironment == "" {
		t.Skip("This environment does not support MRS tests")
	}
}
