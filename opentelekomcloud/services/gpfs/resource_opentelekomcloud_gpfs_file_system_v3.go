package gpfs

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/gpfs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceGPFSFileSystemV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGPFSFileSystemV3Create,
		ReadContext:   resourceGPFSFileSystemV3Read,
		DeleteContext: resourceGPFSFileSystemV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(10 * time.Minute),
			Delete:  schema.DefaultTimeout(10 * time.Minute),
			Default: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"redundancy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "3az",
				ValidateFunc: validation.StringInSlice([]string{"3az"}, false),
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func isGPFSService404(err error) bool {
	if _, ok := err.(golangsdk.ErrDefault404); ok {
		return true
	}
	serviceError, ok := err.(gpfs.Error)
	return ok && strings.HasPrefix(serviceError.Status, "404")
}

func resourceGPFSFileSystemV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NewGPFSClient(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud GPFS client: %w", err)
	}

	name := d.Get("name").(string)
	redundancy := d.Get("redundancy").(string)
	opts := gpfs.CreateFSInput{
		BucketLocation: gpfs.BucketLocation{
			Location: config.GetRegion(d),
		},
		FSName:     name,
		Redundancy: redundancy,
		BucketType: "SFS",
	}
	if _, err = client.CreateFS(&opts); err != nil {
		return diag.Errorf("error creating GPFS file system (%s): %s", name, err)
	}

	d.SetId(name)
	return resourceGPFSFileSystemV3Read(ctx, d, meta)
}

func resourceGPFSFileSystemV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NewGPFSClient(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud GPFS client: %w", err)
	}

	var bucket gpfs.Bucket
	found := false
	output, err := client.ListFS(&gpfs.ListFSInput{BucketType: "SFS"})
	if err != nil {
		return diag.Errorf("error fetching GPFS file systems: %s", err)
	}
	for _, b := range output.Buckets {
		if b.Name == d.Id() {
			bucket = b
			found = true
		}
	}

	if !found {
		return diag.Errorf("error retrieving GPFS file system (%s): %s", d.Id(), err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", bucket.Name),
		d.Set("redundancy", "3az"),
		d.Set("creation_date", bucket.CreationDate.Format(time.RFC3339)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving GPFS file system (%s) fields: %s", d.Id(), err)
	}
	return nil
}

func resourceGPFSFileSystemV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NewGPFSClient(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud GPFS client: %w", err)
	}

	log.Printf("[DEBUG] deleting OBS Bucket: %s", d.Id())
	if _, err = client.DeleteFS(d.Id()); err != nil {
		if isGPFSService404(err) {
			return nil
		}
		return diag.Errorf("error deleting GPFS file system (%s): %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}
