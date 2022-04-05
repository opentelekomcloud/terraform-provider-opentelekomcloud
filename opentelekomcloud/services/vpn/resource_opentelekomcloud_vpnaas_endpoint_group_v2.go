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
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/vpnaas/endpointgroups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceVpnEndpointGroupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpnEndpointGroupV2Create,
		ReadContext:   resourceVpnEndpointGroupV2Read,
		UpdateContext: resourceVpnEndpointGroupV2Update,
		DeleteContext: resourceVpnEndpointGroupV2Delete,
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
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"endpoints": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVpnEndpointGroupV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	var createOpts endpointgroups.CreateOptsBuilder

	endpointType := resourceVpnEndpointGroupV2EndpointType(d.Get("type").(string))
	v := d.Get("endpoints").([]interface{})
	endpoints := make([]string, len(v))
	for i, v := range v {
		endpoints[i] = v.(string)
	}

	createOpts = VpnEndpointGroupCreateOpts{
		endpointgroups.CreateOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			TenantID:    d.Get("tenant_id").(string),
			Endpoints:   endpoints,
			Type:        endpointType,
		},
		common.MapValueSpecs(d),
	}

	log.Printf("[DEBUG] Create group: %#v", createOpts)

	group, err := endpointgroups.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_CREATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForEndpointGroupCreation(networkingClient, group.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] EndpointGroup created: %#v", group)

	d.SetId(group.ID)

	return resourceVpnEndpointGroupV2Read(ctx, d, meta)
}

func resourceVpnEndpointGroupV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Retrieve information about group: %s", d.Id())

	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	group, err := endpointgroups.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "group")
	}

	log.Printf("[DEBUG] Read OpenTelekomCloud Endpoint EndpointGroup %s: %#v", d.Id(), group)

	mErr := multierror.Append(
		d.Set("name", group.Name),
		d.Set("description", group.Description),
		d.Set("tenant_id", group.TenantID),
		d.Set("type", group.Type),
		d.Set("endpoints", group.Endpoints),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVpnEndpointGroupV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	opts := endpointgroups.UpdateOpts{}

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

	log.Printf("[DEBUG] Updating endpoint group with id %s: %#v", d.Id(), opts)

	if hasChange {
		group, err := endpointgroups.Update(networkingClient, d.Id(), opts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"PENDING_UPDATE"},
			Target:     []string{"UPDATED"},
			Refresh:    waitForEndpointGroupUpdate(networkingClient, group.ID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      0,
			MinTimeout: 2 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)

		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] Updated group with id %s", d.Id())
	}

	return resourceVpnEndpointGroupV2Read(ctx, d, meta)
}

func resourceVpnEndpointGroupV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Destroy group: %s", d.Id())

	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	err = endpointgroups.Delete(networkingClient, d.Id()).Err

	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    waitForEndpointGroupDeletion(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	return diag.FromErr(err)
}

func waitForEndpointGroupDeletion(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		group, err := endpointgroups.Get(networkingClient, id).Extract()
		log.Printf("[DEBUG] Got group %s => %#v", id, group)

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] EndpointGroup %s is actually deleted", id)
				return "", "DELETED", nil
			}
			return nil, "", fmt.Errorf("unexpected error: %s", err)
		}

		log.Printf("[DEBUG] EndpointGroup %s deletion is pending", id)
		return group, "DELETING", nil
	}
}

func waitForEndpointGroupCreation(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		group, err := endpointgroups.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_CREATE", nil
		}
		return group, "ACTIVE", nil
	}
}

func waitForEndpointGroupUpdate(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		group, err := endpointgroups.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_UPDATE", nil
		}
		return group, "UPDATED", nil
	}
}

func resourceVpnEndpointGroupV2EndpointType(epType string) endpointgroups.EndpointType {
	switch epType {
	case "subnet":
		return endpointgroups.TypeSubnet
	case "cidr":
		return endpointgroups.TypeCIDR
	case "vlan":
		return endpointgroups.TypeVLAN
	case "router":
		return endpointgroups.TypeRouter
	case "network":
		return endpointgroups.TypeNetwork
	}
	return ""
}
