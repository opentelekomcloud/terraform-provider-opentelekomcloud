package evs

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/evs/v2/snapshots"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceEvsSnapshotV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvsSnapshotV2Create,
		ReadContext:   resourceEvsSnapshotV2Read,
		UpdateContext: resourceEvsSnapshotV2Update,
		DeleteContext: resourceEvsSnapshotV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metadata": {
				Type:             schema.TypeMap,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: common.SuppressMapDiffs(),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
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

func resourceEvsSnapshotV2Metadata(d *schema.ResourceData) map[string]string {
	metadata := make(map[string]string)
	for k, v := range d.Get("metadata").(map[string]interface{}) {
		metadata[k] = v.(string)
	}
	return metadata
}

func resourceEvsSnapshotV2StateRefreshFunc(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		snapshot, err := snapshots.Get(client, id)
		if err != nil {
			var errDefault404 golangsdk.ErrDefault404
			if errors.As(err, &errDefault404) {
				return snapshot, "deleted", nil
			}
			return nil, "", err
		}

		switch snapshot.Status {
		case "error":
			return snapshot, snapshot.Status, errors.New("there was an error creating the EVS snapshot")
		case "error_deleting":
			return snapshot, snapshot.Status, errors.New("there was an error deleting the EVS snapshot")
		default:
			return snapshot, snapshot.Status, nil
		}
	}
}

func waitForEvsSnapshotV2Status(ctx context.Context, client *golangsdk.ServiceClient, id string, timeout time.Duration,
	pending, target []string) error {
	stateConf := &resource.StateChangeConf{
		Pending:      pending,
		Target:       target,
		Refresh:      resourceEvsSnapshotV2StateRefreshFunc(client, id),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		MinTimeout:   3 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceEvsSnapshotV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.BlockStorageV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	snapshot, err := snapshots.Create(client, snapshots.CreateOpts{
		VolumeID:    d.Get("volume_id").(string),
		Force:       d.Get("force").(bool),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Metadata:    resourceEvsSnapshotV2Metadata(d),
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud EVS volume snapshot: %s", err)
	}
	d.SetId(snapshot.ID)

	if err := waitForEvsSnapshotV2Status(ctx, client, d.Id(), d.Timeout(schema.TimeoutCreate),
		[]string{"creating"}, []string{"available"}); err != nil {
		return diag.Errorf("error waiting for EVS snapshot (%s) creation to available: %s", d.Id(), err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceEvsSnapshotV2Read(clientCtx, d, meta)
}

func resourceEvsSnapshotV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.BlockStorageV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	snapshot, err := snapshots.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving EVS snapshot")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("volume_id", snapshot.VolumeID),
		d.Set("name", snapshot.Name),
		d.Set("metadata", snapshot.Metadata),
		d.Set("description", snapshot.Description),
		d.Set("status", snapshot.Status),
		d.Set("size", snapshot.Size),
		d.Set("created_at", snapshot.CreatedAt.Format(time.RFC3339)),
		d.Set("updated_at", snapshot.UpdatedAt.Format(time.RFC3339)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEvsSnapshotV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.BlockStorageV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	_, err = snapshots.Update(client, d.Id(), snapshots.UpdateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	})
	if err != nil {
		return diag.Errorf("error updating EVS snapshot: %s", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceEvsSnapshotV2Read(clientCtx, d, meta)
}

func resourceEvsSnapshotV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.BlockStorageV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	if err := snapshots.Delete(client, d.Id()); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting EVS snapshot")
	}

	if err := waitForEvsSnapshotV2Status(ctx, client, d.Id(), d.Timeout(schema.TimeoutDelete),
		[]string{"available", "deleting"}, []string{"deleted"}); err != nil {
		return diag.Errorf("error waiting for EVS snapshot (%s) deleted: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}
