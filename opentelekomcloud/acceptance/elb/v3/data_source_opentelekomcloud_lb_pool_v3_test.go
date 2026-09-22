package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
)

const dataSourcePoolName = "data.opentelekomcloud_lb_pool_v3.pool"

func TestLBPoolV3DataSource_basic(t *testing.T) {
	t.Parallel()
	quotas.BookMany(t, []*quotas.ExpectedQuota{
		{Q: quotas.LbPool, Count: 1},
		{Q: quotas.LoadBalancer, Count: 1},
	})

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testLBPoolV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testLBPoolV3Basic + `
data "opentelekomcloud_lb_pool_v3" "pool" {
  id = opentelekomcloud_lb_pool_v3.pool.id
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourcePoolName, "id", resourcePoolName, "id"),
					resource.TestCheckResourceAttrPair(dataSourcePoolName, "protocol", resourcePoolName, "protocol"),
					resource.TestCheckResourceAttrPair(dataSourcePoolName, "lb_algorithm", resourcePoolName, "lb_algorithm"),
					resource.TestCheckResourceAttrPair(dataSourcePoolName, "loadbalancer_id", resourcePoolName, "loadbalancer_id"),
					resource.TestCheckResourceAttr(dataSourcePoolName, "loadbalancer_ids.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourcePoolName, "created_at"),
					resource.TestCheckResourceAttrSet(dataSourcePoolName, "updated_at"),
				),
			},
			{
				Config: testLBPoolV3Basic + `
data "opentelekomcloud_lb_pool_v3" "pool" {
  loadbalancer_id = opentelekomcloud_lb_pool_v3.pool.loadbalancer_id
  protocol        = opentelekomcloud_lb_pool_v3.pool.protocol
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourcePoolName, "id", resourcePoolName, "id"),
					resource.TestCheckResourceAttrPair(dataSourcePoolName, "vpc_id", resourcePoolName, "vpc_id"),
					resource.TestCheckResourceAttrPair(dataSourcePoolName, "type", resourcePoolName, "type"),
				),
			},
		},
	})
}
