package bms

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/bms/v2/tags"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceBMSTagsV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBMSTagsV2Create,
		ReadContext:   resourceBMSTagsV2Read,
		DeleteContext: resourceBMSTagsV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceTagsV2(d *schema.ResourceData) []string {
	rawTAGS := d.Get("tags").(*schema.Set)
	tagList := make([]string, rawTAGS.Len())
	for i, raw := range rawTAGS.List() {
		tagList[i] = raw.(string)
	}
	return tagList
}

func resourceBMSTagsV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	bmsClient, err := config.ComputeV2Client(config.GetRegion(d))

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud bms client: %s", err)
	}

	createOpts := tags.CreateOpts{
		Tag: resourceTagsV2(d),
	}

	_, err = tags.Create(bmsClient, d.Get("server_id").(string), createOpts)

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Tags: %s", err)
	}
	d.SetId(d.Get("server_id").(string))

	log.Printf("[INFO] Server ID: %s", d.Id())

	return resourceBMSTagsV2Read(ctx, d, meta)
}

func resourceBMSTagsV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	bmsClient, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud bms client: %s", err)
	}

	n, err := tags.Get(bmsClient, d.Id())
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud tags: %s", err)
	}

	mErr := multierror.Append(
		d.Set("tags", n),
		d.Set("region", config.GetRegion(d)),
		d.Set("server_id", d.Id()),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceBMSTagsV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	bmsClient, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud bms client: %s", err)
	}

	err = tags.Delete(bmsClient, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud tags: %s", err)
	}

	d.SetId("")
	return nil
}
