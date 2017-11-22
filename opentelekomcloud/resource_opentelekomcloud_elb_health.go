package opentelekomcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/healthcheck"
)

func resourceHealth() *schema.Resource {
	return &schema.Resource{
		Create: resourceHealthCreate,
		Read:   resourceHealthRead,
		Update: resourceHealthUpdate,
		Delete: resourceHealthDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delay": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"max_retries": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"url_path": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"http_method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"expected_codes": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"admin_state_up": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},

			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceHealthCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	//adminStateUp := d.Get("admin_state_up").(bool)
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

	timeout := d.Timeout(schema.TimeoutCreate)
	// Wait for LoadBalancer to become active before continuing
	lbID := d.Get("loadbalancer_id").(string)
	err = waitForELBLoadBalancer(networkingClient, lbID, "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	log.Printf("[DEBUG] Attempting to create monitor")
	var health *healthcheck.Health
	err = resource.Retry(timeout, func() *resource.RetryError {
		health, err = healthcheck.Create(networkingClient, createOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	//err = waitForLBV2viaPool(networkingClient, poolID, "ACTIVE", timeout)
	//d.SetId(monitor.ID)

	return resourceHealthRead(d, meta)
}

func resourceHealthRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.otcV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	health, err := healthcheck.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "health")
	}

	log.Printf("[DEBUG] Retrieved health %s: %#v", d.Id(), health)

	d.Set("id", health.ID)
	//d.Set("type", health.Type)
	//d.Set("delay", health.Delay)
	//d.Set("timeout", health.Timeout)
	//d.Set("max_retries", health.MaxRetries)
	//d.Set("url_path", health.URLPath)
	//d.Set("http_method", health.HTTPMethod)
	//d.Set("expected_codes", health.ExpectedCodes)
	//d.Set("admin_state_up", health.AdminStateUp)
	//d.Set("name", health.Name)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceHealthUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.otcV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
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
	timeout := d.Timeout(schema.TimeoutUpdate)
	poolID := d.Get("pool_id").(string)
	err = waitForLBV2viaPool(networkingClient, poolID, "ACTIVE", timeout)
	if err != nil {
		return err
	}

	err = resource.Retry(timeout, func() *resource.RetryError {
		_, err = healthcheck.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Unable to update monitor %s: %s", d.Id(), err)
	}

	// Wait for LB to become active before continuing

	return resourceHealthRead(d, meta)
}

func resourceHealthDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.otcV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	log.Printf("[DEBUG] Deleting health %s", d.Id())
	timeout := d.Timeout(schema.TimeoutUpdate)
	//err = waitForLBV2viaPool(networkingClient, poolID, "ACTIVE", timeout)

	err = resource.Retry(timeout, func() *resource.RetryError {
		err = healthcheck.Delete(networkingClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Unable to delete health %s: %s", d.Id(), err)
	}

	//err = waitForLBV2viaPool(networkingClient, poolID, "ACTIVE", timeout)

	return nil
}
