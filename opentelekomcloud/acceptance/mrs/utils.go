package acceptance

import (
	"testing"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func testAccPreCheckMrs(t *testing.T) {
	common.TestAccPreCheckRequiredEnvVars(t)
}
