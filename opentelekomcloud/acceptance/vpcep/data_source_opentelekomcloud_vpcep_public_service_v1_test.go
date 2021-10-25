package vpcep

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

const dataSourcePublicService = "data.opentelekomcloud_vpcep_public_service_v1.obs"

func TestDataSourcePublicService(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      checkServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePublicService,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourcePublicService, "owner", "OTC"),
					resource.TestCheckResourceAttr(dataSourcePublicService, "service_type", "gateway"),
				),
			},
		},
	})
}

var testDataSourcePublicService = fmt.Sprintf(`
data "opentelekomcloud_vpcep_public_service_v1" "obs" {
  name = "com.t-systems.otc.%s.obs"
}
`, env.OS_REGION_NAME)
