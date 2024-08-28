package ecs

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/volumeattach"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceComputeVolumeAttachV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeVolumeAttachV2Create,
		ReadContext:   resourceComputeVolumeAttachV2Read,
		DeleteContext: resourceComputeVolumeAttachV2Delete,
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

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"device": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceComputeVolumeAttachV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.ComputeV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateV2Client, err)
	}

	instanceId := d.Get("instance_id").(string)
	volumeId := d.Get("volume_id").(string)

	var device string
	if v, ok := d.GetOk("device"); ok {
		device = v.(string)
	}

	attachOpts := volumeattach.CreateOpts{
		Device:   device,
		VolumeID: volumeId,
	}

	log.Printf("[DEBUG] Creating volume attachment: %#v", attachOpts)

	attachment, err := volumeattach.Create(client, instanceId, attachOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ATTACHING"},
		Target:     []string{"ATTACHED"},
		Refresh:    resourceComputeVolumeAttachV2AttachFunc(client, instanceId, attachment.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      30 * time.Second,
		MinTimeout: 15 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmterr.Errorf("error attaching OpenTelekomCloud volume: %s", err)
	}

	log.Printf("[DEBUG] Created volume attachment: %#v", attachment)

	// Use the instance ID and attachment ID as the resource ID.
	// This is because an attachment cannot be retrieved just by its ID alone.
	id := fmt.Sprintf("%s/%s", instanceId, attachment.ID)

	d.SetId(id)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceComputeVolumeAttachV2Read(clientCtx, d, meta)
}

func resourceComputeVolumeAttachV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.ComputeV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateV2Client, err)
	}

	instanceId, attachmentId, err := ParseComputeVolumeAttachmentId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	attachment, err := volumeattach.Get(client, instanceId, attachmentId).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "compute_volume_attach")
	}

	log.Printf("[DEBUG] Retrieved volume attachment: %#v", attachment)

	mErr := multierror.Append(
		d.Set("instance_id", attachment.ServerID),
		d.Set("volume_id", attachment.VolumeID),
		d.Set("device", attachment.Device),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceComputeVolumeAttachV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.ComputeV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateV2Client, err)
	}

	instanceId, attachmentId, err := ParseComputeVolumeAttachmentId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{""},
		Target:     []string{"DETACHED"},
		Refresh:    resourceComputeVolumeAttachV2DetachFunc(client, instanceId, attachmentId),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      15 * time.Second,
		MinTimeout: 15 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmterr.Errorf("error detaching OpenTelekomCloud volume: %s", err)
	}

	return nil
}

func resourceComputeVolumeAttachV2AttachFunc(
	computeClient *golangsdk.ServiceClient, instanceId, attachmentId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		va, err := volumeattach.Get(computeClient, instanceId, attachmentId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return va, "ATTACHING", nil
			}
			return va, "", err
		}

		return va, "ATTACHED", nil
	}
}

func resourceComputeVolumeAttachV2DetachFunc(
	computeClient *golangsdk.ServiceClient, instanceId, attachmentId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to detach OpenTelekomCloud volume %s from instance %s",
			attachmentId, instanceId)

		va, err := volumeattach.Get(computeClient, instanceId, attachmentId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return va, "DETACHED", nil
			}
			return va, "", err
		}

		err = volumeattach.Delete(computeClient, instanceId, attachmentId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return va, "DETACHED", nil
			}

			if _, ok := err.(golangsdk.ErrDefault400); ok {
				return nil, "", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud Volume Attachment (%s) is still active.", attachmentId)
		return nil, "", nil
	}
}

func ParseComputeVolumeAttachmentId(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		return "", "", fmt.Errorf("unable to determine volume attachment ID")
	}

	instanceId := idParts[0]
	attachmentId := idParts[1]

	return instanceId, attachmentId, nil
}
