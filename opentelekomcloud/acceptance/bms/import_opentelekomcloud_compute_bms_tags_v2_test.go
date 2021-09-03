package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccBMSTagsV2_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckBMSTagsV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSTagsV2Basic,
			},
			{
				ResourceName:      resourceTagsName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
