package opentelekomcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v1/flowlogs"
)

func resourceVpcFlowLogV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcFlowLogV1Create,
		Read:   resourceVpcFlowLogV1Read,
		Update: resourceVpcFlowLogV1Update,
		Delete: resourceVpcFlowLogV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validateName,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"port", "vpc", "network",
				}, true),
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"traffic_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"all", "accept", "reject",
				}, true),
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_topic_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"admin_state": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceVpcFlowLogV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vpcClient, err := config.networkingV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud vpc client: %s", err)
	}

	createOpts := flowlogs.CreateOpts{
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		ResourceType: d.Get("resource_type").(string),
		ResourceID:   d.Get("resource_id").(string),
		TrafficType:  d.Get("traffic_type").(string),
		LogGroupID:   d.Get("log_group_id").(string),
		LogTopicID:   d.Get("log_topic_id").(string),
	}

	log.Printf("[DEBUG] Create VPC Flow Log Options: %#v", createOpts)
	fl, err := flowlogs.Create(vpcClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud VPC flow log: %s", err)
	}

	d.SetId(fl.ID)
	return resourceVpcFlowLogV1Read(d, config)
}

func resourceVpcFlowLogV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vpcClient, err := config.networkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud vpc client: %s", err)
	}

	fl, err := flowlogs.Get(vpcClient, d.Id()).Extract()
	if err != nil {
		// ignore ErrDefault404
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud flowlog: %s", err)
	}

	d.Set("name", fl.Name)
	d.Set("description", fl.Description)
	d.Set("resource_type", fl.ResourceType)
	d.Set("resource_id", fl.ResourceID)
	d.Set("traffic_type", fl.TrafficType)
	d.Set("log_group_id", fl.LogGroupID)
	d.Set("log_topic_id", fl.LogTopicID)
	d.Set("admin_state", fl.AdminState)
	d.Set("status", fl.Status)

	return nil
}

func resourceVpcFlowLogV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vpcClient, err := config.networkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud vpc client: %s", err)
	}

	var updateOpts flowlogs.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}

	_, err = flowlogs.Update(vpcClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating OpenTelekomCloud VPC flow log: %s", err)
	}

	return resourceVpcFlowLogV1Read(d, meta)
}

func resourceVpcFlowLogV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vpcClient, err := config.networkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud vpc client: %s", err)
	}

	err = flowlogs.Delete(vpcClient, d.Id()).ExtractErr()
	if err != nil {
		// ignore ErrDefault404
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[INFO] Successfully deleted OpenTelekomCloud vpc flow log %s", d.Id())
			return nil
		}
		return err
	}

	d.SetId("")
	return nil
}
