package vpc

import (
	"context"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/bandwidths"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceBandwidthV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBandwidthV2Create,
		ReadContext:   resourceBandwidthV2Read,
		UpdateContext: resourceBandwidthV2Update,
		DeleteContext: resourceBandwidthV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
					validation.StringMatch(
						regexp.MustCompile(`^[\w-.]+$`),
						"The value is a string that can contain letters, digits, underscores (_), hyphens (-), and periods (.).",
					),
				),
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceBandwidthV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	opts := bandwidths.CreateOpts{
		Name: d.Get("name").(string),
		Size: d.Get("size").(int),
	}
	bandwidth, err := bandwidths.Create(client, opts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating bandwidth: %w", err)
	}
	d.SetId(bandwidth.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceBandwidthV2Read(clientCtx, d, meta)
}

func resourceBandwidthV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	bandwidth, err := bandwidths.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error reading bandwidth")
	}

	mErr := multierror.Append(
		d.Set("name", bandwidth.Name),
		d.Set("size", bandwidth.Size),
		d.Set("status", bandwidth.Status),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting bandwidth fields: %w", err)
	}

	return nil
}

func resourceBandwidthV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	opts := bandwidths.UpdateOpts{}
	if d.HasChange("name") {
		opts.Name = d.Get("name").(string)
	}
	if d.HasChange("size") {
		opts.Size = d.Get("size").(int)
	}
	_, err = bandwidths.Update(client, d.Id(), opts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating bandwidth")
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceBandwidthV2Read(clientCtx, d, meta)
}

func resourceBandwidthV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	if err := bandwidths.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmterr.Errorf("error deleting bandwidth: %w", err)
	}

	return nil
}
