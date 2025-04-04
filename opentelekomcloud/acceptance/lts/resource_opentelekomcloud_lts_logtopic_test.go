package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/streams"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getLtsStreamResourceFunc(config *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.LtsV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS v2 client: %s", err)
	}

	requestResp, err := streams.List(client, state.Primary.Attributes["group_id"])
	if err != nil {
		return nil, err
	}
	if len(requestResp) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	var streamResult streams.LogStream
	for _, stream := range requestResp {
		if stream.LogStreamId == state.Primary.ID {
			streamResult = stream
		}
	}
	if streamResult.LogStreamId == "" {
		return nil, golangsdk.ErrDefault404{}
	}
	return streamResult, nil
}

func TestAccLogTankTopicV2_basic(t *testing.T) {
	var (
		topic        streams.LogStream
		resourceName = "opentelekomcloud_logtank_topic_v2.topic"
		rName        = fmt.Sprintf("lts_topic%s", acctest.RandString(3))
		rc           = common.InitResourceCheck(resourceName, &topic, getLtsStreamResourceFunc)
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLogTankTopicV2_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(
						resourceName, "topic_name", rName),
				),
			},
		},
	})
}

func testAccLogTankTopicV2_basic(name string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_logtank_group_v2" "group" {
  group_name  = "%[1]s"
  ttl_in_days = 7
}

resource "opentelekomcloud_logtank_topic_v2" "topic" {
  group_id   = opentelekomcloud_logtank_group_v2.group.id
  topic_name = "%[1]s"
}
`, name)
}
