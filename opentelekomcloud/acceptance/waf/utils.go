package acceptance

import (
	"fmt"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func getAccWafdClient(client *golangsdk.ServiceClient, err error, config *cfg.Config) (*golangsdk.ServiceClient, error) {
	if env.OS_REGION_NAME != "eu-ch2" {
		client, err = config.WafDedicatedV1Client(env.OS_REGION_NAME)
		if err != nil {
			return nil, fmt.Errorf("error creating OpenTelekomCloud Waf dedicated client: %s", err)
		}
	} else {
		client, err = config.WafDedicatedSwissV1Client(env.OS_REGION_NAME)
		if err != nil {
			return nil, fmt.Errorf("error creating OpenTelekomCloud Waf dedicated client: %s", err)
		}
	}
	return client, err
}
