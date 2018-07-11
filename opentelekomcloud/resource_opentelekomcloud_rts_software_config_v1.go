package opentelekomcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"time"
	"fmt"
	"log"
	"github.com/huaweicloud/golangsdk/openstack/rts/v1/softwareconfig"
	"github.com/huaweicloud/golangsdk"
)

func resourceSoftwareConfigV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceSoftwareConfigV1Create,
		Read:   resourceSoftwareConfigV1Read,
		Delete: resourceSoftwareConfigV1Delete,
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
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
			},
			"config": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
			},
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew:true,
			},
			"options": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew:true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"inputs": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
							Computed: true,
						},
					},
				},
			},
			"outputs": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
							Computed: true,
						},
						"error_output": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
							Computed: true,
						},
					},
				},
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
func resourceInputsV1(d *schema.ResourceData) []softwareconfig.Inputs {
	rawInputs := d.Get("inputs").([]interface{})
	inputs := make([]softwareconfig.Inputs, len(rawInputs))
	for i, raw := range rawInputs {
		rawMap := raw.(map[string]interface{})
		inputs[i] = softwareconfig.Inputs{
			Default:rawMap["default"].(string),
			Type:   rawMap["type"].(string),
			Name:   rawMap["name"].(string),
			Description:rawMap["description"].(string),
		}
	}
	log.Printf("[DEBUG] input %s", inputs)
	return inputs
}

func resourceOutputsV1(d *schema.ResourceData)[]softwareconfig.Outputs {
	rawOutputs := d.Get("outputs").([]interface{})
	outputs := make([]softwareconfig.Outputs, len(rawOutputs))
	for i, raw := range rawOutputs {
		rawMap := raw.(map[string]interface{})
		outputs[i] = softwareconfig.Outputs{
			Type:        rawMap["type"].(string),
			Name:        rawMap["name"].(string),
			ErrorOutput: rawMap["error_output"].(bool),
			Description: rawMap["description"].(string),
		}
	}
	log.Printf("[DEBUG] output %s", outputs)
	return outputs
}
func resourceSoftwareConfigV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orchastrationClient, err := config.orchestrationV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud RTS client: %s", err)
	}

	createOpts := softwareconfig.CreateOpts{
		Name:		d.Get("name").(string),
		Config:		d.Get("config").(string),
		Group:		d.Get("group").(string),
		Inputs:		resourceInputsV1(d),
		Outputs:	resourceOutputsV1(d),
		Options:	resourceOptionsV1(d),
	}

	n, err := softwareconfig.Create(orchastrationClient, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud RTS Software Config: %s", err)
	}
	d.SetId(n.Id)


	return resourceSoftwareConfigV1Read(d, meta)

}

func resourceSoftwareConfigV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orchastrationClient, err := config.orchestrationV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud RTS client: %s", err)
	}

	n, err := softwareconfig.Get(orchastrationClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud Vpc: %s", err)
	}

	d.Set("id", n.Id)
	d.Set("name", n.Name)
	d.Set("config", n.Config)
	d.Set("group", n.Group)
	d.Set("inputs", n.Inputs)
	d.Set("outputs", n.Outputs)
	d.Set("options", n.Options)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceSoftwareConfigV1Delete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	orchastrationClient, err := config.orchestrationV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud vpc: %s", err)
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



