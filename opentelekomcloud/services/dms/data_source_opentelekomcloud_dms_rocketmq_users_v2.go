package dms

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/users"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func DataSourceDmsRocketMQUsersV2() *schema.Resource {
	permSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {Type: schema.TypeString, Computed: true},
			"perm": {Type: schema.TypeString, Computed: true},
		},
	}
	return &schema.Resource{
		ReadContext: dataSourceDmsRocketMQUsersV2Read,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key":           {Type: schema.TypeString, Computed: true},
						"white_remote_address": {Type: schema.TypeString, Computed: true},
						"admin":                {Type: schema.TypeBool, Computed: true},
						"default_topic_perm":   {Type: schema.TypeString, Computed: true},
						"default_group_perm":   {Type: schema.TypeString, Computed: true},
						"topic_perms":          {Type: schema.TypeList, Computed: true, Elem: permSchema},
						"group_perms":          {Type: schema.TypeList, Computed: true, Elem: permSchema},
					},
				},
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func flattenRocketMQUserPerms(perms []users.Permission) []map[string]interface{} {
	result := make([]map[string]interface{}, len(perms))
	for i, p := range perms {
		result[i] = map[string]interface{}{"name": p.Name, "perm": p.Perm}
	}
	return result
}

func dataSourceDmsRocketMQUsersV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)
	client, err := config.DmsV2Client(region)
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID := d.Get("instance_id").(string)
	opts := users.ListOpts{Limit: 50}
	var all []users.User
	for {
		resp, err := users.List(client, instanceID, opts)
		if err != nil {
			return diag.Errorf("error getting users of DMS RocketMQ instance (%s): %s", instanceID, err)
		}
		all = append(all, resp.Users...)
		if len(resp.Users) == 0 || len(all) >= resp.Total {
			break
		}
		opts.Offset += len(resp.Users)
	}

	accessKey := d.Get("access_key").(string)
	ids := []string{instanceID}
	result := make([]map[string]interface{}, 0, len(all))
	for _, v := range all {
		if accessKey != "" && v.AccessKey != accessKey {
			continue
		}
		ids = append(ids, v.AccessKey)
		result = append(result, map[string]interface{}{
			"access_key":           v.AccessKey,
			"white_remote_address": v.WhiteRemoteAddress,
			"admin":                v.Admin,
			"default_topic_perm":   v.DefaultTopicPerm,
			"default_group_perm":   v.DefaultGroupPerm,
			"topic_perms":          flattenRocketMQUserPerms(v.TopicPerms),
			"group_perms":          flattenRocketMQUserPerms(v.GroupPerms),
		})
	}
	log.Printf("[DEBUG] DMS RocketMQ users: %v", ids[1:])
	d.SetId(hashcode.Strings(ids))

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("users", result),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
