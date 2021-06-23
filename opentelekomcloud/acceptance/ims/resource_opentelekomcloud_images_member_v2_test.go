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
)

const resourceName = "opentelekomcloud_images_member_v2.member_1"

func TestAccImagesMemberV2Basic(t *testing.T) {
	var member members.Member

	var projectName2 = os.Getenv("OS_PROJECT_NAME_2")
	var privateImageID = os.Getenv("OS_PRIVATE_IMAGE_ID")
	if projectName2 == "" || privateImageID == "" {
		t.Skip("OS_PROJECT_NAME_2 or OS_PRIVATE_IMAGE_ID are empty")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckImagesMemberV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesMemberV2Basic(projectName2, privateImageID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesMemberV2Exists(resourceName, &member),
					resource.TestCheckResourceAttr(resourceName, "status", "pending"),
				),
			},
		},
	})
}

func testAccCheckImagesMemberV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ImageV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud IMSv2 client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_images_member_v2" {
			continue
		}

		imageID := rs.Primary.Attributes["image_id"]
		_, err := members.Get(client, imageID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("member still exists")
		}
	}

	return nil
}

func testAccCheckImagesMemberV2Exists(n string, member *members.Member) resource.TestCheckFunc {
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
			return fmt.Errorf("error creating OpenTelekomCloud IMSv2 client: %w", err)
		}

		imageID := rs.Primary.Attributes["image_id"]
		found, err := members.Get(client, imageID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.MemberID != rs.Primary.ID {
			return fmt.Errorf("member not found")
		}

		*member = *found

		return nil
	}
}

func testAccImagesMemberV2Basic(projectToShare, privateImageID string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_images_member_v2" "member_1" {
  member   = "%s"
  image_id = "%s"
}
`, projectToShare, privateImageID)
}
