package acceptance

import (
	"os"
	"sync"
	"testing"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/css/v1/flavors"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const sharedFlavorName = "css.medium.8"

var (
	flavor   *flavors.Flavor
	flavOnce sync.Once
)

func findSharedFlavor(t *testing.T) {
	t.Helper()
	if os.Getenv("TF_ACC") == "" {
		t.Skip("findSharedFlavor can be used only in acceptance tests")
	}
	flavOnce.Do(func() {
		t.Helper()
		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.CssV1Client(env.OS_REGION_NAME)
		if err != nil {
			t.Fatalf("error creating CSSv1 client: %s", err)
		}

		versions, err := flavors.List(client)
		if err != nil {
			t.Fatalf("error extracting versions: %s", err)
		}

		flavor = flavors.FindFlavor(versions, flavors.FilterOpts{FlavorName: sharedFlavorName})
	})
}

func sharedFlavorQuotas(t *testing.T, nodeCount int64, volumeSize int64) quotas.MultipleQuotas {
	t.Helper()
	findSharedFlavor(t)

	qts := quotas.MultipleQuotas{
		{Q: quotas.Server, Count: 1},
		{Q: quotas.Volume, Count: 1},
		{Q: quotas.VolumeSize, Count: volumeSize},
		{Q: quotas.RAM, Count: int64(flavor.RAM * 1024)},
		{Q: quotas.CPU, Count: int64(flavor.CPU)},
	}
	qts.X(nodeCount)
	return qts
}
