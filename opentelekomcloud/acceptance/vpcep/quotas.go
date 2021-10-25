package vpcep

import "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

var (
	serviceQuota  = quotas.NewQuota(20)
	endpointQuota = quotas.NewQuota(50)
)

func endpointQuotas() quotas.MultipleQuotas {
	return quotas.MultipleQuotas{
		{Q: serviceQuota, Count: 1},
		{Q: endpointQuota, Count: 1},
	}
}
