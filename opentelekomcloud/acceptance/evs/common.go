package acceptance

import "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

func volumeQuotas(size int64) quotas.MultipleQuotas {
	return quotas.MultipleQuotas{
		{Q: quotas.Volume, Count: 1},
		{Q: quotas.VolumeSize, Count: size},
	}
}
