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
				Computed: true,
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
		VpcID:        d.Get("vpc_id").(string),
		Bandwidth:    d.Get("bandwidth").(int),
		Type:         d.Get("type").(string),
		AdminStateUp: &adminStateUp,
		AZ:           d.Get("az").(string),
		ChargeMode:   d.Get("charge_mode").(string),
		EipType:      d.Get("eip_type").(string),
		VipAddress:   d.Get("vip_address").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	job, err := loadbalancer_elbs.Create(networkingClient, createOpts).ExtractJobResponse()
	if err != nil {
		fmt.Printf("resourceELoadBalancerCreate Error creating LoadBalancer: %s \n", err)

		return err
	}

	//fmt.Printf("job=%+v.\n", job)

	log.Printf("Successfully created loadbalancer %s on subnet %s", createOpts.Name, createOpts.VipSubnetID)
	log.Printf("Waiting for loadbalancer %s to become active", createOpts.Name)

	if err := WaitForJobSuccess(networkingClient, job.URI, loadbalancerActiveTimeoutSeconds); err != nil {
		fmt.Printf("WaitForJobSuccess fails err=%v+  %s to become active", err, createOpts.Name)
		return err
	}

	mlb, err := GetJobEntity(networkingClient, job.URI, "elb")
	log.Printf("LoadBalancer %s is active", createOpts.Name)

	if vid, ok := mlb["id"]; ok {
		fmt.Printf("resourceELoadBalancerCreate: vid=%s.\n", vid)
		if id, ok := vid.(string); ok {
			fmt.Printf("id=%s.\n", id)
			lb, err := loadbalancer_elbs.Get(networkingClient, id).Extract()
			if err != nil {
				fmt.Printf("loadbalancer_elbs Extract Error: %s.\n", err.Error())
				return err
			}
			fmt.Printf("got lb=%+v.\n", lb)
			//return lb, err

			fmt.Printf("@@@@@@@@@@@@@@@@ resourceELoadBalancerCreate  LoadBalancer:  created %v+ %s \n", lb, lb.ID)

			// Once the loadbalancer has been created, apply any requested security groups
			// to the port that was created behind the scenes.
			//?

			// If all has been successful, set the ID on the resource
			d.SetId(lb.ID)

			return resourceELoadBalancerRead(d, meta)
		}
	}
	return nil
}

func resourceELoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV1Client(GetRegion(d, config))
	if err != nil {
		fmt.Printf("@@@@@@@@@@@@@@@@ resourceELoadBalancerRead Error creating OpenTelekomCloud networking client: %s \n", err)

		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	lb, err := loadbalancer_elbs.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		fmt.Printf("@@@@@@@@@@@@@@@@ resourceELoadBalancerRead Extract: %s \n", err)

		return CheckDeleted(d, err, "loadbalancer")
	}

	log.Printf("[DEBUG] Retrieved loadbalancer %s: %#v", d.Id(), lb)

	//fmt.Printf("@@@@@@@@@@@@@@@@ resourceELoadBalancerRead Retrieved loadbalancer %s: %#v \n", d.Id(), lb)

	//?
	d.Set("name", lb.Name)
	d.Set("description", lb.Description)
	d.Set("vip_subnet_id", lb.VipSubnetID)
	d.Set("tenant_id", lb.TenantID)
	d.Set("vip_address", lb.VipAddress)
	d.Set("vpc_id", lb.VpcID)
	d.Set("admin_state_up", lb.AdminStateUp)
	d.Set("az", lb.AZ)
	d.Set("vip_address", lb.VipAddress)
	d.Set("eip_type", lb.EipType)
	d.Set("bandwidth", lb.Bandwidth)
	d.Set("charge_mode", lb.ChargeMode)
	d.Set("region", GetRegion(d, config))

	// Get any security groups on the VIP Port

	return nil
}

func resourceELoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	fmt.Printf("[resourceELoadBalancerUpdate] ########## update loadbalancer %s", d.Id())
	config := meta.(*Config)
	networkingClient, err := config.networkingV1Client(GetRegion(d, config))
	if err != nil {
		fmt.Printf("Error creating OpenTelekomCloud networking client: %s", err)
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	var updateOpts loadbalancer_elbs.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}
	if d.HasChange("bandwidth") {
		updateOpts.Bandwidth = d.Get("bandwidth").(int)
	}
	if d.HasChange("admin_state_up") {
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}
	timeout := d.Timeout(schema.TimeoutUpdate)

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

	return resourceELoadBalancerRead(d, meta)
}

func resourceELoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	log.Printf("[DEBUG] Deleting loadbalancer %s", d.Id())
	fmt.Printf("[resourceELoadBalancerDelete] ##########nn Deleting loadbalancer %s", d.Id())
	timeout := d.Timeout(schema.TimeoutDelete)
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = loadbalancer_elbs.Delete(networkingClient, d.Id()).ExtractErr()
		if err != nil {
			fmt.Printf("[resourceELoadBalancerDelete] ##########nn Deleting err %+v \n", err)

			return checkForRetryableError(err)
		}
		return nil
	})

	// Wait for LoadBalancer to become delete
	/*pending := []string{"PENDING_UPDATE", "PENDING_DELETE", "ACTIVE"}
	fmt.Printf("[resourceELoadBalancerDelete] ##########nn waiting loadbalancer %s", d.Id())

	err = waitForELBLoadBalancer(networkingClient, d.Id(), "DELETED", pending, timeout)
	if err != nil {
		fmt.Printf("[resourceELoadBalancerDelete] ##########nn Deleting waitForELBLoadBalancer err %+v \n", err)
		return err
	}

	fmt.Printf("[resourceELoadBalancerDelete] ##########nn done loadbalancer %s", d.Id()) */
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
