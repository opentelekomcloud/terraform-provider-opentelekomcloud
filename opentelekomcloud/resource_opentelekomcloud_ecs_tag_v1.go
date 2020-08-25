package opentelekomcloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tags "github.com/huaweicloud/golangsdk/openstack/ecs/v1/cloudservertags"
)

func setTagForInstance(d *schema.ResourceData, meta interface{}, instanceID string, tagsMap map[string]interface{}) error {
	config := meta.(*Config)
	client, err := config.computeV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud compute v1 client: %s", err)
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
		return fmt.Errorf("Error creating OpenTelekomCloud instance tags: %s", createTags.Err)
	}

	return nil
}
