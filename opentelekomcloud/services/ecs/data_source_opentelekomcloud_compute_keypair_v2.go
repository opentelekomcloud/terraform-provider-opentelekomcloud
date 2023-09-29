package ecs

import (
	"context"
	"log"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/keypairs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceComputeKeypairV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeKeypairV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceComputeKeypairV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreateV2Client, err)
	}

	var allKeyPairs []keypairs.KeyPair
	if keypairName := d.Get("name").(string); keypairName != "" {
		kp, err := keypairs.Get(client, keypairName).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmterr.Errorf("no keypair found")
			}
			return fmterr.Errorf("unable to retrieve OpenTelekomCloud %s keypair: %w", keypairName, err)
		}

		allKeyPairs = append(allKeyPairs, *kp)
	} else {
		allPages, err := keypairs.List(client).AllPages()
		if err != nil {
			return fmterr.Errorf("Unable to query OpenTelekomCloud keypairs: %w", err)
		}

		allKeyPairs, err = keypairs.ExtractKeyPairs(allPages)
		if err != nil {
			return fmterr.Errorf("Unable to retrieve OpenTelekomCloud keypairs: %w", err)
		}
	}

	var filteredKeyPairs []keypairs.KeyPair
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		r := regexp.MustCompile(nameRegex.(string))
		for _, kp := range allKeyPairs {
			if r.MatchString(kp.Name) {
				filteredKeyPairs = append(filteredKeyPairs, kp)
			}
		}
		allKeyPairs = filteredKeyPairs
	}

	if len(allKeyPairs) < 1 {
		return fmterr.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	keypair := allKeyPairs[0]

	log.Printf("[DEBUG] Retrieved opentelekoncloud_compute_keypair_v2 %s: %#v", keypair.Name, keypair)

	d.SetId(keypair.Name)
	mErr := multierror.Append(
		d.Set("name", keypair.Name),
		d.Set("user_id", keypair.UserID),
		d.Set("fingerprint", keypair.Fingerprint),
		d.Set("public_key", keypair.PublicKey),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
