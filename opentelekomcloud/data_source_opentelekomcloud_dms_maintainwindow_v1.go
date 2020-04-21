package opentelekomcloud

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/dms/v1/maintainwindows"
)

func dataSourceDmsMaintainWindowV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDmsMaintainWindowV1Read,

		Schema: map[string]*schema.Schema{
			"seq": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"begin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"end": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceDmsMaintainWindowV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dmsV1Client, err := config.dmsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud dms client: %s", err)
	}

	v, err := maintainwindows.Get(dmsV1Client).Extract()
	if err != nil {
		return err
	}

	seq := d.Get("seq").(int)
	begin := d.Get("begin").(string)
	end := d.Get("end").(string)
	df := d.Get("default").(bool)

	maintainWindows := v.MaintainWindows
	var filteredMVs []maintainwindows.MaintainWindow
	for _, mv := range maintainWindows {
		if seq != 0 && mv.ID != seq {
			continue
		}

		if begin != "" && mv.Begin != begin {
			continue
		}
		if end != "" && mv.End != end {
			continue
		}

		if mv.Default != df {
			continue
		}
		filteredMVs = append(filteredMVs, mv)
	}
	if len(filteredMVs) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your filters and try again.")
	}
	mw := filteredMVs[0]
	d.SetId(strconv.Itoa(mw.ID))
	d.Set("begin", mw.Begin)
	d.Set("end", mw.End)
	d.Set("default", mw.Default)
	log.Printf("[DEBUG] Dms MaintainWindow : %+v", mw)

	return nil
}
