package acceptance

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccImagesImageAccessAcceptV2ImportBasic(t *testing.T) {
	acceptResourceName := "opentelekomcloud_images_image_access_accept_v2.accept_1"

	privateImageID := os.Getenv("OS_PRIVATE_IMAGE_ID")
	shareProjectID := os.Getenv("OS_PROJECT_ID_2")
	if privateImageID == "" || shareProjectID == "" {
		t.Skip("OS_PRIVATE_IMAGE_ID or OS_PROJECT_ID_2 are empty, but test requires")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckImagesImageAccessAcceptV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageAccessAcceptV2Basic(privateImageID, shareProjectID),
			},
			{
				ResourceName:      acceptResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
