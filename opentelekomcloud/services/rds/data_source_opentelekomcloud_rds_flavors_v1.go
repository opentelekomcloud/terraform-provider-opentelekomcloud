package rds

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v1/datastores"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v1/flavors"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func DataSourceRdsFlavorV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcedataSourceRdsFlavorV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"datastore_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"datastore_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"speccode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourcedataSourceRdsFlavorV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)

	rdsClient, err := config.RdsV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating OpenTelekomCloud rds client: %s", err)
	}

	datastoresList, err := datastores.List(rdsClient, d.Get("datastore_name").(string)).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve datastores: %s ", err)
	}

	if len(datastoresList) < 1 {
		return diag.Errorf("Returned no datastore result. ")
	}
	var datastoreId string
	for _, datastore := range datastoresList {
		if datastore.Name == d.Get("datastore_version").(string) {
			datastoreId = datastore.ID
			break
		}
	}
	if datastoreId == "" {
		return diag.Errorf("Returned no datastore ID. ")
	}
	log.Printf("[DEBUG] Received datastore Id: %s", datastoreId)

	flavorsList, err := flavors.List(rdsClient, datastoreId, d.Get("region").(string)).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve flavors: %s", err)
	}
	if len(flavorsList) < 1 {
		return diag.Errorf("Returned no flavor result. ")
	}

	var rdsFlavor flavors.Flavor
	if d.Get("speccode").(string) == "" {
		rdsFlavor = flavorsList[0]
	} else {
		for _, flavor := range flavorsList {
			if flavor.SpecCode == d.Get("speccode").(string) {
				rdsFlavor = flavor
				break
			}
		}
	}
	log.Printf("[DEBUG] Retrieved flavorId %s: %+v ", rdsFlavor.ID, rdsFlavor)
	if rdsFlavor.ID == "" {
		return diag.Errorf("Returned no flavor Id. ")
	}

	d.SetId(rdsFlavor.ID)

	d.Set("name", rdsFlavor.Name)
	d.Set("ram", rdsFlavor.Ram)
	d.Set("speccode", rdsFlavor.SpecCode)
	d.Set("region", config.GetRegion(d))

	return nil
}
