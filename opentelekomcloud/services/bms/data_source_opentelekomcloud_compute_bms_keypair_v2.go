package bms

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/bms/v2/keypairs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceBMSKeyPairV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBMSKeyPairV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBMSKeyPairV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	bmsClient, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekom bms client: %s", err)
	}

	listOpts := keypairs.ListOpts{
		Name: d.Get("name").(string),
	}

	refinedKeypairs, err := keypairs.List(bmsClient, listOpts)
	if err != nil {
		return fmterr.Errorf("unable to retrieve keypairs: %s", err)
	}

	if len(refinedKeypairs) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedKeypairs) > 1 {
		return fmterr.Errorf("your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Keypairs := refinedKeypairs[0]

	log.Printf("[INFO] Retrieved Keypairs using given filter %s: %+v", Keypairs.Name, Keypairs)
	d.SetId(Keypairs.Name)

	mErr := multierror.Append(
		d.Set("public_key", Keypairs.PublicKey),
		d.Set("fingerprint", Keypairs.Fingerprint),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
