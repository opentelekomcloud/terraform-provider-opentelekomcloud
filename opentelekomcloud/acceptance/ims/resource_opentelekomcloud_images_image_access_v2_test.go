package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/imageservice/v2/members"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceName = "opentelekomcloud_images_image_access_v2.image_access_1"

func TestAccImagesImageAccessV2_basic(t *testing.T) {
	var member members.Member

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckImagesImageAccessV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageAccessV2Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageAccessV2Exists(resourceName, &member),
					resource.TestCheckResourceAttrPtr(resourceName, "status", &member.Status),
					resource.TestCheckResourceAttr(resourceName, "status", "pending"),
				),
			},
			{
				Config: testAccImagesImageAccessV2Update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageAccessV2Exists(resourceName, &member),
					resource.TestCheckResourceAttrPtr(resourceName, "status", &member.Status),
					resource.TestCheckResourceAttr(resourceName, "status", "accepted"),
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

		imageID, memberID, err := resourceImagesImageAccessV2ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = members.Get(client, imageID, memberID).Extract()
		if err == nil {
			return fmt.Errorf("image share still exists")
		}
	}

	return nil
}

func testAccCheckImagesImageAccessV2Exists(n string, member *members.Member) resource.TestCheckFunc {
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

		imageID, memberID, err := resourceImagesImageAccessV2ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := members.Get(client, imageID, memberID).Extract()
		if err != nil {
			return err
		}

		id := fmt.Sprintf("%s/%s", found.ImageID, found.MemberID)
		if id != rs.Primary.ID {
			return fmt.Errorf("image member not found")
		}

		*member = *found

		return nil
	}
}

const testAccImagesImageAccessV2 = `
data "opentelekomcloud_identity_auth_scope_v3" "scope" {
  name = "scope"
}
resource "opentelekomcloud_images_image_v2" "image_1" {
  name             = "CirrOS-tf_1"
  image_source_url = "https://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
  container_format = "bare"
  disk_format      = "qcow2"
  visibility       = "shared"
}`

func testAccImagesImageAccessV2Basic() string {
	return fmt.Sprintf(`
%s
resource "opentelekomcloud_images_image_access_v2" "image_access_1" {
  image_id  = opentelekomcloud_images_image_v2.image_1.id
  member_id = data.opentelekomcloud_identity_auth_scope_v3.scope.project_id
}
`, testAccImagesImageAccessV2)
}

func testAccImagesImageAccessV2Update() string {
	return fmt.Sprintf(`
%s
resource "opentelekomcloud_images_image_access_v2" "image_access_1" {
  image_id  = openstack_images_image_v2.image_1.id
  member_id = data.opentelekomcloud_identity_auth_scope_v3.scope.project_id
  status    = "accepted"
}
`, testAccImagesImageAccessV2)
}

func resourceImagesImageAccessV2ParseID(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		return "", "", fmt.Errorf("unable to determine image share access ID")
	}

	imageID := idParts[0]
	memberID := idParts[1]

	return imageID, memberID, nil
}
