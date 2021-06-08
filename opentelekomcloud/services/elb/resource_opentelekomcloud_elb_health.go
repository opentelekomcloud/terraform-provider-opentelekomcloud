package elb

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/elbaas/healthcheck"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceHealth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHealthCreate,
		ReadContext:   resourceHealthRead,
		UpdateContext: resourceHealthUpdate,
		DeleteContext: resourceHealthDelete,

		DeprecationMessage: classicLBDeprecated,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"healthcheck_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"HTTP", "TCP"}, false),
			},
			"healthcheck_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"healthcheck_connect_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"healthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"unhealthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"healthcheck_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 50),
			},
			"healthcheck_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 5),
			},
		},
	}
}

func resourceHealthCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ElbV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	// adminStateUp := d.Get("admin_state_up").(bool)
	createOpts := healthcheck.CreateOpts{
		ListenerID:             d.Get("listener_id").(string),
		HealthcheckProtocol:    d.Get("healthcheck_protocol").(string),
		HealthcheckUri:         d.Get("healthcheck_uri").(string),
		HealthcheckConnectPort: d.Get("healthcheck_connect_port").(int),
		HealthyThreshold:       d.Get("healthy_threshold").(int),
		UnhealthyThreshold:     d.Get("unhealthy_threshold").(int),
		HealthcheckTimeout:     d.Get("healthcheck_timeout").(int),
		HealthcheckInterval:    d.Get("healthcheck_interval").(int),
	}

	health, err := healthcheck.Create(client, createOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(health.ID)

	log.Printf("Successfully created healthcheck %s.", health.ID)

	return resourceHealthRead(ctx, d, meta)
}

func resourceHealthRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.ElbV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	health, err := healthcheck.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "health"))
	}

	log.Printf("[DEBUG] Retrieved health %s: %+v", d.Id(), health)

	d.Set("listener_id", health.ListenerID)
	d.Set("healthcheck_protocol", health.HealthcheckProtocol)
	d.Set("healthcheck_uri", health.HealthcheckUri)
	d.Set("healtcheck_connect_port", health.HealthcheckConnectPort)
	d.Set("healthy_threshold", health.HealthyThreshold)
	d.Set("unhealthy_threshold", health.UnhealthyThreshold)
	d.Set("healthcheck_timeout", health.HealthcheckTimeout)
	d.Set("healthcheck_interval", health.HealthcheckInterval)

	d.Set("region", config.GetRegion(d))

	return nil
}

func resourceHealthUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.ElbV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	var updateOpts healthcheck.UpdateOpts
	if d.HasChange("healthcheck_protocol") {
		updateOpts.HealthcheckProtocol = d.Get("healthcheck_protocol").(string)
	}
	if d.HasChange("healthcheck_uri") {
		updateOpts.HealthcheckUri = d.Get("healthcheck_uri").(string)
	}
	if d.HasChange("healthcheck_connect_port") {
		updateOpts.HealthyThreshold = d.Get("healthcheck_connect_port").(int)
	}
	if d.HasChange("healthy_threshold") {
		updateOpts.HealthyThreshold = d.Get("healthy_threshold").(int)
	}
	if d.HasChange("unhealthy_threshold") {
		updateOpts.UnhealthyThreshold = d.Get("unhealthy_threshold").(int)
	}
	if d.HasChange("healthcheck_timeout") {
		updateOpts.HealthcheckTimeout = d.Get("healthcheck_timeout").(int)
	}
	if d.HasChange("healthcheck_interval") {
		updateOpts.HealthcheckInterval = d.Get("healthcheck_interval").(int)
	}

	log.Printf("[DEBUG] Updating health %s with options: %#v", d.Id(), updateOpts)

	_, err = healthcheck.Update(networkingClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceHealthRead(ctx, d, meta)
}

func resourceHealthDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ElbV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	id := d.Id()
	log.Printf("[DEBUG] Deleting health %s", id)

	if err := healthcheck.Delete(client, id).ExtractErr(); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("Successfully deleted health %s", id)
	return nil
}
