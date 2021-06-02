package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestAccDDSInstanceV3_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_dds_instance_v3.instance"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDDSV3InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccDDSInstanceV3Config_basic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"flavor",
					"password",
					"availability_zone",
				},
			},
		},
	})
}
