package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const resourceFloatingIpAssociateName = "opentelekomcloud_compute_floatingip_associate_v2.fip_1"

func TestAccComputeV2FloatingIPAssociate_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2FloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FloatingIPAssociateBasic,
			},
			{
				ResourceName:      resourceFloatingIpAssociateName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
