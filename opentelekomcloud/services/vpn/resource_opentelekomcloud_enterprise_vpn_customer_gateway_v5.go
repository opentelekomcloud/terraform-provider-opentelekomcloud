package vpn

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	cgw "github.com/opentelekomcloud/gophertelekomcloud/openstack/evpn/v5/customer-gateway"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceEnterpriseCustomerGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvpnCustomerGatewayCreate,
		UpdateContext: resourceEvpnCustomerGatewayUpdate,
		ReadContext:   resourceEvpnCustomerGatewayRead,
		DeleteContext: resourceEvpnCustomerGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"asn": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  65000,
				ForceNew: true,
			},
			"id_value": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"id_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ip",
				ForceNew: true,
			},
			"tags": common.TagsSchema(),
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The create time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceEvpnCustomerGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV5, func() (*golangsdk.ServiceClient, error) {
		return config.EvpnV5Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV5Client, err)
	}

	gatewayTags := d.Get("tags").(map[string]interface{})
	var tagSlice []tags.ResourceTag
	for k, v := range gatewayTags {
		tagSlice = append(tagSlice, tags.ResourceTag{Key: k, Value: v.(string)})
	}
	createOpts := cgw.CreateOpts{
		Name:    d.Get("name").(string),
		BgpAsn:  pointerto.Int(d.Get("asn").(int)),
		Tags:    tagSlice,
		IdType:  d.Get("id_type").(string),
		IdValue: d.Get("id_value").(string),
	}

	n, err := cgw.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud EVPN customer gateway: %w", err)
	}
	d.SetId(n.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV5)
	return resourceEvpnCustomerGatewayRead(clientCtx, d, meta)
}

func resourceEvpnCustomerGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV5, func() (*golangsdk.ServiceClient, error) {
		return config.EvpnV5Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV5Client, err)
	}

	gw, err := cgw.Get(client, d.Id())
	if err != nil {
		return diag.Errorf("error retrieving OpenTelekomCloud EVPN customer gateway (%s): %s", d.Id(), err)
	}

	tagsMap := make(map[string]string)
	for _, tag := range gw.Tags {
		tagsMap[tag.Key] = tag.Value
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", gw.Name),
		d.Set("id_value", gw.IdValue),
		d.Set("asn", gw.BgpAsn),
		d.Set("id_type", gw.IdType),
		d.Set("created_at", gw.CreatedAt),
		d.Set("updated_at", gw.UpdatedAt),
		d.Set("tags", tagsMap),
		d.Set("ip", gw.Ip),
		d.Set("route_mode", gw.RouteMode),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEvpnCustomerGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV5, func() (*golangsdk.ServiceClient, error) {
		return config.EvpnV5Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV5Client, err)
	}

	if d.HasChange("name") {
		_, err := cgw.Update(client, cgw.UpdateOpts{
			GatewayID: d.Id(),
			Name:      d.Get("name").(string),
		})
		if err != nil {
			return diag.Errorf("error updating OpenTelekomCloud EVPN customer gateway: %s", err)
		}
	}

	if d.HasChange("tags") {
		if err = updateTags(client, d, "customer-gateway", d.Id()); err != nil {
			return diag.Errorf("error updating tags of OpenTelekomCloud EVPN customer gateway (%s): %s", d.Id(), err)
		}
	}
	clientCtx := common.CtxWithClient(ctx, client, keyClientV5)
	return resourceEvpnCustomerGatewayRead(clientCtx, d, meta)
}

func resourceEvpnCustomerGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV5, func() (*golangsdk.ServiceClient, error) {
		return config.EvpnV5Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV5Client, err)
	}

	err = cgw.Delete(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting OpenTelekomCloud EVPN customer gateway")
	}

	return nil
}
