package acceptance

import (
	"testing"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccIdentityV3AgencyPreCheck(t *testing.T) {
	if env.OS_TENANT_NAME == "" {
		t.Skip("OS_TENANT_NAME must be set for acceptance tests")
	}
}
