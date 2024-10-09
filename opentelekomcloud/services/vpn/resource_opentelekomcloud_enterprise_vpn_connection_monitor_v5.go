package vpn

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	connection_monitoring "github.com/opentelekomcloud/gophertelekomcloud/openstack/evpn/v5/connection-monitoring"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceEnterpriseConnectionMonitor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvpnConnectionMonitorCreate,
		ReadContext:   resourceEvpnConnectionMonitorRead,
		DeleteContext: resourceEvpnConnectionMonitorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"connection_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_ip": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceEvpnConnectionMonitorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV5, func() (*golangsdk.ServiceClient, error) {
		return config.EvpnV5Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV5Client, err)
	}

	opts := connection_monitoring.CreateOpts{
		ConnectionID: d.Get("connection_id").(string),
	}
	n, err := connection_monitoring.Create(client, opts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud EVPN customer gateway: %w", err)
	}
	d.SetId(n.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV5)
	return resourceEvpnConnectionMonitorRead(clientCtx, d, meta)
}

func resourceEvpnConnectionMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV5, func() (*golangsdk.ServiceClient, error) {
		return config.EvpnV5Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV5Client, err)
	}

	cm, err := connection_monitoring.Get(client, d.Id())
	if err != nil {
		return diag.Errorf("error retrieving OpenTelekomCloud EVPN connection monitor (%s): %s", d.Id(), err)
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("destination_ip", cm.DestinationIp),
		d.Set("source_ip", cm.SourceIp),
		d.Set("status", cm.Status),
		d.Set("connection_id", cm.ConnectionId),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEvpnConnectionMonitorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV5, func() (*golangsdk.ServiceClient, error) {
		return config.EvpnV5Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV5Client, err)
	}

	err = connection_monitoring.Delete(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting OpenTelekomCloud EVPN connection monitor")
	}

	return nil
}
