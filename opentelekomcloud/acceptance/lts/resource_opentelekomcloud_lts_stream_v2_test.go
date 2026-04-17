package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/streams"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccLtsStream_basic(t *testing.T) {
	var (
		stream       streams.LogStream
		resourceName = "opentelekomcloud_lts_stream_v2.stream"
		rName        = fmt.Sprintf("lts_stream%s", acctest.RandString(3))
		rc           = common.InitResourceCheck(resourceName, &stream, getLtsStreamResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLtsV2Stream_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "stream_name", rName),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "-1"),
					resource.TestCheckResourceAttr(resourceName, "filter_count", "0"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrPair(resourceName, "group_id", "opentelekomcloud_lts_group_v2.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.terraform", ""),
				),
			},
			{
				Config: testAccLtsV2Stream_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "60"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testLtsStreamImportState(resourceName),
			},
		},
	})
}

func testLtsStreamImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		streamID := rs.Primary.ID
		groupID := rs.Primary.Attributes["group_id"]

		return fmt.Sprintf("%s/%s", groupID, streamID), nil
	}
}

func testAccLtsV2Stream_basic(rName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_group_v2" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30

  tags = {
    owner = "terraform"
  }
}

resource "opentelekomcloud_lts_stream_v2" "stream" {
  group_id    = opentelekomcloud_lts_group_v2.test.id
  stream_name = "%[1]s"

  tags = {
    foo       = "bar"
    terraform = ""
  }
}
`, rName)
}

func testAccLtsV2Stream_update(rName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_group_v2" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30

  tags = {
    owner = "terraform"
  }
}

resource "opentelekomcloud_lts_stream_v2" "stream" {
  group_id    = opentelekomcloud_lts_group_v2.test.id
  stream_name = "%[1]s"
  ttl_in_days = 60

  tags = {
    owner = "terraform"
  }
}
`, rName)
}
