package acceptance

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const resourceFloatingIpAssociateName = "opentelekomcloud_compute_floatingip_associate_v2.fip_1"

func TestAccComputeV2FloatingIPAssociate_importBasic(t *testing.T) {
	t.Parallel()
	qts := simpleServerWithIPQuotas(1)
	th.AssertNoErr(t, quotas.AcquireMultipleQuotas(qts, 5*time.Second))
	defer quotas.ReleaseMultipleQuotas(qts)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
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
