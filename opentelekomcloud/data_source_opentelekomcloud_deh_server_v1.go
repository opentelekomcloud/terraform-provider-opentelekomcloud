package opentelekomcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/deh/v1/hosts"
	"log"
)

func dataSourceDEHServersV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDEHServersV1Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"dedicated_host_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"addresses": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"fixed_ip_v4": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDEHServersV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dehClient, err := config.dehV1Client(GetRegion(d, config))

	listServerOpts := hosts.ListServerOpts{
		ID:     d.Get("server_id").(string),
		Name:   d.Get("name").(string),
		Status: d.Get("status").(string),
		UserID: d.Get("user_id").(string),
	}
	pages, err := hosts.ListServer(dehClient, d.Get("dedicated_host_id").(string), listServerOpts)

	if err != nil {
		return fmt.Errorf("Unable to retrieve deh server: %s", err)
	}

	if len(pages) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}
	if len(pages) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	DehServer := pages[0]

	log.Printf("[INFO] Retrieved Deh Server using given filter %s: %+v", DehServer.ID, DehServer)
	d.SetId(DehServer.ID)

	d.Set("server_id", DehServer.ID)
	d.Set("user_id", DehServer.UserID)
	d.Set("name", DehServer.Name)
	d.Set("status", DehServer.Status)
	d.Set("flavor", DehServer.Flavor)
	d.Set("addresses", DehServer.Addresses)
	d.Set("metadata", DehServer.Metadata)
	d.Set("tenant_id", DehServer.TenantID)
	d.Set("region", GetRegion(d, config))
	networks, err := flattenInstanceNetwork(d, meta)
	if err != nil {
		return err
	}
	if err := d.Set("addresses", networks); err != nil {
		return fmt.Errorf("[DEBUG] Error saving network to state for OpenTelekomCloud server (%s): %s", d.Id(), err)
	}

	return nil

}
