package cbr

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cbr/v3/backups"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCBRBackupsIdsV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCBRBackupsIdsV3Read,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"checkpoint_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"image_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vault_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_az": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceCBRBackupsIdsV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	cbrClient, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("unable to create opentelekomcloud CBR client : %s", err)
	}

	listOpts := backups.ListOpts{
		ID:           d.Get("id").(string),
		CheckpointID: d.Get("checkpoint_id").(string),
		ImageType:    d.Get("image_type").(string),
		Name:         d.Get("name").(string),
		VaultID:      d.Get("vault_id").(string),
		ParentID:     d.Get("parent_id").(string),
		ResourceAZ:   d.Get("resource_az").(string),
		ResourceID:   d.Get("resource_id").(string),
		ResourceName: d.Get("resource_name").(string),
		ResourceType: d.Get("resource_type").(string),
		Status:       d.Get("status").(string),
	}
	extractedBackups, err := backups.List(cbrClient, listOpts)

	if err != nil {
		return fmterr.Errorf("unable to retrieve Backups: %s", err)
	}

	if len(extractedBackups) < 1 {
		return fmterr.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	allBackups := make([]string, 0)
	for _, singleBackup := range extractedBackups {
		allBackups = append(allBackups, singleBackup.ID)
	}
	d.SetId("Filter")
	mErr := multierror.Append(
		d.Set("ids", allBackups),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
