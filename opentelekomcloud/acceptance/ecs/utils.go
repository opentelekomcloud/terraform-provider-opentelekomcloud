package acceptance

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/flavors"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/servers"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccCheckComputeV2InstanceDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	computeClient, err := config.ComputeV2Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud ComputeV2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_compute_instance_v2" {
			continue
		}

		server, err := servers.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			if server.Status != "SOFT_DELETED" {
				return fmt.Errorf("instance still exists")
			}
		}
	}

	return nil
}

func getFlavors() (map[string][]*quotas.ExpectedQuota, error) {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.ComputeV2Client(env.OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud ComputeV2 client: %s", err)
	}
	pgs, err := flavors.ListDetail(client, flavors.ListOpts{}).AllPages()
	if err != nil {
		return nil, fmt.Errorf("error listing flavors pages: %w", err)
	}
	flavs, err := flavors.ExtractFlavors(pgs)
	if err != nil {
		return nil, fmt.Errorf("error extracting flavors pages: %w", err)
	}
	resultsQ := map[string][]*quotas.ExpectedQuota{}
	for _, flv := range flavs {
		exp := []*quotas.ExpectedQuota{
			{
				Q:     quotas.CPU,
				Count: int64(flv.VCPUs),
			},
			{
				Q:     quotas.RAM,
				Count: int64(flv.RAM),
			},
		}
		resultsQ[flv.ID] = exp
	}
	return resultsQ, nil
}

var flavorsQuota map[string][]*quotas.ExpectedQuota

func init() {
	qs, err := getFlavors()
	if err != nil {
		panic("failed to get server flavors")
	}
	flavorsQuota = qs
}

func quotasForFlavor(flavorRef string) []*quotas.ExpectedQuota {
	var qts []*quotas.ExpectedQuota
	qts = append(qts,
		flavorsQuota[flavorRef]...,
	)
	return qts
}

func serverQuotas(volume int64, flavor string) []*quotas.ExpectedQuota {
	qts := []*quotas.ExpectedQuota{
		{Q: quotas.Volume, Count: 1},
		{Q: quotas.VolumeSize, Count: volume},
		{Q: quotas.Server, Count: 1},
	}
	qts = append(qts, quotasForFlavor(flavor)...)
	return qts
}

func simpleServerWithIPQuotas(eipCount int64) []*quotas.ExpectedQuota {
	qts := serverQuotas(4, env.OsFlavorID)
	qts = append(qts, &quotas.ExpectedQuota{Q: quotas.FloatingIP, Count: eipCount})
	return qts
}
