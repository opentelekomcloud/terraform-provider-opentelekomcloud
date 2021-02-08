package opentelekomcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
)

// tagsSchema returns the schema to use for tags.
func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
	}
}

func tagsSchemaComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Computed: true,
	}
}

// UpdateResourceTags is a helper to update the tags for a resource.
// It expects the tags field to be named "tags"
func UpdateResourceTags(client *golangsdk.ServiceClient, d *schema.ResourceData, resourceType, id string) error {
	if d.HasChange("tags") {
		oldMapRaw, newMapRaw := d.GetChange("tags")
		oldMap := oldMapRaw.(map[string]interface{})
		newMap := newMapRaw.(map[string]interface{})

		// remove old tags
		if len(oldMap) > 0 {
			tagList := expandResourceTags(oldMap)
			err := tags.Delete(client, resourceType, id, tagList).ExtractErr()
			if err != nil {
				return err
			}
		}

		// set new tags
		if len(newMap) > 0 {
			tagList := expandResourceTags(newMap)
			err := tags.Create(client, resourceType, id, tagList).ExtractErr()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// tagsToMap returns the list of tags into a map.
func tagsToMap(tags []tags.ResourceTag) map[string]string {
	result := make(map[string]string)
	for _, val := range tags {
		result[val.Key] = val.Value
	}

	return result
}

// expandResourceTags returns the tags for the given map of data.
func expandResourceTags(tagMap map[string]interface{}) []tags.ResourceTag {
	var tagList []tags.ResourceTag

	for k, v := range tagMap {
		tag := tags.ResourceTag{
			Key:   k,
			Value: v.(string),
		}
		tagList = append(tagList, tag)
	}

	return tagList
}
