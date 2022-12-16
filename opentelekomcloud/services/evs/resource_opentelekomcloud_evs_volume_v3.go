package evs

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	cinderV3 "github.com/opentelekomcloud/gophertelekomcloud/openstack/blockstorage/v3/volumes"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/evs/v3/volumes"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceEvsStorageVolumeV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvsVolumeV3Create,
		ReadContext:   resourceEvsVolumeV3Read,
		UpdateContext: resourceEvsVolumeV3Update,
		DeleteContext: resourceBlockStorageVolumeV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		CustomizeDiff: common.MultipleCustomizeDiffs(
			common.ValidateVolumeType("volume_type"),
			customdiff.ForceNewIfChange("size", isDownScale),
		),

		Schema: map[string]*schema.Schema{
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"device_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "VBD",
				ValidateFunc: validation.StringInSlice([]string{"VBD", "SCSI"}, true),
			},
			"tags": common.TagsSchema(),
			"attachment": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: resourceVolumeAttachmentHash,
			},
			"multiattach": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"kms_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cascade": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"wwn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceEvsVolumeV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.BlockStorageV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	if !common.HasFilledOpt(d, "backup_id") && !common.HasFilledOpt(d, "size") {
		return fmterr.Errorf("missing required argument: 'size' is required, but no definition was found")
	}

	createOpts := &volumes.CreateOpts{
		BackupID:         d.Get("backup_id").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
		Description:      d.Get("description").(string),
		Size:             d.Get("size").(int),
		Name:             d.Get("name").(string),
		SnapshotID:       d.Get("snapshot_id").(string),
		ImageRef:         d.Get("image_id").(string),
		VolumeType:       d.Get("volume_type").(string),
		Multiattach:      d.Get("multiattach").(bool),
	}
	m := make(map[string]string)
	if v, ok := d.GetOk("kms_id"); ok {
		m["__system__cmkid"] = v.(string)
		m["__system__encrypted"] = "1"
	}
	if d.Get("device_type").(string) == "SCSI" {
		m["hw:passthrough"] = "true"
	}
	if len(m) > 0 {
		createOpts.Metadata = m
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := volumes.Create(client, createOpts).ExtractJobResponse()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud EVS volume: %s", err)
	}
	log.Printf("[INFO] Volume Job ID: %s", v.JobID)

	// Wait for the volume to become available.
	log.Printf("[DEBUG] Waiting for volume to become available")
	err = volumes.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutCreate)/time.Second), v.JobID)
	if err != nil {
		return diag.FromErr(err)
	}

	entity, err := volumes.GetJobEntity(client, v.JobID, "volume_id")
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := entity.(string); ok {
		log.Printf("[INFO] Volume ID: %s", id)
		// Store the ID now
		d.SetId(id)

		// set tags
		tagRaw := d.Get("tags").(map[string]interface{})
		if len(tagRaw) > 0 {
			tagList := common.ExpandResourceTags(tagRaw)
			if err := tags.Create(client, "os-vendor-volumes", id, tagList).ExtractErr(); err != nil {
				return fmterr.Errorf("error setting tags for EVSv3 Volume: %w", err)
			}
		}
		clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
		return resourceEvsVolumeV3Read(clientCtx, d, meta)
	}
	return fmterr.Errorf("unexpected conversion error in resourceEvsVolumeV3Create")
}

func resourceEvsVolumeV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.BlockStorageV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	v, err := volumes.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "volume")
	}

	log.Printf("[DEBUG] Retrieved volume %s: %+v", d.Id(), v)

	mErr := multierror.Append(
		d.Set("size", v.Size),
		d.Set("description", v.Description),
		d.Set("availability_zone", v.AvailabilityZone),
		d.Set("name", v.Name),
		d.Set("snapshot_id", v.SnapshotID),
		d.Set("volume_type", v.VolumeType),
		d.Set("wwn", v.WWN),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	// save tags
	resourceTags, err := tags.Get(client, "os-vendor-volumes", d.Id()).Extract()
	if err != nil {
		return fmterr.Errorf("error fetching OpenTelekomCloud SFS File System tags: %s", err)
	}
	tagMap := common.TagsToMap(resourceTags)
	if err := d.Set("tags", tagMap); err != nil {
		return fmterr.Errorf("error saving tags for OpenTelekomCloud EVSv3 Volume: %s", err)
	}

	// set attachments
	attachments := make([]map[string]interface{}, len(v.Attachments))
	for i, attachment := range v.Attachments {
		attachments[i] = make(map[string]interface{})
		attachments[i]["id"] = attachment.ID
		attachments[i]["instance_id"] = attachment.ServerID
		attachments[i]["device"] = attachment.Device
		log.Printf("[DEBUG] attachment: %v", attachment)
	}
	if err := d.Set("attachment", attachments); err != nil {
		return fmterr.Errorf("[DEBUG] Error saving attachment to state for OpenTelekomCloud evs storage (%s): %s", d.Id(), err)
	}

	return nil
}

// using OpenStack Cinder API v2 to update volume resource
func resourceEvsVolumeV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.BlockStorageV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	updateOpts := cinderV3.UpdateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	_, err = cinderV3.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud volume: %s", err)
	}

	// update tags
	if d.HasChange("tags") {
		if err := common.UpdateResourceTags(client, d, "os-vendor-volumes", d.Id()); err != nil {
			return fmterr.Errorf("error updating tags for EVSv3 Volume %s: %w", d.Id(), err)
		}
	}

	if d.HasChange("size") {
		if err := extendSize(d, client); err != nil {
			return diag.FromErr(err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"extending"},
			Target:     []string{"available", "in-use"},
			Refresh:    volumeV3StateRefreshFunc(client, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutDelete),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmterr.Errorf("error waiting for volume (%s) to become ready after resize: %s", d.Id(), err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceEvsVolumeV3Read(clientCtx, d, meta)
}

func resourceVolumeAttachmentHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if m["instance_id"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["instance_id"].(string)))
	}
	return hashcode.String(buf.String())
}

func volumeV3StateRefreshFunc(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := cinderV3.Get(client, id).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "deleted", nil
			}
			return nil, "", err
		}
		if v.Status == "error" {
			return v, v.Status, fmt.Errorf("volume is in the error state")
		}
		return v, v.Status, nil
	}
}
