package vbs

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vbs/v2/policies"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vbs/v2/tags"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceVBSBackupPolicyV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVBSPolicyV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"frequency": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"remain_first_backup": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rentention_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"policy_resource_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"filter_tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
					},
				},
			},
		},
	}
}

func dataSourceVBSPolicyV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	vbsClient, err := config.VbsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating VBSv2 client: %w", err)
	}

	policyID := d.Id()
	rawTags := d.Get("filter_tags").(*schema.Set).List()
	if len(rawTags) > 0 {
		tagsOpts := tags.ListOpts{Action: "filter", Tags: getVBSFilterTagsV2(d)}
		querytags, err := tags.ListResources(vbsClient, tagsOpts).ExtractResources()
		if err != nil {
			return fmterr.Errorf("error Querying backup policy using tags: %s ", err)
		}
		if querytags.TotalCount > 1 {
			return fmterr.Errorf("your tags query returned more than one result." +
				" Please try a more specific search criteria")
		}
		if querytags.TotalCount > 0 {
			policyID = querytags.Resource[0].ResourceID
		}
	}
	listOpts := policies.ListOpts{
		ID:     policyID,
		Name:   d.Get("name").(string),
		Status: d.Get("status").(string),
	}

	refinedPolicies, err := policies.List(vbsClient, listOpts)
	if err != nil {
		return fmterr.Errorf("unable to retrieve policies: %s", err)
	}

	if len(refinedPolicies) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedPolicies) > 1 {
		return fmterr.Errorf("your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Policy := refinedPolicies[0]

	log.Printf("[INFO] Retrieved Policy using given filter %s: %+v", Policy.ID, Policy)
	d.SetId(Policy.ID)

	mErr := multierror.Append(
		d.Set("name", Policy.Name),
		d.Set("policy_resource_count", Policy.ResourceCount),
		d.Set("frequency", Policy.ScheduledPolicy.Frequency),
		d.Set("remain_first_backup", Policy.ScheduledPolicy.RemainFirstBackup),
		d.Set("rentention_num", Policy.ScheduledPolicy.RententionNum),
		d.Set("start_time", Policy.ScheduledPolicy.StartTime),
		d.Set("status", Policy.ScheduledPolicy.Status),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	n, err := tags.Get(vbsClient, Policy.ID).Extract()
	if err != nil {
		return fmterr.Errorf("error getting tags: %w", err)
	}
	var tag []map[string]interface{}
	for _, policy := range n.Tags {
		mapping := map[string]interface{}{
			"key":   policy.Key,
			"value": policy.Value,
		}
		tag = append(tag, mapping)
	}

	if err := d.Set("tags", tag); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getVBSFilterTagsV2(d *schema.ResourceData) []tags.Tags {
	rawTags := d.Get("filter_tags").(*schema.Set).List()
	filterTags := make([]tags.Tags, len(rawTags))
	for i, raw := range rawTags {
		rawMap := raw.(map[string]interface{})
		rawValues := rawMap["values"].(*schema.Set)
		values := make([]string, rawValues.Len())
		for i, list := range rawValues.List() {
			values[i] = list.(string)
		}
		filterTags[i] = tags.Tags{
			Key:    rawMap["key"].(string),
			Values: values,
		}
	}
	return filterTags
}
