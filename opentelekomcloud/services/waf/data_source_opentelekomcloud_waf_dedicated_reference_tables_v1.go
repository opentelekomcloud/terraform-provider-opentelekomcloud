package waf

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf-premium/v1/rules"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func DataSourceWafDedicatedRefTablesV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWafDedicatedRefTablesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"conditions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceWafDedicatedRefTablesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.WafDedicatedV1Client(config.GetRegion(d))
	if err != nil {
		return diag.FromErr(err)
	}

	refTables, err := rules.ListReferenceTable(client, rules.ListReferenceTableOpts{})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error obtain WAF Dedicated reference table information")
	}

	if len(refTables) == 0 {
		return nil
	}
	// filter data by name
	filterData, err := common.FilterSliceWithField(refTables, map[string]interface{}{
		"Name": d.Get("name").(string),
	})
	if err != nil {
		return diag.Errorf("error filtering OpenTelekomCloud WAF Dedicated reference tables: %s", err)
	}
	tables := make([]map[string]interface{}, 0, len(filterData))
	ids := make([]string, 0, len(refTables))
	for _, t := range filterData {
		v := t.(rules.ReferenceTable)
		tab := map[string]interface{}{
			"id":          v.ID,
			"name":        v.Name,
			"type":        v.Type,
			"conditions":  v.Values,
			"description": v.Description,
			"created_at":  time.Unix(v.CreatedAt/1000, 0).Format("2006-01-02 15:04:05"),
		}
		tables = append(tables, tab)
		ids = append(ids, v.ID)
	}

	d.SetId(hashcode.Strings(ids))
	mErr := multierror.Append(nil, d.Set("tables", tables))

	return diag.FromErr(mErr.ErrorOrNil())
}
