package opentelekomcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/evs/v2/tags"
)

func resourceEVSTagV2Create(d *schema.ResourceData, meta interface{}, resource_type, resource_id string, tag map[string]string) (*tags.Tags, error) {
	config := meta.(*Config)
	client, err := chooseEVSV2Client(d, config)
	if err != nil {
		return nil, fmt.Errorf("Error creating OpenTelekomCloud client: %s", err)
	}

	createOpts := tags.CreateOpts{Tags: tag}
	return tags.Create(client, resource_type, resource_id, createOpts).Extract()
}

func resourceEVSTagV2Get(d *schema.ResourceData, meta interface{}, resource_type, resource_id string) (*tags.Tags, error) {
	config := meta.(*Config)
	client, err := chooseEVSV2Client(d, config)
	if err != nil {
		return nil, fmt.Errorf("Error creating OpenTelekomCloud client: %s", err)
	}

	return tags.Get(client, resource_type, resource_id).Extract()
}
