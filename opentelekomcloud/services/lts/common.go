package lts

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rt "github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/tags"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const (
	keyClientV2         = "lts-v2-client"
	errCreationV2Client = "error creating OpenTelekomCloud LTS V2 client: %w"
	errCreationV1Client = "error creating OpenTelekomCloud LTS V1 client: %w"
)

func ltsTags(d *schema.ResourceData) []rt.ResourceTag {
	t := d.Get("tags").(map[string]interface{})
	var tagSlice []rt.ResourceTag
	for k, v := range t {
		tagSlice = append(tagSlice, rt.ResourceTag{Key: k, Value: v.(string)})
	}
	return tagSlice
}

func ignoreSysEpsTag(tags map[string]string) map[string]string {
	delete(tags, "_sys_enterprise_project_id")
	return tags
}

func updateTags(d *schema.ResourceData, meta interface{}, resourceType, id string) error {
	config := meta.(*cfg.Config)
	client, err := config.LtsV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf(errCreationV1Client, err)
	}
	if d.HasChange("tags") {
		oldMapRaw, newMapRaw := d.GetChange("tags")
		oldMap := oldMapRaw.(map[string]interface{})
		newMap := newMapRaw.(map[string]interface{})

		// remove old tags
		if len(oldMap) > 0 {
			tagList := common.ExpandResourceTags(oldMap)
			err := tags.Manage(client, resourceType, id, tags.TagOpts{
				Action: "delete",
				IsOpen: true,
				Tags:   tagList,
			})
			if err != nil {
				return err
			}
		}

		// set new tags
		if len(newMap) > 0 {
			tagList := common.ExpandResourceTags(newMap)
			err := tags.Manage(client, resourceType, id, tags.TagOpts{
				Action: "create",
				IsOpen: true,
				Tags:   tagList,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
