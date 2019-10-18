package opentelekomcloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/evs/v2/tags"
)

func resourceEVSTagV2Create(d *schema.ResourceData, meta interface{}, resourceType, resourceID string, tag map[string]string) (*tags.Tags, error) {
	config := meta.(*Config)
	client, err := config.loadEVSV2Client(GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating OpenTelekomCloud client: %s", err)
	}

	createOpts := tags.CreateOpts{Tags: tag}
	return tags.Create(client, resourceType, resourceID, createOpts).Extract()
}

func resourceEVSTagV2Get(d *schema.ResourceData, meta interface{}, resourceType, resourceID string) (*tags.Tags, error) {
	config := meta.(*Config)
	client, err := config.loadEVSV2Client(GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating OpenTelekomCloud client: %s", err)
	}

	return tags.Get(client, resourceType, resourceID).Extract()
}
