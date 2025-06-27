package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/smn/v2/templates"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getSmnMessageTemplateResourceFunc(conf *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.SmnV2Client(cfg.ProjectName(env.OS_REGION_NAME))
	if err != nil {
		return nil, fmt.Errorf("error creating SmnV2 client: %s", err)
	}

	return templates.Get(client, state.Primary.ID)
}

func TestAccSmnMessageTemplate_basic(t *testing.T) {
	var obj interface{}

	name := fmt.Sprintf("smn-acc-api%s", acctest.RandString(5))
	rName := "opentelekomcloud_smn_message_template_v2.test"

	rc := common.InitResourceCheck(
		rName,
		&obj,
		getSmnMessageTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSmnMessageTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "protocol", "default"),
					resource.TestCheckResourceAttr(rName, "content",
						"Test content, contains {content1} and {content2}"),
					resource.TestCheckResourceAttr(rName, "tag_names.#", "2"),
					resource.TestCheckResourceAttr(rName, "tag_names.0", "content1"),
					resource.TestCheckResourceAttr(rName, "tag_names.1", "content2"),
				),
			},
			{
				Config: testSmnMessageTemplate_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "protocol", "default"),
					resource.TestCheckResourceAttr(rName, "content",
						"Test update content, contains {content1}, {content2} and {content3}"),
					resource.TestCheckResourceAttr(rName, "tag_names.#", "3"),
					resource.TestCheckResourceAttr(rName, "tag_names.0", "content1"),
					resource.TestCheckResourceAttr(rName, "tag_names.1", "content2"),
					resource.TestCheckResourceAttr(rName, "tag_names.2", "content3"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testSmnMessageTemplate_basic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_smn_message_template_v2" "test" {
  name     = "%s"
  protocol = "default"
  content  = "Test content, contains {content1} and {content2}"
}
`, name)
}

func testSmnMessageTemplate_basic_update(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_smn_message_template_v2" "test" {
  name     = "%s"
  protocol = "default"
  content  = "Test update content, contains {content1}, {content2} and {content3}"
}
`, name)
}
