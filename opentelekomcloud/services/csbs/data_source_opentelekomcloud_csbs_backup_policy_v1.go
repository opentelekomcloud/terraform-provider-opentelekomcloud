package csbs

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/csbs/v1/policies"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCSBSBackupPolicyV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCSBSBackupPolicyV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"common": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"resource": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"scheduled_operation": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"max_backups": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"retention_duration_days": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"permanent": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"trigger_pattern": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operation_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
		},
	}
}

func dataSourceCSBSBackupPolicyV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	policyClient, err := config.CsbsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating CSBSv1 client: %w", err)
	}

	listOpts := policies.ListOpts{
		ID:     d.Id(),
		Name:   d.Get("name").(string),
		Status: d.Get("status").(string),
	}

	refinedPolicies, err := policies.List(policyClient, listOpts)

	if err != nil {
		return fmterr.Errorf("unable to retrieve backup policies: %s", err)
	}

	if len(refinedPolicies) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedPolicies) > 1 {
		return fmterr.Errorf("your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	backupPolicy := refinedPolicies[0]

	log.Printf("[INFO] Retrieved backup policy %s using given filter", backupPolicy.ID)

	d.SetId(backupPolicy.ID)

	if err := d.Set("resource", flattenCSBSPolicyResources(backupPolicy)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("scheduled_operation", flattenCSBSScheduledOperations(backupPolicy)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tags", flattenCSBSPolicyTags(backupPolicy)); err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("name", backupPolicy.Name),
		d.Set("id", backupPolicy.ID),
		d.Set("common", backupPolicy.Parameters.Common),
		d.Set("status", backupPolicy.Status),
		d.Set("description", backupPolicy.Description),
		d.Set("provider_id", backupPolicy.ProviderId),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
