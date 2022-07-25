package common

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
)

// TagsSchema returns the schema to use for tags.
func TagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeMap,
		Optional:     true,
		ValidateFunc: ValidateTags,
		Elem:         &schema.Schema{Type: schema.TypeString},
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
			tagList := ExpandResourceTags(oldMap)
			err := tags.Delete(client, resourceType, id, tagList).ExtractErr()
			if err != nil {
				return err
			}
		}

		// set new tags
		if len(newMap) > 0 {
			tagList := ExpandResourceTags(newMap)
			err := tags.Create(client, resourceType, id, tagList).ExtractErr()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// TagsToMap returns the list of tags into a map.
func TagsToMap(tags []tags.ResourceTag) map[string]string {
	result := make(map[string]string)
	for _, val := range tags {
		result[val.Key] = val.Value
	}

	return result
}

// ExpandResourceTags returns the tags for the given map of data.
func ExpandResourceTags(tagMap map[string]interface{}) []tags.ResourceTag {
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

func Contains(tagSlice []tags.ResourceTag, tag tags.ResourceTag) bool {
	for _, v := range tagSlice {
		if v == tag {
			return true
		}
	}

	return false
}
