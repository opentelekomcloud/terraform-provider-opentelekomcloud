package opentelekomcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/rts/v1/deployment"
	"log"
	"time"

	"github.com/huaweicloud/golangsdk"
)

func resourceRTSSoftwareDeploymentV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceRTSSoftwareDeploymentV1Create,
		Read:   resourceRTSSoftwareDeploymentV1Read,
		Update: resourceRTSSoftwareDeploymentV1Update,
		Delete: resourceRTSSoftwareDeploymentV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ //request and response parameters
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"config_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
			},
			"server_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
			},
			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"action": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"input_values": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"status_reason": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"output_values": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceInputValuesV1(d *schema.ResourceData) map[string]interface{} {
	m := make(map[string]interface{})
	for key, val := range d.Get("input_values").(map[string]interface{}) {
		m[key] = val.(string)
	}

	return m
}

func resourceOutputValuesV1(d *schema.ResourceData) map[string]interface{} {
	m := make(map[string]interface{})
	for key, val := range d.Get("output_values").(map[string]interface{}) {
		m[key] = val.(string)
	}

	return m
}

func resourceRTSSoftwareDeploymentV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orchestrationClient, err := config.orchestrationV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Orchestration client: %s", err)
	}

	createOpts := deployment.CreateOpts{
		Action: d.Get("action").(string),
		ConfigId: d.Get("config_id").(string),
		ServerId: d.Get("server_id").(string),
		StatusReason: d.Get("status_reason").(string),
		Status: d.Get("status").(string),
		TenantId: d.Get("tenant_id").(string),
		InputValues:resourceInputValuesV1(d),
	}

	n, err := deployment.Create(orchestrationClient, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Software Deployment: %s", err)
	}

	d.SetId(n.Id)

	log.Printf("[INFO] Software Deployment ID: %s", n.Id)

	return resourceRTSSoftwareDeploymentV1Read(d, meta)

}

func resourceRTSSoftwareDeploymentV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orchestrationClient, err := config.orchestrationV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Orchestration client: %s", err)
	}

	n, err := deployment.Get(orchestrationClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud Software Deployment: %s", err)
	}

	d.Set("id", n.Id)
	d.Set("config_id", n.ConfigId)
	d.Set("status", n.Status)
	d.Set("status_reason", n.StatusReason)
	d.Set("server_id", n.ServerId)
	d.Set("output_values", n.OutputValues)
	d.Set("input_values", n.InputValues)
	d.Set("action", n.Action)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceRTSSoftwareDeploymentV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orchestrationClient, err := config.orchestrationV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud orchestration client: %s", err)
	}

	var updateOpts deployment.UpdateOpts

	updateOpts.ConfigId = d.Get("config_id").(string)
	updateOpts.OutputValues = resourceOutputValuesV1(d)

	if d.HasChange("status") {
		updateOpts.Status = d.Get("status").(string)
	}
	if d.HasChange("action") {
		updateOpts.Action = d.Get("action").(string)
	}
	if d.HasChange("status_reason") {
		updateOpts.StatusReason = d.Get("status_reason").(string)
	}
	if d.HasChange("input_values") {
		updateOpts.InputValues = resourceInputValuesV1(d)
	}

	_, err = deployment.Update(orchestrationClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating OpenTelekomCloud Software Deployment: %s", err)
	}

	return resourceRTSSoftwareDeploymentV1Read(d, meta)
}

func resourceRTSSoftwareDeploymentV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orchestrationClient, err := config.orchestrationV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Orchestration: %s", err)
	}

	err = deployment.Delete(orchestrationClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud Software Deployment: %s", err)
	}

	d.SetId("")
	return nil
}
