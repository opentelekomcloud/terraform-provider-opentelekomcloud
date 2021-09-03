package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDedicatedHostV1_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDeHV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeHV1Basic,
			},
			{
				ResourceName:      resourceHostName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
