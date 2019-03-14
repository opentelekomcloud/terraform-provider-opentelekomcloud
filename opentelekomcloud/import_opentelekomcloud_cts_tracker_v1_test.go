package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccCTSTrackerV1_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_cts_tracker_v1.tracker_v1"
	var bucketName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCTSTrackerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV1_basic(bucketName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
