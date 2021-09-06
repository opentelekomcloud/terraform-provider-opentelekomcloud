package acceptance

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccImagesImageAccessV2ImportBasic(t *testing.T) {
	accessResourceName := "opentelekomcloud_images_image_access_v2.access_1"

	privateImageID := os.Getenv("OS_PRIVATE_IMAGE_ID")
	shareProjectID := os.Getenv("OS_PROJECT_ID_2")
	if privateImageID == "" || shareProjectID == "" {
		t.Skip("OS_PRIVATE_IMAGE_ID or OS_PROJECT_ID_2 are empty, but test requires")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckImagesImageAccessV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageAccessV2Basic(privateImageID, shareProjectID),
			},
			{
				ResourceName:      accessResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
