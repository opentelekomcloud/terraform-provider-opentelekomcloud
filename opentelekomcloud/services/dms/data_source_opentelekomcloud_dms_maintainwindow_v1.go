package dms

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v1/maintainwindows"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceDmsMaintainWindowV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsMaintainWindowV1Read,

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

func dataSourceDmsMaintainWindowV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	DmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud dms client: %s", err)
	}

	v, err := maintainwindows.Get(DmsV1Client).Extract()
	if err != nil {
		return diag.FromErr(err)
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
		return fmterr.Errorf("your query returned no results. Please change your filters and try again")
	}
	mw := filteredMVs[0]
	d.SetId(strconv.Itoa(mw.ID))
	mErr := multierror.Append(
		d.Set("begin", mw.Begin),
		d.Set("end", mw.End),
		d.Set("default", mw.Default),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Dms MaintainWindow : %+v", mw)

	return nil
}
