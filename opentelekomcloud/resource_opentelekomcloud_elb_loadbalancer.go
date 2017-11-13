package opentelekomcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/loadbalancer_elbs"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

func resourceELoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceELoadBalancerCreate,
		Read:   resourceELoadBalancerRead,
		Update: resourceELoadBalancerUpdate,
		Delete: resourceELoadBalancerDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"bandwidth": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"admin_state_up": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},

			"vip_subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"az": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"charge_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"eip_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"vip_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceELoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	createOpts := loadbalancer_elbs.CreateOpts{
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		VipSubnetID:  d.Get("vip_subnet_id").(string),
		Tenant_ID:    d.Get("tenant_id").(string),
		VpcID:        d.Get("vpd_id").(string),
		Bandwidth:    d.Get("bandwidth").(int),
		Type:         d.Get("type").(string),
		AdminStateUp: &adminStateUp,
		AZ:           d.Get("az").(string),
		ChargeMode:   d.Get("charge_mode").(string),
		EipType:      d.Get("eip_type").(string),
		VipAddress:   d.Get("vip_address").(string),
		TenantID:     d.Get("tenantId").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	lb, err := loadbalancer_elbs.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating LoadBalancer: %s", err)
	}

	// Wait for LoadBalancer to become active before continuing
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForELBLoadBalancer(networkingClient, lb.ID, "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	// Once the loadbalancer has been created, apply any requested security groups
	// to the port that was created behind the scenes.
	//?

	// If all has been successful, set the ID on the resource
	d.SetId(lb.ID)

	return resourceELoadBalancerRead(d, meta)
}

func resourceELoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	lb, err := loadbalancer_elbs.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "loadbalancer")
	}

	log.Printf("[DEBUG] Retrieved loadbalancer %s: %#v", d.Id(), lb)

	//?
	d.Set("region", GetRegion(d, config))

	// Get any security groups on the VIP Port

	return nil
}

func resourceELoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	var updateOpts loadbalancer_elbs.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}
	if d.HasChange("admin_state_up") {
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	// Wait for LoadBalancer to become active before continuing
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForLBV2LoadBalancer(networkingClient, d.Id(), "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating loadbalancer %s with options: %#v", d.Id(), updateOpts)
	err = resource.Retry(timeout, func() *resource.RetryError {
		_, err = loadbalancer_elbs.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	// Wait for LoadBalancer to become active before continuing
	err = waitForELBLoadBalancer(networkingClient, d.Id(), "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	// Security Groups get updated separately
	if d.HasChange("security_group_ids") {
		vipPortID := d.Get("vip_port_id").(string)
		if err := resourceLoadBalancerV2SecurityGroups(networkingClient, vipPortID, d); err != nil {
			return err
		}
	}

	return resourceLoadBalancerV2Read(d, meta)
}

func resourceELoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	log.Printf("[DEBUG] Deleting loadbalancer %s", d.Id())
	timeout := d.Timeout(schema.TimeoutDelete)
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = loadbalancer_elbs.Delete(networkingClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	// Wait for LoadBalancer to become delete
	pending := []string{"PENDING_UPDATE", "PENDING_DELETE", "ACTIVE"}
	err = waitForELBLoadBalancer(networkingClient, d.Id(), "DELETED", pending, timeout)
	if err != nil {
		return err
	}

	return nil
}

func resourceELoadBalancerSecurityGroups(networkingClient *gophercloud.ServiceClient, vipPortID string, d *schema.ResourceData) error {
	if vipPortID != "" {
		if _, ok := d.GetOk("security_group_ids"); ok {
			updateOpts := ports.UpdateOpts{
				SecurityGroups: resourcePortSecurityGroupsV2(d),
			}

			log.Printf("[DEBUG] Adding security groups to loadbalancer "+
				"VIP Port %s: %#v", vipPortID, updateOpts)

			_, err := ports.Update(networkingClient, vipPortID, updateOpts).Extract()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
