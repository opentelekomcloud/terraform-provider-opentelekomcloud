package opentelekomcloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/nodes"
)

func dataSourceCceNodeIdsV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCceNodeIdsV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func dataSourceCceNodeIdsV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Unable to create opentelekomcloud CCE client : %s", err)
	}

	var listOpts nodes.ListOpts
	refinedNodes, err := nodes.List(cceClient, d.Get("cluster_id").(string), listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve Nodes: %s", err)
	}

	if len(refinedNodes) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	Nodes := make([]string, 0)
	for _, node := range refinedNodes {
		Nodes = append(Nodes, node.Metadata.Id)
	}

	d.SetId(d.Get("cluster_id").(string))
	d.Set("ids", Nodes)
	d.Set("region", GetRegion(d, config))

	return nil
}
