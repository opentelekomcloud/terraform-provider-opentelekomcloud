package sdrs

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/sdrs/v1/replications"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceSDRSReplicationPairV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSDRSReplicationPairV1Create,
		ReadContext:   resourceSDRSReplicationPairV1Read,
		UpdateContext: resourceSDRSReplicationPairV1Update,
		DeleteContext: resourceSDRSReplicationPairV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"delete_target_volume": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"replication_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fault_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_volume_id": {
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

func resourceSDRSReplicationPairV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, sdrsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.SdrsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	createOpts := replications.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		GroupID:     d.Get("group_id").(string),
		VolumeID:    d.Get("volume_id").(string),
	}

	n, err := replications.Create(client, createOpts).ExtractJobResponse()
	if err != nil {
		return diag.Errorf("error creating SDRS replication pair: %s", err)
	}

	createTimeoutSec := int(d.Timeout(schema.TimeoutCreate).Seconds())
	if err = replications.WaitForJobSuccess(client, createTimeoutSec, n.JobID); err != nil {
		return diag.FromErr(err)
	}

	replicationPairID, err := replications.GetJobEntity(client, n.JobID, "replication_pair_id")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(replicationPairID.(string))
	return resourceSDRSReplicationPairV1Read(ctx, d, meta)
}

func resourceSDRSReplicationPairV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, sdrsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.SdrsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	n, err := replications.Get(client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	volumes := strings.Split(n.VolumeIDs, ",")
	if len(volumes) != 2 {
		return diag.Errorf("error retrieving volumes of replication pair: Invalid format. "+
			"except retrieving 2 volumes, but got %d", len(volumes))
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", n.Name),
		d.Set("group_id", n.GroupID),
		d.Set("volume_id", volumes[0]),
		d.Set("description", n.Description),
		d.Set("replication_model", n.ReplicaModel),
		d.Set("fault_level", n.FaultLevel),
		d.Set("status", n.Status),
		d.Set("target_volume_id", volumes[1]),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSDRSReplicationPairV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, sdrsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.SdrsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	if d.HasChange("name") {
		updateOpts := replications.UpdateOpts{
			Name: d.Get("name").(string),
		}
		_, err = replications.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating SDRS replication pair, %s", err)
		}
	}
	return resourceSDRSReplicationPairV1Read(ctx, d, meta)
}

func resourceSDRSReplicationPairV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, sdrsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.SdrsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}

	deleteOpts := replications.DeleteOpts{
		GroupID:      d.Get("group_id").(string),
		DeleteVolume: d.Get("delete_target_volume").(bool),
	}
	n, err := replications.Delete(client, d.Id(), deleteOpts).ExtractJobResponse()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteTimeoutSec := int(d.Timeout(schema.TimeoutDelete).Seconds())
	if err := replications.WaitForJobSuccess(client, deleteTimeoutSec, n.JobID); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
