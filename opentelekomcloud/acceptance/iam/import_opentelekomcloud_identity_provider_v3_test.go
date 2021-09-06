package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccIdentityV3ProviderImport(t *testing.T) {
	fullName := fmt.Sprintf("%s.%s", providerResource, "provider")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3ProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3ProviderBasic,
			},
			{
				ResourceName:      fullName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
