package acceptance

import (
	"testing"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func testAccPreCheckS3(t *testing.T) {
	common.TestAccPreCheckRequiredEnvVars(t)

	if env.OS_ACCESS_KEY == "" || env.OS_SECRET_KEY == "" {
		t.Skip("OS_ACCESS_KEY and OS_SECRET_KEY must be set for OBS/S3 acceptance tests")
	}
}
