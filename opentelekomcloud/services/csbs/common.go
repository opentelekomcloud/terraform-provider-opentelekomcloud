package csbs

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
)

func resourceCSBSTagsV1(d *schema.ResourceData) []tags.ResourceTag {
	backupTags := d.Get("tags").(map[string]interface{})
	var tagSlice []tags.ResourceTag
	for k, v := range backupTags {
		tagSlice = append(tagSlice, tags.ResourceTag{Key: k, Value: v.(string)})
	}
	return tagSlice
}
