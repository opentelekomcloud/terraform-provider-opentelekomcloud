package cfw

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	cfwips "github.com/opentelekomcloud/gophertelekomcloud/openstack/cfw/v1/ips"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCfwIpsProtectionV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCFWIpsProtectionV1Create,
		ReadContext:   resourceCFWIpsProtectionV1Read,
		DeleteContext: resourceCFWIpsProtectionV1Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"object_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"ips_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{2}),
			},
			"feature_status": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"mode": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 3}),
			},
			"basic_defense_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ips_switch_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ips_protection_mode_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCFWIpsProtectionV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	objectId := d.Get("object_id").(string)
	virtPatchStatus := InterfaceToIntPtr(d.Get("feature_status"))

	featureOpts := cfwips.SetFeatureStatusOpts{
		ObjectID: objectId,
		Status:   virtPatchStatus,
		IpsType:  d.Get("ips_type").(int),
	}

	err = cfwips.SetIPSFeatureStatus(client, featureOpts)
	if err != nil {
		return fmterr.Errorf("error enabling/disabling IPS feature status: %w", err)
	}

	mode := InterfaceToIntPtr(d.Get("mode"))
	protectionModeOpts := cfwips.SetProtectionModeOpts{
		ObjectID: objectId,
		Mode:     mode,
	}
	err = cfwips.SetProtectionMode(client, protectionModeOpts)
	if err != nil {
		return fmterr.Errorf("error setting IPS protection mode: %w", err)
	}
	log.Printf("[DEBUG] IPS protection  set to mode %d and virtual patching status: %d", *mode, *virtPatchStatus)

	d.SetId(objectId)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceCFWIpsProtectionV1Read(clientCtx, d, meta)
}

func resourceCFWIpsProtectionV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	ipsFeatResp, err := cfwips.GetIPSFeatureStatus(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error fetching IPS protection feature status: %w", err)
	}
	ipsProtectionResp, err := cfwips.GetProtectionMode(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error fetching IPS protection mode: %w", err)
	}

	mErr := multierror.Append(nil,
		d.Set("ips_switch_id", ipsFeatResp.ID),
		d.Set("basic_defense_status", ipsFeatResp.BasicDefenseStatus),
		d.Set("feature_status", ipsFeatResp.VirtualPatchesStatus),
		d.Set("mode", ipsProtectionResp.Mode),
		d.Set("ips_protection_mode_id", ipsProtectionResp.ID),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCFWIpsProtectionV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Deleting action resource is not supported. The action resource is only removed from the state the task remains in the cloud.",
		},
	}
}
