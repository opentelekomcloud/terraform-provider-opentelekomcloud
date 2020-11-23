package opentelekomcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDDSInstanceV3_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_dds_instance_v3.instance"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCDeHV1Destroy,
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
