package v3

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/log"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceLtsLogV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLtsLogV3Create,
		ReadContext:   resourceLtsLogV3Read,
		DeleteContext: resourceLtsLogV3Delete,
		UpdateContext: resourceLtsLogV3Update,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_stream_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLtsLogV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return diag.FromErr(err)
	}

	l, err := log.Create(client, log.CreateOpts{
		LoadbalancerId: d.Get("loadbalancer_id").(string),
		LogGroupId:     d.Get("log_group_id").(string),
		LogStreamId:    d.Get("log_stream_id").(string),
	})
	if err != nil {
		return fmterr.Errorf("error creating LoadBalancerV3 LogTank: %w", err)
	}

	d.SetId(l.Logtank.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceLtsLogV3Read(clientCtx, d, meta)
}

func resourceLtsLogV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return diag.FromErr(err)
	}

	l, err := log.Get(client, d.Id())
	if err != nil {
		return diag.Errorf("error retrieving LoadBalancerV3 LogTank %s", d.Id())
	}

	mErr := multierror.Append(nil,
		d.Set("loadbalancer_id", l.Logtank.LoadbalancerId),
		d.Set("log_group_id", l.Logtank.LogGroupId),
		d.Set("log_stream_id", l.Logtank.LogStreamId),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Elb LogTank fields: %s", err)
	}
	return nil
}

func resourceLtsLogV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return diag.FromErr(err)
	}

	updateOpts := log.UpdateOpts{
		LogGroupId:  d.Get("log_group_id").(string),
		LogStreamId: d.Get("log_stream_id").(string),
	}
	_, err = log.Update(client, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("unable to update LoadBalancerV3 LogTank %s: %s", d.Id(), err)
	}
	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceLtsLogV3Read(clientCtx, d, meta)
}

func resourceLtsLogV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return diag.FromErr(err)
	}

	err = log.Delete(client, d.Id())
	if err != nil {
		return diag.Errorf("unable to delete LoadBalancerV3 LogTank %s: %s", d.Id(), err)
	}
	return nil
}
