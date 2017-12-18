package opentelekomcloud
import (
"fmt"
"log"

"github.com/gophercloud/gophercloud/openstack/networking/v1/extensions/vpcs"

"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceVirtualPrivateCloudVpcV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVirtualPrivateCloudV1Read,

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
			"cidr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_shared_snat": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}


func dataSourceVirtualPrivateCloudV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vpcClient, err := config.vpcV1Client(GetRegion(d, config))

	listOpts := vpcs.ListOpts{
		ID:       d.Get("id").(string),
		Name:     d.Get("name").(string),
		Status:   d.Get("status").(string),
		CIDR:     d.Get("cidr").(string),
	}

	vpcdict, err := vpcs.List(vpcClient, listOpts).AllPages()
	log.Printf("[DEBUG] Value of vpc response : %#v", vpcdict)
	allVpcs, err := vpcs.ExtractVpcs(vpcdict)
	log.Printf("[DEBUG] Value of allVpcs: %#v", allVpcs)
	if err != nil {
		return fmt.Errorf("Unable to retrieve vpcs: %s", err)
	}

	if len(allVpcs) < 1 {
		return fmt.Errorf("No Vpc found with name: %s", d.Get("name"))
	}

	Vpc := allVpcs[0]

	log.Printf("[DEBUG] Retrieved Vpcs using given filter %s: %+v", Vpc.ID, Vpc)
	d.SetId(Vpc.ID)

	d.Set("name", Vpc.Name)
	d.Set("cidr", Vpc.CIDR)
	d.Set("status", Vpc.Status)
	d.Set("id", Vpc.ID)
	d.Set("enable_shared_snat", Vpc.EnableSharedSnat)
	d.Set("region", GetRegion(d, config))

	return nil
}