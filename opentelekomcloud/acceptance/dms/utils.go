package acceptance

import (
	"testing"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func testAccPreCheckDms(t *testing.T) {
	common.TestAccPreCheckRequiredEnvVars(t)
}
