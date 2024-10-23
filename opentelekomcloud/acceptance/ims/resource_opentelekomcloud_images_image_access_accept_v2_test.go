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

func TestAccImagesImageAccessAcceptV2_basic(t *testing.T) {
	var member ims2.Member
	acceptResourceName := "opentelekomcloud_images_image_access_accept_v2.accept_1"
	shareProjectID := os.Getenv("OS_PROJECT_ID_2")
	shareProjectName := os.Getenv("OS_PROJECT_NAME_2")
	shareCloudID := os.Getenv("OS_CLOUD_2")
	if shareProjectID == "" || shareCloudID == "" || shareProjectName == "" {
		t.Skip("OS_PRIVATE_IMAGE_ID or OS_PROJECT_ID_2 or OS_CLOUD_2 are empty, but test requires")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckImagesImageAccessAcceptV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageAccessAcceptV2Basic(shareProjectID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageAccessV2Exists(acceptResourceName, &member),
					resource.TestCheckResourceAttrPtr(acceptResourceName, "status", &member.Status),
					resource.TestCheckResourceAttr(acceptResourceName, "status", "accepted"),
				),
			},
			{
				Config: testAccImagesImageAccessAcceptV2Update(shareProjectID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageAccessV2Exists(acceptResourceName, &member),
					resource.TestCheckResourceAttrPtr(acceptResourceName, "status", &member.Status),
					resource.TestCheckResourceAttr(acceptResourceName, "status", "rejected"),
				),
			},
			{
				ResourceName:      acceptResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckImagesImageAccessAcceptV2Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ImageV2Client(env.OS_REGION_NAME)
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

		_, err = members.Get(client, members.MemberOpts{
			ImageId:  imageID,
			MemberId: memberID,
		})
		if err == nil {
			return fmt.Errorf("image membership still exists")
		}
	}

	return nil
}

func testAccImagesImageAccessAcceptV2Basic(projectToShare string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_images_image_v2" "rancheros" {
  name             = "RancherOS"
  image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
  container_format = "bare"
  disk_format      = "qcow2"
}

resource "opentelekomcloud_images_image_access_v2" "access_1" {
  provider = "opentelekomcloud"

  image_id  = opentelekomcloud_images_image_v2.rancheros.id
  member_id = "%[1]s"
}

%[2]s

resource "opentelekomcloud_images_image_access_accept_v2" "accept_1" {
  provider = "%[3]s"

  depends_on = [opentelekomcloud_images_image_access_v2.access_1]

  image_id  = opentelekomcloud_images_image_v2.rancheros.id
  member_id = "%[1]s"
  status    = "accepted"
}
`, projectToShare, common.AlternativeProviderConfig, common.AlternativeProviderAlias)
}

func testAccImagesImageAccessAcceptV2Update(projectToShare string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_images_image_v2" "rancheros" {
  name             = "RancherOS"
  image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
  container_format = "bare"
  disk_format      = "qcow2"
}

resource "opentelekomcloud_images_image_access_v2" "access_1" {
  provider = "opentelekomcloud"

  image_id  = opentelekomcloud_images_image_v2.rancheros.id
  member_id = "%[1]s"
}

%[2]s

resource "opentelekomcloud_images_image_access_accept_v2" "accept_1" {
  provider = "%[3]s"

  depends_on = [opentelekomcloud_images_image_access_v2.access_1]

  image_id  = opentelekomcloud_images_image_v2.rancheros.id
  member_id = "%[1]s"
  status    = "rejected"
}
`, projectToShare, common.AlternativeProviderConfig, common.AlternativeProviderAlias)
}
