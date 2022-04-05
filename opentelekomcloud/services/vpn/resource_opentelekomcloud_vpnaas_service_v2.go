package vpn

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/vpnaas/services"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceVpnServiceV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpnServiceV2Create,
		ReadContext:   resourceVpnServiceV2Read,
		UpdateContext: resourceVpnServiceV2Update,
		DeleteContext: resourceVpnServiceV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				Computed: false,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_v6_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_v4_ip": {
				Type:     schema.TypeString,
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

func resourceVpnServiceV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	var createOpts services.CreateOptsBuilder

	adminStateUp := d.Get("admin_state_up").(bool)
	createOpts = VpnServiceCreateOpts{
		services.CreateOpts{
			Name:         d.Get("name").(string),
			Description:  d.Get("description").(string),
			AdminStateUp: &adminStateUp,
			TenantID:     d.Get("tenant_id").(string),
			SubnetID:     d.Get("subnet_id").(string),
			RouterID:     d.Get("router_id").(string),
		},
		common.MapValueSpecs(d),
	}

	log.Printf("[DEBUG] Create service: %#v", createOpts)

	service, err := services.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"NOT_CREATED"},
		Target:     []string{"PENDING_CREATE"},
		Refresh:    waitForServiceCreation(networkingClient, service.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Service created: %#v", service)

	d.SetId(service.ID)

	return resourceVpnServiceV2Read(ctx, d, meta)
}

func resourceVpnServiceV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Retrieve information about service: %s", d.Id())

	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	service, err := services.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "service")
	}

	log.Printf("[DEBUG] Read OpenTelekomCloud Service %s: %#v", d.Id(), service)

	mErr := multierror.Append(
		d.Set("name", service.Name),
		d.Set("description", service.Description),
		d.Set("subnet_id", service.SubnetID),
		d.Set("admin_state_up", service.AdminStateUp),
		d.Set("tenant_id", service.TenantID),
		d.Set("router_id", service.RouterID),
		d.Set("status", service.Status),
		d.Set("external_v6_ip", service.ExternalV6IP),
		d.Set("external_v4_ip", service.ExternalV4IP),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVpnServiceV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	opts := services.UpdateOpts{}

	var hasChange bool

	if d.HasChange("name") {
		name := d.Get("name").(string)
		opts.Name = &name
		hasChange = true
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		opts.Description = &description
		hasChange = true
	}

	if d.HasChange("admin_state_up") {
		adminStateUp := d.Get("admin_state_up").(bool)
		opts.AdminStateUp = &adminStateUp
		hasChange = true
	}

	log.Printf("[DEBUG] Updating service with id %s: %#v", d.Id(), opts)

	if hasChange {
		service, err := services.Update(networkingClient, d.Id(), opts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"PENDING_UPDATE"},
			Target:     []string{"UPDATED"},
			Refresh:    waitForServiceUpdate(networkingClient, service.ID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      0,
			MinTimeout: 2 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)

		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] Updated service with id %s", d.Id())
	}

	return resourceVpnServiceV2Read(ctx, d, meta)
}

func resourceVpnServiceV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Destroy service: %s", d.Id())

	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	err = services.Delete(networkingClient, d.Id()).Err

	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    waitForServiceDeletion(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)

	return diag.FromErr(err)
}

func waitForServiceDeletion(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		serv, err := services.Get(networkingClient, id).Extract()
		log.Printf("[DEBUG] Got service %s => %#v", id, serv)

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Service %s is actually deleted", id)
				return "", "DELETED", nil
			}
			return nil, "", fmt.Errorf("unexpected error: %s", err)
		}

		log.Printf("[DEBUG] Service %s deletion is pending", id)
		return serv, "DELETING", nil
	}
}

func waitForServiceCreation(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		service, err := services.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "NOT_CREATED", nil
		}
		return service, "PENDING_CREATE", nil
	}
}

func waitForServiceUpdate(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		service, err := services.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_UPDATE", nil
		}
		return service, "UPDATED", nil
	}
}
