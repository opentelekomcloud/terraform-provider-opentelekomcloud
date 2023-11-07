package dcaas

import (
	"context"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	virtual_gateway "github.com/opentelekomcloud/gophertelekomcloud/openstack/dcaas/v2/virtual-gateway"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceVirtualGatewayV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVirtualGatewayV2Create,
		ReadContext:   resourceVirtualGatewayV2Read,
		UpdateContext: resourceVirtualGatewayV2Update,
		DeleteContext: resourceVirtualGatewayV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"local_ep_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[^<>]*$`),
						"The angle brackets (< and >) are not allowed."),
					validation.StringLenBetween(0, 128),
				),
			},
			"asn": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"device_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"redundant_device_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVirtualGatewayV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DCaaSV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	opts := virtual_gateway.CreateOpts{
		VpcId:                d.Get("vpc_id").(string),
		LocalEndpointGroupId: d.Get("local_ep_group_id").(string),
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		BgpAsn:               d.Get("asn").(int),
		DeviceId:             d.Get("device_id").(string),
		RedundantDeviceId:    d.Get("redundant_device_id").(string),
	}
	vg, err := virtual_gateway.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating opentelekomcloud virtual gateway: %s", err)
	}
	d.SetId(vg.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceVirtualGatewayV2Read(clientCtx, d, meta)
}

func resourceVirtualGatewayV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DCaaSV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	vg, err := virtual_gateway.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "virtual gateway")
	}

	mErr := multierror.Append(nil,
		d.Set("vpc_id", vg.VPCID),
		d.Set("local_ep_group_id", vg.LocalEPGroupID),
		d.Set("name", vg.Name),
		d.Set("description", vg.Description),
		d.Set("asn", vg.BGPASN),
		d.Set("device_id", vg.DeviceID),
		d.Set("redundant_device_id", vg.RedundantDeviceID),
		d.Set("project_id", vg.TenantID),
		d.Set("status", vg.Status),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving opentelekomcloud virtual gateway fields: %s", err)
	}
	return nil
}

func resourceVirtualGatewayV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DCaaSV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	opts := virtual_gateway.UpdateOpts{
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		LocalEndpointGroupId: d.Get("local_ep_group_id").(string),
	}
	err = virtual_gateway.Update(client, d.Id(), opts)
	if err != nil {
		return diag.Errorf("error updating opentelekomcloud virtual gateway (%s): %s", d.Id(), err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceVirtualGatewayV2Read(clientCtx, d, meta)
}

func resourceVirtualGatewayV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DCaaSV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	err = virtual_gateway.Delete(client, d.Id())
	if err != nil {
		return diag.Errorf("error deleting opentelekomcloud virtual gateway (%s): %s", d.Id(), err)
	}

	return nil
}
