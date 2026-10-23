package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
)

func TestLBMemberV3DataSource_basic(t *testing.T) {
	t.Parallel()
	quotas.BookMany(t, []*quotas.ExpectedQuota{
		{Q: quotas.LbPool, Count: 1},
		{Q: quotas.LoadBalancer, Count: 1},
	})
	const dataName = "data.opentelekomcloud_lb_member_v3.test"
	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(dataName, "member_id", resourceMemberName, "member_id"),
		resource.TestCheckResourceAttrPair(dataName, "address", resourceMemberName, "address"),
		resource.TestCheckResourceAttrPair(dataName, "id", resourceMemberName, "id"),
		resource.TestCheckResourceAttrPair(dataName, "member_type", resourceMemberName, "member_type"),
		resource.TestCheckResourceAttr(dataName, "weight", "0"),
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testLBMemberDestroy,
		Steps: []resource.TestStep{
			{Config: testLBMemberV3ZeroWeight + `
data "opentelekomcloud_lb_member_v3" "test" {
  pool_id   = opentelekomcloud_lb_member_v3.member.pool_id
  member_id = opentelekomcloud_lb_member_v3.member.member_id
}
`, Check: check},
			{Config: testLBMemberV3ZeroWeight + `
data "opentelekomcloud_lb_member_v3" "test" {
  pool_id = opentelekomcloud_lb_member_v3.member.pool_id
  name    = opentelekomcloud_lb_member_v3.member.name
  weight  = 0
}
`, Check: check},
		},
	})
}
