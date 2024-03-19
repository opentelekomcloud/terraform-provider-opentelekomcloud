package sdrs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/sdrs/v1/attachreplication"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/sdrs/v1/protectedinstances"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceSDRSReplicationAttachV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSDRSReplicationAttachV1Create,
		ReadContext:   resourceSDRSReplicationAttachV1Read,
		DeleteContext: resourceSDRSReplicationAttachV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceReplicationAttachImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"replication_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"device": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSDRSReplicationAttachV1Create(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, sdrsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.SdrsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	instanceID := d.Get("instance_id").(string)
	replicationID := d.Get("replication_id").(string)
	attachOpts := attachreplication.CreateOpts{
		ReplicationID: replicationID,
		Device:        d.Get("device").(string),
	}

	n, err := attachreplication.Create(client, instanceID, attachOpts).ExtractJobResponse()
	if err != nil {
		return diag.Errorf("error creating SDRS replication attach: %s", err)
	}

	createTimeoutSec := int(d.Timeout(schema.TimeoutCreate).Seconds())
	if err = attachreplication.WaitForJobSuccess(client, createTimeoutSec, n.JobID); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(formatAttachId(instanceID, replicationID))
	return resourceSDRSReplicationAttachV1Read(ctx, d, meta)
}

func resourceSDRSReplicationAttachV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, sdrsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.SdrsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	instanceID := d.Get("instance_id").(string)
	replicationID := d.Get("replication_id").(string)
	n, err := protectedinstances.Get(client, instanceID).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	attachment, err := flattenReplicationAttach(n, replicationID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SDRS replication attach")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("device", attachment.Device),
		d.Set("replication_id", attachment.Replication),
		d.Set("status", n.Status),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenReplicationAttach(instance *protectedinstances.Instance,
	replicationID string) (*protectedinstances.Attachment, error) {
	for _, attach := range instance.Attachment {
		if attach.Replication == replicationID {
			// find the target attachment
			return &attach, nil
		}
	}
	return nil, golangsdk.ErrDefault404{}
}

func resourceSDRSReplicationAttachV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, sdrsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.SdrsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	instanceID := d.Get("instance_id").(string)
	replicationID := d.Get("replication_id").(string)
	n, err := attachreplication.Delete(client, instanceID, replicationID).ExtractJobResponse()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteTimeoutSec := int(d.Timeout(schema.TimeoutDelete).Seconds())
	if err := attachreplication.WaitForJobSuccess(client, deleteTimeoutSec, n.JobID); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func formatAttachId(instanceID string, replicationID string) string {
	return fmt.Sprintf("%s/%s", instanceID, replicationID)
}

func extractAttachId(resourceID string) (instanceID, replicationID string, err error) {
	rgs := strings.Split(resourceID, "/")
	if len(rgs) != 2 {
		err = fmt.Errorf("invalid format specified for replication attach id," +
			" must be <protected_instance_id>/<replication_id>")
		return
	}

	instanceID = rgs[0]
	replicationID = rgs[1]
	return
}

func resourceReplicationAttachImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	instanceID, replicationID, err := extractAttachId(d.Id())
	if err != nil {
		return nil, err
	}
	mErr := multierror.Append(
		nil,
		d.Set("instance_id", instanceID),
		d.Set("replication_id", replicationID),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import replication attach id, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
