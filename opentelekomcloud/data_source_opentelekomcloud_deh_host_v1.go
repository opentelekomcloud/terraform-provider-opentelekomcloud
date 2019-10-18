package opentelekomcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/deh/v1/hosts"
	"log"
)

func dataSourceDEHHostV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDEHHostV1Read,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_type_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_instance_capacities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flavor": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_placement": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_vcpus": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"available_memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cores": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sockets": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_uuids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vcpus": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceDEHHostV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dehClient, err := config.dehV1Client(GetRegion(d, config))

	listOpts := hosts.ListOpts{
		ID:    d.Id(),
		Name:  d.Get("name").(string),
		State: d.Get("status").(string),
		Az:    d.Get("availability_zone").(string),
	}

	deh, err := hosts.List(dehClient, listOpts).AllPages()
	refinedDeh, err := hosts.ExtractHosts(deh)
	if err != nil {
		return fmt.Errorf("Unable to retrieve dedicated hosts: %s", err)
	}

	if len(refinedDeh) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedDeh) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Deh := refinedDeh[0]

	log.Printf("[INFO] Retrieved Deh using given filter %s: %+v", Deh.ID, Deh)
	d.SetId(Deh.ID)

	d.Set("name", Deh.Name)
	d.Set("id", Deh.ID)
	d.Set("auto_placement", Deh.AutoPlacement)
	d.Set("availability_zone", Deh.Az)
	d.Set("tenant_id", Deh.TenantId)
	d.Set("status", Deh.State)
	d.Set("available_vcpus", Deh.AvailableVcpus)
	d.Set("available_memory", Deh.AvailableMemory)
	d.Set("instance_total", Deh.InstanceTotal)
	d.Set("host_type_name", Deh.HostProperties.HostTypeName)
	d.Set("host_type", Deh.HostProperties.HostType)
	d.Set("cores", Deh.HostProperties.Cores)
	d.Set("sockets", Deh.HostProperties.Sockets)
	d.Set("vcpus", Deh.HostProperties.Vcpus)
	d.Set("memory", Deh.HostProperties.Memory)
	d.Set("available_instance_capacities", getInstanceProperties(&Deh))
	d.Set("instance_uuids", Deh.InstanceUuids)
	d.Set("region", GetRegion(d, config))
	return nil
}
