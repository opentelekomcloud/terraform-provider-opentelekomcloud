package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccImagesImageAccessV2ImportBasic(t *testing.T) {
	memberName := "data.opentelekomcloud_identity_auth_scope_v3.scope"
	imageName := "opentelekomcloud_images_image_v2.image_1"
	resourceName := "opentelekomcloud_images_image_access_v2.image_access_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckImagesImageAccessV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageAccessV2Basic(),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccImagesImageAccessV2ImportID(imageName, memberName),
			},
		},
	})
}

func testAccImagesImageAccessV2ImportID(imageName, memberName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		image, ok := s.RootModule().Resources[imageName]
		if !ok {
			return "", fmt.Errorf("image not found: %s", imageName)
		}

		member, ok := s.RootModule().Resources[memberName]
		if !ok {
			return "", fmt.Errorf("member not found: %s", memberName)
		}

		return fmt.Sprintf("%s/%s", image.Primary.ID, member.Primary.Attributes["project_id"]), nil
	}
}
