package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDNSV2RecordSet_importBasic(t *testing.T) {
	zoneName := randomZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSetBasic(zoneName),
			},
			{
				ResourceName:      resourceRecordSetName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
