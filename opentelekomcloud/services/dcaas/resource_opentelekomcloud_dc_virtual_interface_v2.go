package dcaas

import (
	"context"
	"log"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	virtual_interface "github.com/opentelekomcloud/gophertelekomcloud/openstack/dcaas/v2/virtual-interface"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceVirtualInterfaceV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVirtualInterfaceV2Create,
		ReadContext:   resourceVirtualInterfaceV2Read,
		UpdateContext: resourceVirtualInterfaceV2Update,
		DeleteContext: resourceVirtualInterfaceV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"direct_connect_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"virtual_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
				),
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"private",
				}, false),
			},
			"route_mode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"static", "bgp",
				}, false),
			},
			"vlan": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 3999),
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"remote_ep_group_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"service_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"local_gateway_v4_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"remote_gateway_v4_ip"},
			},
			"remote_gateway_v4_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"asn": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntNotInSlice([]int{64512}),
			},
			"bgp_md5": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enable_bfd": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable_nqa": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"lag_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVirtualInterfaceV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DCaaSV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	opts := virtual_interface.CreateOpts{
		TenantID:          d.Get("project_id").(string),
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		DirectConnectID:   d.Get("direct_connect_id").(string),
		VgwID:             d.Get("virtual_gateway_id").(string),
		Type:              d.Get("type").(string),
		ServiceType:       d.Get("service_type").(string),
		VLAN:              d.Get("vlan").(int),
		Bandwidth:         d.Get("bandwidth").(int),
		LocalGatewayV4IP:  d.Get("local_gateway_v4_ip").(string),
		RemoteGatewayV4IP: d.Get("remote_gateway_v4_ip").(string),
		RouteMode:         d.Get("route_mode").(string),
		BGPASN:            d.Get("asn").(int),
		BGPMD5:            d.Get("bgp_md5").(string),
		RemoteEPGroupID:   d.Get("remote_ep_group_id").(string),
	}
	vi, err := virtual_interface.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating opentelekomcloud virtual interface: %s", err)
	}
	d.SetId(vi.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceVirtualInterfaceV2Read(clientCtx, d, meta)
}

func resourceVirtualInterfaceV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DCaaSV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	vi, err := virtual_interface.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "virtual interface")
	}
	log.Printf("[DEBUG] The virtual interface response is: %#v", vi)

	mErr := multierror.Append(nil,
		d.Set("vgw_id", vi.VgwID),
		d.Set("type", vi.Type),
		d.Set("route_mode", vi.RouteMode),
		d.Set("vlan", vi.VLAN),
		d.Set("bandwidth", vi.Bandwidth),
		d.Set("remote_ep_group_id", vi.RemoteEPGroupID),
		d.Set("name", vi.Name),
		d.Set("description", vi.Description),
		d.Set("direct_connect_id", vi.DirectConnectID),
		d.Set("service_type", vi.ServiceType),
		d.Set("local_gateway_v4_ip", vi.LocalGatewayV4IP),
		d.Set("remote_gateway_v4_ip", vi.RemoteGatewayV4IP),
		d.Set("asn", vi.BGPASN),
		d.Set("bgp_md5", vi.BGPMD5),
		d.Set("enable_bfd", vi.EnableBFD),
		d.Set("enable_nqa", vi.EnableNQA),
		d.Set("lag_id", vi.LagID),
		d.Set("status", vi.Status),
		d.Set("created_at", vi.CreateTime),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving opentelekomcloud virtual interface fields: %s", err)
	}
	return nil
}

func resourceVirtualInterfaceV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DCaaSV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	if d.HasChanges("name", "description", "bandwidth", "remote_ep_group_id") {
		opts := virtual_interface.UpdateOpts{
			Name:                  d.Get("name").(string),
			Description:           d.Get("description").(string),
			Bandwidth:             d.Get("bandwidth").(int),
			RemoteEndpointGroupId: d.Get("remote_ep_group_id").(string),
		}

		err := virtual_interface.Update(client, d.Id(), opts)
		if err != nil {
			return diag.Errorf("error closing network detection of the opentelekomcloud virtual interface (%s): %s", d.Id(), err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceVirtualInterfaceV2Read(clientCtx, d, meta)
}

func resourceVirtualInterfaceV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DCaaSV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	err = virtual_interface.Delete(client, d.Id())
	if err != nil {
		return diag.Errorf("error deleting opentelekomcloud virtual interface (%s): %s", d.Id(), err)
	}

	return nil
}
