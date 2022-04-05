package vpc

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/layer3/routers"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceNetworkingRouterV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingRouterV2Create,
		ReadContext:   resourceNetworkingRouterV2Read,
		UpdateContext: resourceNetworkingRouterV2Update,
		DeleteContext: resourceNetworkingRouterV2Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"distributed": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"external_gateway": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: common.SuppressExternalGateway,
			},
			"enable_snat": {
				Type:         schema.TypeBool,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"external_gateway"},
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkingRouterV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	createOpts := RouterCreateOpts{
		routers.CreateOpts{
			Name:     d.Get("name").(string),
			TenantID: d.Get("tenant_id").(string),
		},
		common.MapValueSpecs(d),
	}

	if asuRaw, ok := d.GetOk("admin_state_up"); ok {
		asu := asuRaw.(bool)
		createOpts.AdminStateUp = &asu
	}

	if dRaw, ok := d.GetOk("distributed"); ok {
		d := dRaw.(bool)
		createOpts.Distributed = &d
	}

	externalGateway := d.Get("external_gateway").(string)
	snat := d.Get("enable_snat").(bool)
	if externalGateway != "" {
		createOpts.GatewayInfo = &routers.GatewayInfo{
			NetworkID:  externalGateway,
			EnableSNAT: &snat,
		}
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	router, err := routers.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Neutron router: %s", err)
	}
	log.Printf("[INFO] Router ID: %s", router.ID)

	log.Printf("[DEBUG] Waiting for OpenTelekomCloud Neutron Router (%s) to become available", router.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD", "PENDING_CREATE", "PENDING_UPDATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    resourceNetworkingRouterV2StateRefreshFunc(client, router.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting Neutron Router to become available: %w", err)
	}

	d.SetId(router.ID)

	return resourceNetworkingRouterV2Read(ctx, d, meta)
}

func resourceNetworkingRouterV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	router, err := routers.Get(client, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Neutron Router: %w", err)
	}

	log.Printf("[DEBUG] Retrieved Router %s: %+v", d.Id(), router)
	mErr := multierror.Append(
		d.Set("name", router.Name),
		d.Set("admin_state_up", router.AdminStateUp),
		d.Set("distributed", router.Distributed),
		d.Set("tenant_id", router.TenantID),
		d.Set("external_gateway", router.GatewayInfo.NetworkID),
		d.Set("enable_snat", router.GatewayInfo.EnableSNAT),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetworkingRouterV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	var updateOpts routers.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("admin_state_up") {
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	// Gateway settings
	var updateGatewaySettings bool
	var externalGateway string
	gatewayInfo := routers.GatewayInfo{}

	if d.HasChange("external_gateway") {
		updateGatewaySettings = true

		externalGateway = d.Get("external_gateway").(string)
		gatewayInfo.NetworkID = externalGateway
	}

	if d.HasChange("enable_snat") {
		updateGatewaySettings = true

		enableSNAT := d.Get("enable_snat").(bool)
		gatewayInfo.EnableSNAT = &enableSNAT
	}

	if updateGatewaySettings {
		updateOpts.GatewayInfo = &gatewayInfo
	}

	log.Printf("[DEBUG] Updating Router %s with options: %+v", d.Id(), updateOpts)

	_, err = routers.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud Neutron Router: %w", err)
	}

	return resourceNetworkingRouterV2Read(ctx, d, meta)
}

func resourceNetworkingRouterV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	if err := routers.Delete(client, d.Id()).ExtractErr(); err != nil {
		return common.CheckDeletedDiag(d, err, "networking_router_v2")
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    resourceNetworkingRouterV2StateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud Neutron Router: %w", err)
	}

	d.SetId("")
	return nil
}

func resourceNetworkingRouterV2StateRefreshFunc(client *golangsdk.ServiceClient, routerID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := routers.Get(client, routerID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return n, "DELETED", nil
			}

			return n, "", err
		}

		return n, n.Status, nil
	}
}
