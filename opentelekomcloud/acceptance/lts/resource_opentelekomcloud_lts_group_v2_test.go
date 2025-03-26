package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/groups"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccLtsV2Group_basic(t *testing.T) {
	var (
		group        groups.LogGroup
		resourceName = "opentelekomcloud_lts_group_v2.group"
		rName        = fmt.Sprintf("lts_group%s", acctest.RandString(3))
		rc           = common.InitResourceCheck(resourceName, &group, getLtsGroupResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLtsGroup_basic(rName, 30),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "group_name", rName),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "30"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLtsGroup_basic(rName, 7),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_name", rName),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "7"),
				),
			},
			{
				Config: testAccLtsGroup_update(rName, 60),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_name", rName),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "60"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.terraform", ""),
				),
			},
			{
				Config: testAccLtsGroup_update2(rName, 60),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func testAccLtsGroup_basic(name string, ttl int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_group_v2" "group" {
  group_name  = "%s"
  ttl_in_days = %d

  tags = {
    owner = "terraform"
  }
}
`, name, ttl)
}

func testAccLtsGroup_update(name string, ttl int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_group_v2" "group" {
  group_name  = "%s"
  ttl_in_days = %d

  tags = {
    foo       = "bar"
    key       = "value"
    terraform = ""
  }
}
`, name, ttl)
}

func testAccLtsGroup_update2(name string, ttl int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_lts_group_v2" "group" {
  group_name  = "%s"
  ttl_in_days = %d
}
`, name, ttl)
}
