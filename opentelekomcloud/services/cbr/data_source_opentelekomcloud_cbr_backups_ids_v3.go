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
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vault_id": {
				Required: true,
				Type:     schema.TypeString,
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
		VaultID: d.Get("vault_id").(string),
	}
	pages, err := backups.List(cbrClient, listOpts).AllPages()

	if err != nil {
		return fmterr.Errorf("unable to retrieve Backups: %s", err)
	}

	extractedBackups, err := backups.ExtractBackups(pages)
	if err != nil {
		return fmterr.Errorf("unable to retrieve backups: %s", err)
	}

	if len(extractedBackups) < 1 {
		return fmterr.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	allBackups := make([]string, 0)
	for _, singleBackup := range extractedBackups {
		allBackups = append(allBackups, singleBackup.ID)
	}
	vaultID := d.Get("vault_id").(string)
	d.SetId(vaultID)
	mErr := multierror.Append(
		d.Set("ids", allBackups),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
