package dis

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	_ "github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dis/v2/apps"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDisAppV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDisAppV2Create,
		ReadContext:   resourceDisAppV2Read,
		DeleteContext: resourceDisAppV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"app_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 200),
					validation.StringMatch(
						regexp.MustCompile(`^[A-Za-z0-9\-_]+$`),
						"Only letters, digits, underscores (_), and hyphens (-) are allowed.",
					),
				),
			},
			"created": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"commit_checkpoint_stream_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"partition_consuming_states": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sequence_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"latest_offset": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"earliest_offset": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"checkpoint_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceDisAppV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DisV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	opts := apps.CreateAppOpts{
		AppName: d.Get("app_name").(string),
	}

	log.Printf("[DEBUG] Creating new App: %s", opts.AppName)
	err = apps.CreateApp(client, opts)
	if err != nil {
		return fmterr.Errorf("error creating DIS streams: %s", err)
	}

	d.SetId(opts.AppName)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceDisAppV2Read(clientCtx, d, meta)
}

func resourceDisAppV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DisV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	app, err := apps.GetApp(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error getting app details: %s", err)
	}

	mErr := multierror.Append(
		d.Set("app_name", app.AppName),
		d.Set("app_id", app.AppId),
		d.Set("created", app.CreateTime),
		d.Set("commit_checkpoint_stream_names", app.CommitCheckPointStreamNames),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	var result []map[string]interface{}
	for _, v := range app.CommitCheckPointStreamNames {
		status, errStat := apps.GetAppStatus(client, apps.GetAppStatusOpts{
			AppName:        app.AppName,
			StreamName:     v,
			CheckpointType: checkpointType,
		})
		if errStat != nil {
			return fmterr.Errorf("error getting app status: %s", errStat)
		}

		for _, partition := range status.PartitionConsumingStates {
			result = append(result, map[string]interface{}{
				"id":              partition.PartitionId,
				"status":          partition.Status,
				"sequence_number": partition.SequenceNumber,
				"latest_offset":   partition.LatestOffset,
				"earliest_offset": partition.EarliestOffset,
				"checkpoint_type": partition.CheckpointType,
			})
		}
	}
	err = d.Set("partition_consuming_states", result)
	if err != nil {
		fmterr.Errorf("error setting DISv2 consuming states: %s", err)
	}

	return nil
}

func resourceDisAppV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DisV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	name := d.Id()
	err = apps.DeleteApp(client, name)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud DIS app: %s", err)
	}

	d.SetId("")
	return nil
}
