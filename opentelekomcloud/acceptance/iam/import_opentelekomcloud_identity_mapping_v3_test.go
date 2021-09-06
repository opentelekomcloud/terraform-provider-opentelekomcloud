package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccIdentityV3Mapping_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_identity_mapping_v3.mapping"
	var mappingName = fmt.Sprintf("acctest-%s", acctest.RandString(3))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheckRequiredEnvVars(t)
			common.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3MappingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3MappingBasic(mappingName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
