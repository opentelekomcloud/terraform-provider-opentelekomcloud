package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDNSV2Zone_importBasic(t *testing.T) {
	var zoneName = fmt.Sprintf("accepttest%s.com.", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDNSV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2ZoneBasic(zoneName),
			},

			{
				ResourceName:      resourceZoneName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
