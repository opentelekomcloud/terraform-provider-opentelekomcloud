package ecs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tags "github.com/opentelekomcloud/gophertelekomcloud/openstack/ecs/v1/cloudservertags"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func SetTagForInstance(d *schema.ResourceData, meta interface{}, instanceID string, tagsMap map[string]interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.ComputeV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud compute v1 client: %s", err)
	}

	rId := instanceID
	var tagsList []tags.Tag
	for k, v := range tagsMap {
		tag := tags.Tag{
			Key:   k,
			Value: v.(string),
		}
		tagsList = append(tagsList, tag)
	}

	createOpts := tags.BatchOpts{Action: tags.ActionCreate, Tags: tagsList}
	createTags := tags.BatchAction(client, rId, createOpts)
	if createTags.Err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud instance tags: %s", createTags.Err)
	}

	return nil
}
