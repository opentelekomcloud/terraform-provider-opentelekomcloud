package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/imageservice/v2/members"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/ims"
)

func TestAccImagesImageAccessAcceptV2_basic(t *testing.T) {
	var member members.Member
	acceptResourceName := "opentelekomcloud_images_image_access_accept_v2.accept_1"
	privateImageID := os.Getenv("OS_PRIVATE_IMAGE_ID")
	shareProjectID := os.Getenv("OS_PROJECT_ID_2")
	if privateImageID == "" || shareProjectID == "" {
		t.Skip("OS_PRIVATE_IMAGE_ID or OS_PROJECT_ID_2 are empty, but test requires")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckImagesImageAccessAcceptV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageAccessAcceptV2Basic(privateImageID, shareProjectID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageAccessV2Exists(acceptResourceName, &member),
					resource.TestCheckResourceAttrPtr(acceptResourceName, "status", &member.Status),
					resource.TestCheckResourceAttr(acceptResourceName, "status", "accepted"),
				),
			},
			{
				Config: testAccImagesImageAccessAcceptV2Update(privateImageID, shareProjectID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageAccessV2Exists(acceptResourceName, &member),
					resource.TestCheckResourceAttrPtr(acceptResourceName, "status", &member.Status),
					resource.TestCheckResourceAttr(acceptResourceName, "status", "rejected"),
				),
			},
		},
	})
}

func testAccCheckImagesImageAccessAcceptV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ImageV2Client(env.OsRegionName)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud IMSv2: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_images_image_access_accept_v2" {
			continue
		}

		imageID, memberID, err := ims.ResourceImagesImageAccessV2ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = members.Get(client, imageID, memberID).Extract()
		if err == nil {
			return fmt.Errorf("image membership still exists")
		}
	}

	return nil
}

func testAccImagesImageAccessAcceptV2Basic(privateImageID, projectToShare string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_images_image_access_v2" "access_1" {
  image_id  = "%[1]s"
  member_id = "%[2]s"
}

%[3]s

resource "opentelekomcloud_images_image_access_accept_v2" "accept_1" {
  provider = "%s"

  depends_on = [opentelekomcloud_images_image_access_v2.access_1]

  image_id  = "%[1]s"
  member_id = "%[2]s"
  status    = "accepted"
}
`, privateImageID, projectToShare, common.AlternativeProviderConfig, common.AlternativeProviderAlias)
}

func testAccImagesImageAccessAcceptV2Update(privateImageID, projectToShare string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_images_image_access_v2" "access_1" {
  image_id  = "%[1]s"
  member_id = "%[2]s"
}

%[3]s

resource "opentelekomcloud_images_image_access_accept_v2" "accept_1" {
  provider = "%[4]s"

  depends_on = [opentelekomcloud_images_image_access_v2.access_1]

  image_id  = "%[1]s"
  member_id = "%[2]s"
  status    = "rejected"
}
`, privateImageID, projectToShare, common.AlternativeProviderConfig, common.AlternativeProviderAlias)
}
