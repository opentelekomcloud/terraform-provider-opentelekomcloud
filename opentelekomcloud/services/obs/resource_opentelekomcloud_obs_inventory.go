package obs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/obs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceObsBucketInventory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceObsBucketInventoryPut,
		ReadContext:   resourceObsBucketInventoryRead,
		UpdateContext: resourceObsBucketInventoryPut,
		DeleteContext: resourceObsBucketInventoryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceOBSInventoryImportState,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"configuration_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"frequency": {
				Type:     schema.TypeString,
				Required: true,
			},
			"destination": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"format": {
							Type:     schema.TypeString,
							Required: true,
						},
						"bucket": {
							Type:     schema.TypeString,
							Required: true,
						},
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"included_object_versions": {
				Type:     schema.TypeString,
				Required: true,
			},
			"filter_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceObsBucketInventoryPut(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NewObjectStorageClient(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OBS client: %s", err)
	}

	configId := d.Get("configuration_id").(string)
	bucket := d.Get("bucket").(string)

	log.Printf("[DEBUG] OBS bucket: %s, put inventory: %s", bucket, configId)

	destination := d.Get("destination.0").(map[string]interface{})

	params := &obs.SetBucketInventoryInput{
		Bucket:            bucket,
		InventoryConfigId: configId,
		BucketInventoryConfiguration: obs.BucketInventoryConfiguration{
			Id:        configId,
			IsEnabled: d.Get("is_enabled").(bool),
			Schedule: obs.InventorySchedule{
				Frequency: d.Get("frequency").(string),
			},
			Destination: obs.InventoryDestination{
				Bucket: destination["bucket"].(string),
				Format: destination["format"].(string),
				Prefix: destination["prefix"].(string),
			},
			IncludedObjectVersions: d.Get("included_object_versions").(string),
		},
	}

	if filter, ok := d.GetOk("filter_prefix"); ok {
		params.Filter = obs.InventoryFilter{
			Prefix: filter.(string),
		}
	}

	_, err = client.SetBucketInventory(params)
	if err != nil {
		return fmterr.Errorf("error putting OBS inventory: %s", err)
	}

	d.SetId(configId)

	return resourceObsBucketInventoryRead(ctx, d, meta)
}

func resourceObsBucketInventoryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NewObjectStorageClient(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OBS client: %s", err)
	}

	log.Printf("[DEBUG] OBS bucket inventory, read for bucket: %s", d.Id())
	inputGet := obs.GetBucketInventoryInput{
		BucketName:        d.Get("bucket").(string),
		InventoryConfigId: d.Id(),
	}
	inv, err := client.GetBucketInventory(inputGet)

	if err != nil {
		return fmterr.Errorf("error getting bucket inventory: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("is_enabled", inv.IsEnabled),
		d.Set("frequency", inv.Schedule.Frequency),
		d.Set("destination", flattenDestinationBucketInventory(inv.Destination)),
		d.Set("included_object_versions", inv.IncludedObjectVersions),
		d.Set("filter_prefix", inv.Filter.Prefix),
	)

	if inv.Filter.Prefix != "" {
		mErr = multierror.Append(mErr, d.Set("filter_prefix", inv.Filter.Prefix))
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting OBS bucket replication fields: %s", err)
	}

	return nil
}

func resourceObsBucketInventoryDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NewObjectStorageClient(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OBS client: %s", err)
	}

	inputDelete := obs.DeleteBucketInventoryInput{
		Bucket:            d.Get("bucket").(string),
		InventoryConfigId: d.Id(),
	}

	log.Printf("[DEBUG] OBS bucket: %s, delete inventory", inputDelete.Bucket)
	_, err = client.DeleteBucketInventory(&inputDelete)

	if err != nil {
		return fmterr.Errorf("error deleting OBS inventory: %s", err)
	}
	return nil
}

func flattenDestinationBucketInventory(inv obs.InventoryDestination) (destination []map[string]string) {
	dest := make(map[string]string)
	dest["format"] = inv.Format
	dest["bucket"] = inv.Bucket
	dest["prefix"] = inv.Prefix
	destination = append(destination, dest)
	return
}

func resourceOBSInventoryImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<bucket>/<config_id>', but '%s'", importedId)
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("bucket", parts[0])
}
