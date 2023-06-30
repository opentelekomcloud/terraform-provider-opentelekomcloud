package csbs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
)

func resourceCSBSTagsV1(d *schema.ResourceData) []tags.ResourceTag {
	rawTags := d.Get("tags").(*schema.Set).List()
	tagsRaw := make([]tags.ResourceTag, len(rawTags))
	for i, raw := range rawTags {
		rawMap := raw.(map[string]interface{})
		tagsRaw[i] = tags.ResourceTag{
			Key:   rawMap["key"].(string),
			Value: rawMap["value"].(string),
		}
	}
	return tagsRaw
}

func flattenCSBSTags(resourceTags []tags.ResourceTag) []map[string]interface{} {
	var tagsList []map[string]interface{}
	for _, tag := range resourceTags {
		mapping := map[string]interface{}{
			"key":   tag.Key,
			"value": tag.Value,
		}
		tagsList = append(tagsList, mapping)
	}

	return tagsList
}
