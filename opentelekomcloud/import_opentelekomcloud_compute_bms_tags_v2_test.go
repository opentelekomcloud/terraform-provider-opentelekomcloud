package opentelekomcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOTCBMSTagsV2_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_compute_bms_tags_v2.tags_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckRequiredEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCBMSTagsV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccBMSTagsV2_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
