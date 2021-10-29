package v1

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/elbaas/loadbalancer_elbs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceELoadBalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceELoadBalancerCreate,
		ReadContext:   resourceELoadBalancerRead,
		UpdateContext: resourceELoadBalancerUpdate,
		DeleteContext: resourceELoadBalancerDelete,

		DeprecationMessage: classicLBDeprecated,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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

			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"bandwidth": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 1000),
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Internal", "External",
				}, false),
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},

			"vip_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"az": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"vip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceELoadBalancerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ElbV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	createOpts := loadbalancer_elbs.CreateOpts{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		VpcID:           d.Get("vpc_id").(string),
		Bandwidth:       d.Get("bandwidth").(int),
		Type:            d.Get("type").(string),
		AdminStateUp:    &adminStateUp,
		VipSubnetID:     d.Get("vip_subnet_id").(string),
		AZ:              d.Get("az").(string),
		SecurityGroupID: d.Get("security_group_id").(string),
		VipAddress:      d.Get("vip_address").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	job, err := loadbalancer_elbs.Create(client, createOpts).ExtractJobResponse()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := waitForLBV1JobSuccess(client, job.URI, int(d.Timeout(schema.TimeoutCreate)/time.Second)); err != nil {
		return diag.FromErr(err)
	}

	entity, err := getLBV1JobEntity(client, job.URI, "elb")
	if err != nil {
		return fmterr.Errorf("error getting job entity: %w", err)
	}

	if mlb, ok := entity.(map[string]interface{}); ok {
		if vid, ok := mlb["id"]; ok {
			if id, ok := vid.(string); ok {
				// If all has been successful, set the ID on the resource, return Read of it
				d.SetId(id)
				return resourceELoadBalancerRead(ctx, d, meta)
			}
		}
	}

	return fmterr.Errorf("unexpected conversion error in resourceELoadBalancerCreate.")
}

func resourceELoadBalancerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.ElbV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	lb, err := loadbalancer_elbs.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "loadbalancer"))
	}

	log.Printf("[DEBUG] Retrieved loadbalancer %s: %#v", d.Id(), lb)

	basu := false
	// Can be 0 (not up) or 2 (frozen)
	if lb.AdminStateUp == 1 {
		basu = true
	}

	mErr := multierror.Append(
		d.Set("name", lb.Name),
		d.Set("description", lb.Description),
		d.Set("vpc_id", lb.VpcID),
		d.Set("bandwidth", lb.Bandwidth),
		d.Set("type", lb.Type),
		d.Set("admin_state_up", basu),
		d.Set("vip_subnet_id", lb.VipSubnetID),
		d.Set("az", lb.AZ),
		d.Set("vip_address", lb.VipAddress),
		d.Set("security_group_id", lb.SecurityGroupID),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceELoadBalancerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ElbV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
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

	log.Printf("[DEBUG] Updating loadbalancer %s with options: %#v", d.Id(), updateOpts)
	job, err := loadbalancer_elbs.Update(client, d.Id(), updateOpts).ExtractJobResponse()
	if err != nil {
		return fmterr.Errorf("error getting job entity: %w", err)
	}
	if err := golangsdk.WaitForJobSuccess(client, job.URI, int(d.Timeout(schema.TimeoutUpdate)/time.Second)); err != nil {
		return diag.FromErr(err)
	}

	return resourceELoadBalancerRead(ctx, d, meta)
}

func resourceELoadBalancerDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ElbV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	id := d.Id()
	log.Printf("[DEBUG] Deleting loadbalancer %s", d.Id())
	job, err := loadbalancer_elbs.Delete(client, id, false).ExtractJobResponse()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("Waiting for loadbalancer %s to delete", id)

	if err := golangsdk.WaitForJobSuccess(client, job.URI, int(d.Timeout(schema.TimeoutDelete)/time.Second)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("Successfully deleted loadbalancer %s", id)
	return nil
}

func waitForLBV1JobSuccess(client *golangsdk.ServiceClient, uri string, secs int) error {
	return golangsdk.WaitFor(secs, func() (bool, error) {
		job := new(golangsdk.JobStatus)
		_, err := client.Get(golangsdk.GetJobEndpoint(client.Endpoint)+uri, job, nil)
		if err != nil {
			return false, err
		}

		if job.Status == "SUCCESS" {
			return true, nil
		}
		if job.Status == "FAIL" {
			err = fmt.Errorf("Job failed with code %s: %s.\n", job.ErrorCode, job.FailReason)
			return false, err
		}

		return false, nil
	})
}

func getLBV1JobEntity(client *golangsdk.ServiceClient, uri string, label string) (interface{}, error) {
	job := new(golangsdk.JobStatus)
	_, err := client.Get(golangsdk.GetJobEndpoint(client.Endpoint)+uri, &job, nil)
	if err != nil {
		return nil, err
	}

	if job.Status == "SUCCESS" {
		if e := job.Entities[label]; e != nil {
			return e, nil
		}
	}

	return nil, fmt.Errorf("unexpected conversion error in GetJobEntity")
}
