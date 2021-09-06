package acceptance

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccKmsV1GrantImportBasic(t *testing.T) {
	resourceName := "opentelekomcloud_kms_grant_v1.grant_1"

	granteePrincipal := os.Getenv("OS_USER_ID")
	if granteePrincipal == "" {
		t.Skip("OS_USER_ID must be set for acceptance test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccKMSPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckKmsV1GrantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsGrantV1Basic(granteePrincipal),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
