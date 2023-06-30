package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/image/v2/members"
	ims2 "github.com/opentelekomcloud/gophertelekomcloud/openstack/ims/v2/members"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/ims"
)

func TestAccImagesImageAccessV2Basic(t *testing.T) {
	var member ims2.Member
	accessResourceName := "opentelekomcloud_images_image_access_v2.access_1"

	privateImageID := os.Getenv("OS_PRIVATE_IMAGE_ID")
	shareProjectID := os.Getenv("OS_PROJECT_ID_2")
	if privateImageID == "" || shareProjectID == "" {
		t.Skip("OS_PRIVATE_IMAGE_ID or OS_PROJECT_ID_2 are empty, but test requires")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckImagesImageAccessV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageAccessV2Basic(privateImageID, shareProjectID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageAccessV2Exists(accessResourceName, &member),
					resource.TestCheckResourceAttrPtr(accessResourceName, "status", &member.Status),
					resource.TestCheckResourceAttr(accessResourceName, "status", "pending"),
				),
			},
		},
	})
}

func testAccCheckImagesImageAccessV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ImageV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud IMSv2: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_images_image_access_v2" {
			continue
		}

		imageID, memberID, err := ims.ResourceImagesImageAccessV2ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = members.Get(client, members.MemberOpts{
			ImageId:  imageID,
			MemberId: memberID,
		})
		if err == nil {
			return fmt.Errorf("image share still exists")
		}
	}

	return nil
}

func testAccCheckImagesImageAccessV2Exists(n string, member *ims2.Member) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.ImageV2Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud IMSv2: %w", err)
		}

		imageID, memberID, err := ims.ResourceImagesImageAccessV2ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := members.Get(client, members.MemberOpts{
			ImageId:  imageID,
			MemberId: memberID,
		})
		if err != nil {
			return err
		}

		id := fmt.Sprintf("%s/%s", found.ImageId, found.MemberId)
		if id != rs.Primary.ID {
			return fmt.Errorf("image member not found")
		}

		*member = *found

		return nil
	}
}

func testAccImagesImageAccessV2Basic(privateImageID, projectToShare string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_images_image_access_v2" "access_1" {
  image_id  = "%s"
  member_id = "%s"
}
`, privateImageID, projectToShare)
}
