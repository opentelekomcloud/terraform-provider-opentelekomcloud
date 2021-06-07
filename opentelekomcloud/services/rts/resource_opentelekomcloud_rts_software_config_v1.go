package rts

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rts/v1/softwareconfig"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceSoftwareConfigV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSoftwareConfigV1Create,
		ReadContext:   resourceSoftwareConfigV1Read,
		DeleteContext: resourceSoftwareConfigV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ // request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"config": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"group": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"options": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"input_values": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
			"output_values": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
		},
	}
}

func resourceOptionsV1(d *schema.ResourceData) map[string]interface{} {
	m := make(map[string]interface{})
	for key, val := range d.Get("options").(map[string]interface{}) {
		m[key] = val.(string)
	}

	return m
}

func resourceSoftwareConfigV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	orchastrationClient, err := config.OrchestrationV1Client(config.GetRegion(d))

	if err != nil {
		return fmterr.Errorf("Error creating OpenTelekomCloud RTS client: %s", err)
	}
	input := d.Get("input_values").([]interface{})

	inputs := make([]map[string]interface{}, len(input))
	for i, v := range input {
		inputs[i] = v.(map[string]interface{})
	}

	output := d.Get("output_values").([]interface{})

	outputs := make([]map[string]interface{}, len(output))
	for i, v := range output {
		outputs[i] = v.(map[string]interface{})
	}
	createOpts := softwareconfig.CreateOpts{
		Name:    d.Get("name").(string),
		Config:  d.Get("config").(string),
		Group:   d.Get("group").(string),
		Inputs:  inputs,
		Outputs: outputs,
		Options: resourceOptionsV1(d),
	}

	n, err := softwareconfig.Create(orchastrationClient, createOpts).Extract()

	if err != nil {
		return fmterr.Errorf("Error creating OpenTelekomCloud RTS Software Config: %s", err)
	}
	d.SetId(n.Id)

	return resourceSoftwareConfigV1Read(ctx, d, meta)
}

func resourceSoftwareConfigV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	orchastrationClient, err := config.OrchestrationV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("Error creating OpenTelekomCloud RTS client: %s", err)
	}

	n, err := softwareconfig.Get(orchastrationClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("Error retrieving OpenTelekomCloud Vpc: %s", err)
	}

	d.Set("id", n.Id)
	d.Set("name", n.Name)
	d.Set("config", n.Config)
	d.Set("group", n.Group)
	d.Set("options", n.Options)
	d.Set("region", config.GetRegion(d))
	if err := d.Set("input_values", n.Inputs); err != nil {
		return fmterr.Errorf("[DEBUG] Error saving inputs to state for OpenTelekomCloud RTS Software Config (%s): %s", d.Id(), err)
	}
	if err := d.Set("output_values", n.Outputs); err != nil {
		return fmterr.Errorf("[DEBUG] Error saving outputs to state for OpenTelekomCloud RTS Software Config (%s): %s", d.Id(), err)
	}
	return nil
}

func resourceSoftwareConfigV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	orchastrationClient, err := config.OrchestrationV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("Error creating OpenTelekomCloud vpc: %s", err)
	}
	err = softwareconfig.Delete(orchastrationClient, d.Id()).ExtractErr()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[INFO] Successfully deleted OpenTelekomCloud RTS Software Config %s", d.Id())

		}
		if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
			if errCode.Actual == 409 {
				log.Printf("[INFO] Error deleting OpenTelekomCloud RTS Software Config %s", d.Id())
			}
		}
		log.Printf("[INFO] Successfully deleted OpenTelekomCloud RTS Software Config %s", d.Id())
	}

	d.SetId("")
	return nil
}
