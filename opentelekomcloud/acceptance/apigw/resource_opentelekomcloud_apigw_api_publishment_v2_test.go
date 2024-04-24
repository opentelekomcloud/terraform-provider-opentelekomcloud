package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	apis "github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/api"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/apigw"
)

const resourceApigwApiPublishName = "opentelekomcloud_apigw_api_publishment_v2.pub"

func getPublishmentResourceFunc(cfg *cfg.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.APIGWV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return apigw.GetVersionHistory(client,
		state.Primary.Attributes["gateway_id"],
		state.Primary.Attributes["environment_id"],
		state.Primary.Attributes["api_id"])
}

func TestAccApiPublishment_basic(t *testing.T) {
	var history []apis.VersionResp
	rName := fmt.Sprintf("apigw_acc_api_pub%s", acctest.RandString(5))
	rc := common.InitResourceCheck(
		resourceApigwApiPublishName,
		&history,
		getPublishmentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			common.TestAccPreCheck(t)
		},
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApiPublishment_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceApigwApiPublishName, "environment_name"),
					resource.TestCheckResourceAttrSet(resourceApigwApiPublishName, "published_at"),
					resource.TestCheckResourceAttrSet(resourceApigwApiPublishName, "publish_id"),
				),
			},
			{
				ResourceName:      resourceApigwApiPublishName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccApiPublishment_basic(name string) string {
	relatedConfig := testAccApigwApi_basic(testAccApigwApi_base(name), name)

	return fmt.Sprintf(`
%[1]s

resource "opentelekomcloud_apigw_api_publishment_v2" "pub" {
  gateway_id     = opentelekomcloud_apigw_gateway_v2.gateway.id
  environment_id = opentelekomcloud_apigw_environment_v2.env.id
  api_id         = opentelekomcloud_apigw_api_v2.api.id
}
`, relatedConfig, name)
}
