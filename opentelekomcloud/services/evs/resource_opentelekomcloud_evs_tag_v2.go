package evs

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/evs/v2/tags"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func resourceEVSTagV2Create(_ context.Context, d *schema.ResourceData, meta interface{}, resourceType, resourceID string, tag map[string]string) (*tags.Tags, error) {
	config := meta.(*cfg.Config)
	client, err := config.BlockStorageV2Client(config.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud client: %s", err)
	}

	createOpts := tags.CreateOpts{Tags: tag}
	return tags.Create(client, resourceType, resourceID, createOpts).Extract()
}

func resourceEVSTagV2Get(d *schema.ResourceData, meta interface{}, resourceType, resourceID string) (*tags.Tags, error) {
	config := meta.(*cfg.Config)
	client, err := config.BlockStorageV2Client(config.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud client: %s", err)
	}

	return tags.Get(client, resourceType, resourceID).Extract()
}
