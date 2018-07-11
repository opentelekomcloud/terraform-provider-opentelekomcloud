package opentelekomcloud

import (
	"fmt"
	"log"

	"github.com/huaweicloud/golangsdk/openstack/rts/v1/softwareconfig"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceRtsSoftwareConfigV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRtsSoftwareConfigV1Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"inputs": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"outputs": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_output": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"config": &schema.Schema{
				Type:         schema.TypeString,
				Computed: true,
			},
			"options": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceRtsSoftwareConfigV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orchestrationClient, err := config.orchestrationV1Client(GetRegion(d, config))

	listOpts := softwareconfig.ListOpts{
		Id:   d.Get("id").(string),
		Name: d.Get("name").(string),
	}

	refinedConfigs, err := softwareconfig.List(orchestrationClient, listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve RTS Software Configs: %s", err)
	}

	if len(refinedConfigs) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedConfigs) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Config := refinedConfigs[0]

	log.Printf("[INFO] Retrieved RTS Software Config using given filter %s: %+v", Config.Id, Config)
	d.SetId(Config.Id)

	d.Set("name", Config.Name)
	d.Set("group", Config.Group)
	d.Set("region", GetRegion(d, config))

	n, err := softwareconfig.Get(orchestrationClient,Config.Id).Extract()
	if err != nil {
		return fmt.Errorf("Unable to retrieve RTS Software Config: %s", err)
	}

	var inputvalues []map[string]interface{}
	for _, input := range n.Inputs{
		mapping := map[string]interface{}{
			"description": input.Description,
			"default":     input.Default,
			"type":        input.Type,
			"name":        input.Name,
		}
		inputvalues = append(inputvalues, mapping)
	}

	var outputvalues []map[string]interface{}
	for _, output := range n.Outputs{
		mapping := map[string]interface{}{
			"description":  output.Description,
			"error_output": output.ErrorOutput,
			"type":         output.Type,
			"name":         output.Name,
		}
		outputvalues = append(outputvalues, mapping)
	}

	d.Set("config", n.Config)
	d.Set("options", n.Options)
	if err := d.Set("inputs", inputvalues); err != nil {
		return err
	}
	if err := d.Set("outputs", outputvalues); err != nil {
		return err
	}

	return nil
}