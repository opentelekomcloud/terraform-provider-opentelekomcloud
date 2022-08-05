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
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/blockstorage/extensions/volumeactions"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/blockstorage/v2/volumes"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/volumeattach"
	volumes_v3 "github.com/opentelekomcloud/gophertelekomcloud/openstack/evs/v3/volumes"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceBlockStorageVolumeV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBlockStorageVolumeV2Create,
		ReadContext:   resourceBlockStorageVolumeV2Read,
		UpdateContext: resourceBlockStorageVolumeV2Update,
		DeleteContext: resourceBlockStorageVolumeV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.ForceNewIfChange("size", isDownScale),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"metadata": {
				Type:             schema.TypeMap,
				Optional:         true,
				ForceNew:         false,
				Computed:         true,
				DiffSuppressFunc: suppressMetadataComputed,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_vol_id": {
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
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"device_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "VBD",
				ValidateFunc: validation.StringInSlice([]string{"VBD", "SCSI"}, true),
			},
			"consistency_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_replica": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
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
				Set: resourceVolumeV2AttachmentHash,
			},
			"cascade": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"wwn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceContainerMetadataV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("metadata").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func resourceContainerTags(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("tags").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func resourceBlockStorageVolumeV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.BlockStorageV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	metadata := resourceContainerMetadataV2(d)
	if d.Get("device_type").(string) == "SCSI" {
		metadata["hw:passthrough"] = "true"
	}

	createOpts := &volumes.CreateOpts{
		AvailabilityZone:   d.Get("availability_zone").(string),
		ConsistencyGroupID: d.Get("consistency_group_id").(string),
		Description:        d.Get("description").(string),
		ImageID:            d.Get("image_id").(string),
		Metadata:           metadata,
		Name:               d.Get("name").(string),
		Size:               d.Get("size").(int),
		SnapshotID:         d.Get("snapshot_id").(string),
		SourceReplica:      d.Get("source_replica").(string),
		SourceVolID:        d.Get("source_vol_id").(string),
		VolumeType:         d.Get("volume_type").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := volumes.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud volume: %s", err)
	}
	log.Printf("[INFO] Volume ID: %s", v.ID)

	// Wait for the volume to become available.
	log.Printf(
		"[DEBUG] Waiting for volume (%s) to become available",
		v.ID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"downloading", "creating"},
		Target:       []string{"available"},
		Refresh:      VolumeV2StateRefreshFunc(client, v.ID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		MinTimeout:   3 * time.Second,
		PollInterval: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf(
			"Error waiting for volume (%s) to become ready: %s",
			v.ID, err)
	}
	_, err = resourceEVSTagV2Create(ctx, d, meta, "volumes", v.ID, resourceContainerTags(d))
	if err != nil {
		return fmterr.Errorf("error creating tags for volume (%s): %s", v.ID, err)
	}

	// Store the ID now
	d.SetId(v.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceBlockStorageVolumeV2Read(clientCtx, d, meta)
}

func resourceBlockStorageVolumeV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.BlockStorageV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	v, err := volumes.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "volume")
	}

	log.Printf("[DEBUG] Retrieved volume %s: %+v", d.Id(), v)
	mErr := multierror.Append(nil,
		d.Set("size", v.Size),
		d.Set("description", v.Description),
		d.Set("availability_zone", v.AvailabilityZone),
		d.Set("name", v.Name),
		d.Set("snapshot_id", v.SnapshotID),
		d.Set("source_vol_id", v.SourceVolID),
		d.Set("volume_type", v.VolumeType),
		d.Set("region", config.GetRegion(d)),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}
	// NOTE: This tries to remove system metadata.
	md := make(map[string]string)
	var sys_keys = [1]string{"hw:passthrough"}

OUTER:
	for key, val := range v.Metadata {
		for i := range sys_keys {
			if key == sys_keys[i] {
				continue OUTER
			}
		}
		md[key] = val
	}

	if err := d.Set("metadata", md); err != nil {
		return fmterr.Errorf("[DEBUG] Error saving metadata to state for OpenTelekomCloud block storage (%s): %s", d.Id(), err)
	}

	attachments := make([]map[string]interface{}, len(v.Attachments))
	for i, attachment := range v.Attachments {
		attachments[i] = make(map[string]interface{})
		attachments[i]["id"] = attachment.ID
		attachments[i]["instance_id"] = attachment.ServerID
		attachments[i]["device"] = attachment.Device
		log.Printf("[DEBUG] attachment: %v", attachment)
	}
	if err := d.Set("attachment", attachments); err != nil {
		return fmterr.Errorf("[DEBUG] Error saving attachment to state for OpenTelekomCloud block storage (%s): %s", d.Id(), err)
	}
	taglist, err := resourceEVSTagV2Get(d, meta, "volumes", v.ID)
	if err != nil {
		return fmterr.Errorf("error fetching tags for volume (%s): %s", v.ID, err)
	}
	mErr = multierror.Append(mErr, d.Set("tags", taglist.Tags))

	// This is useful for import
	if d.Get("device_type").(string) == "" {
		mErr = multierror.Append(mErr, d.Set("device_type", "VBD"))
	}

	if d.Get("device_type").(string) == "SCSI" {
		blockStorageClientV3, err := config.BlockStorageV3Client(config.GetRegion(d))
		if err != nil {
			return fmterr.Errorf("error creating OpenTelekomCloud block storage client: %s", err)
		}

		v, err := volumes_v3.Get(blockStorageClientV3, d.Id()).Extract()
		if err != nil {
			return common.CheckDeletedDiag(d, err, "volume")
		}

		log.Printf("[DEBUG] Retrieved volume %s: %+v", d.Id(), v)

		mErr = multierror.Append(mErr, d.Set("wwn", v.WWN))
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceBlockStorageVolumeV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.BlockStorageV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	updateOpts := volumes.UpdateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	if d.HasChange("metadata") {
		updateOpts.Metadata = resourceVolumeMetadataV2(d)
	}

	_, err = volumes.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud volume: %s", err)
	}
	if d.HasChange("tags") {
		_, err = resourceEVSTagV2Create(ctx, d, meta, "volumes", d.Id(), resourceContainerTags(d))
		if err != nil {
			return fmterr.Errorf("error creating tags for the volume: %w", err)
		}
	}

	if d.HasChange("size") {
		if err := extendSize(d, client); err != nil {
			return diag.FromErr(err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"extending"},
			Target:       []string{"available", "in-use"},
			Refresh:      VolumeV2StateRefreshFunc(client, d.Id()),
			Timeout:      d.Timeout(schema.TimeoutDelete),
			Delay:        10 * time.Second,
			MinTimeout:   3 * time.Second,
			PollInterval: 2 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmterr.Errorf("error waiting for volume (%s) to become ready after resize: %s", d.Id(), err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceBlockStorageVolumeV2Read(clientCtx, d, meta)
}

func extendSize(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	oldSize, newSize := d.GetChange("size")
	newSizeInt := newSize.(int)
	if oldSize.(int) > newSizeInt {
		return fmt.Errorf("can't decrease volume size. This is internal provider error, this point should never be reached")
	}
	opts := volumeactions.ExtendSizeOpts{NewSize: newSizeInt}
	if err := volumeactions.ExtendSize(client, d.Id(), opts).ExtractErr(); err != nil {
		return fmt.Errorf("failed to extend disk size: %s", err)
	}
	return nil
}

func resourceBlockStorageVolumeV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.BlockStorageV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	v, err := volumes.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "volume")
	}

	// make sure this volume is detached from all instances before deleting
	if len(v.Attachments) > 0 {
		log.Printf("[DEBUG] detaching volumes")
		if computeClient, err := config.ComputeV2Client(config.GetRegion(d)); err != nil {
			return diag.FromErr(err)
		} else {
			for _, volumeAttachment := range v.Attachments {
				log.Printf("[DEBUG] Attachment: %v", volumeAttachment)
				if err := volumeattach.Delete(computeClient, volumeAttachment.ServerID, volumeAttachment.ID).ExtractErr(); err != nil {
					return diag.FromErr(err)
				}
			}

			stateConf := &resource.StateChangeConf{
				Pending:      []string{"in-use", "attaching", "detaching"},
				Target:       []string{"available"},
				Refresh:      VolumeV2StateRefreshFunc(client, d.Id()),
				Timeout:      d.Timeout(schema.TimeoutDelete),
				Delay:        10 * time.Second,
				MinTimeout:   3 * time.Second,
				PollInterval: 2 * time.Second,
			}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return fmterr.Errorf(
					"Error waiting for volume (%s) to become available: %s",
					d.Id(), err)
			}
		}
	}

	// The snapshots associated with the disk are deleted together with the EVS disk if cascade value is true
	deleteOpts := volumes.DeleteOpts{
		Cascade: d.Get("cascade").(bool),
	}
	// It's possible that this volume was used as a boot device and is currently
	// in a "deleting" state from when the instance was terminated.
	// If this is true, just move on. It'll eventually delete.
	if v.Status != "deleting" {
		if err := volumes.Delete(client, d.Id(), deleteOpts).ExtractErr(); err != nil {
			return common.CheckDeletedDiag(d, err, "volume")
		}
	}

	// Wait for the volume to delete before moving on.
	log.Printf("[DEBUG] Waiting for volume (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"deleting", "downloading", "available"},
		Target:       []string{"deleted"},
		Refresh:      VolumeV2StateRefreshFunc(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		MinTimeout:   3 * time.Second,
		PollInterval: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf(
			"Error waiting for volume (%s) to delete: %s",
			d.Id(), err)
	}

	d.SetId("")
	return nil
}

func resourceVolumeMetadataV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("metadata").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

// VolumeV2StateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an OpenTelekomCloud volume.
func VolumeV2StateRefreshFunc(client *golangsdk.ServiceClient, volumeID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := volumes.Get(client, volumeID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "deleted", nil
			}
			return nil, "", err
		}

		if v.Status == "error" {
			return v, v.Status, fmt.Errorf("there was an error creating the volume. " +
				"Please check with your cloud admin or check the Block Storage " +
				"API logs to see why this error occurred")
		}

		return v, v.Status, nil
	}
}

func resourceVolumeV2AttachmentHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if m["instance_id"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["instance_id"].(string)))
	}
	return hashcode.String(buf.String())
}

var suppressedVolumeFields = []string{
	"metadata.policy",
	"metadata.backupId",
}

func suppressMetadataComputed(k, old, new string, _ *schema.ResourceData) bool {
	if common.StringInSlice(k, suppressedVolumeFields) {
		return true
	}
	return old == new
}

func isDownScale(_ context.Context, old, new, _ interface{}) bool {
	return old.(int) > new.(int)
}
