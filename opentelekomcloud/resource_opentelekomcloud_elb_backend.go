package opentelekomcloud

import (
	"fmt"
	"log"
	"time"

	// "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/backendmember"
)

func resourceBackend() *schema.Resource {
	return &schema.Resource{
		Create: resourceBackendCreate,
		Read:   resourceBackendRead,
		Delete: resourceBackendDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"listener_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"server_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 1 {
						errors = append(errors, fmt.Errorf(
							"Only numbers greater than 0 are supported values for 'weight'"))
					}
					return
				},
			},

			"admin_state_up": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceBackendCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	//adminStateUp := d.Get("admin_state_up").(bool)
	createOpts := backendmember.CreateOpts{
		ListenerId: d.Get("listener_id").(string),
		ServerId:   d.Get("server_id").(string),
		Address:    d.Get("address").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	// Wait for backend  to become active before continuing
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForELBBackend(networkingClient, createOpts.ServerId, "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	// Wait for LB to become ACTIVE again
	// Wait for LoadBalancer to become active again before continuing
	lbID := d.Get("loadbalancer_id").(string)
	err = waitForELBLoadBalancer(networkingClient, lbID, "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	// ? d.SetId(member.ID)

	return resourceBackendRead(d, meta)
}

func resourceBackendRead(d *schema.ResourceData, meta interface{}) error {
	// config := meta.(*Config)

	uri := d.Get("uri").(string)
	job_id := d.Get("job_id").(string)

	log.Printf("[DEBUG] Retrieved uri %s: job_id %#v", uri, job_id)

	//d.Set("region", GetRegion(d, config))

	return nil
}

func resourceBackendDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	log.Printf("[DEBUG] Attempting to delete backend member %s", d.Id())
	// ??

	// Wait for LB to become ACTIVE
	lbID := d.Get("loadbalancer_id").(string)
	// Wait for backend  to become active before continuing
	timeout := d.Timeout(schema.TimeoutDelete)
	err = waitForELBLoadBalancer(networkingClient, lbID, "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	return nil
}
