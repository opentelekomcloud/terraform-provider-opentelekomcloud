package acceptance

import (
	"testing"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	acceptance "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/ecs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func testAccCCEKeyPairPreCheck(t *testing.T) {
	common.TestAccPreCheckRequiredEnvVars(t)
	if env.OS_KEYPAIR_NAME == "" {
		t.Skip("OS_KEYPAIR_NAME must be set for acceptance tests")
	}
}

var clusterNodesQuota = quotas.NewQuota(5)
var singleNodeQuotas = func() quotas.MultipleQuotas {
	qts := acceptance.QuotasForFlavor("s3.medium.1")
	qts = append(qts,
		&quotas.ExpectedQuota{Q: clusterNodesQuota, Count: 1},
		&quotas.ExpectedQuota{Q: quotas.Server, Count: 1},
	)
	return qts
}()
