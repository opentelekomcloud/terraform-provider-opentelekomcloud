package acceptance

import "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

func lbQuotas() []*quotas.ExpectedQuota {
	return []*quotas.ExpectedQuota{
		{Q: quotas.LoadBalancer, Count: 1},
		{Q: quotas.FloatingIP, Count: 1},
	}
}
