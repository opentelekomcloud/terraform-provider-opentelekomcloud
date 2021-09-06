package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccCTSTrackerV1_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_cts_tracker_v1.tracker_v1"
	var bucketName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCTSTrackerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV1Basic(bucketName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
