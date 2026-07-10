package vpc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v3/addressgroup"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceVpcIPAddressGroupV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcIPAddressGroupV3Read,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_capacity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ip_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ip_extra_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remarks": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVpcIPAddressGroupV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	d.SetId(d.Get("id").(string))
	addressGroup, err := addressgroup.Get(client, d.Id())
	if err != nil {
		return fmterr.Errorf("unable to retrieve ip address group: %s", err)
	}

	ipExtraSet := make([]interface{}, 0, len(addressGroup.IPExtraSet))
	for _, v := range addressGroup.IPExtraSet {
		ipExtraSet = append(ipExtraSet, map[string]interface{}{
			"ip":      v.IP,
			"remarks": v.Remark,
		})
	}

	mErr := multierror.Append(nil,
		d.Set("name", addressGroup.Name),
		d.Set("description", addressGroup.Description),
		d.Set("ip_version", addressGroup.IPVersion),
		d.Set("enterprise_project_id", addressGroup.EnterpriseProjectID),
		d.Set("max_capacity", addressGroup.MaxCapacity),
		d.Set("ip_set", addressGroup.IPSet),
		d.Set("ip_extra_set", ipExtraSet),
		d.Set("project_id", addressGroup.TenantID),
		d.Set("created_at", addressGroup.CreatedAt),
		d.Set("updated_at", addressGroup.UpdatedAt),
		d.Set("status", addressGroup.Status),
		d.Set("status_message", addressGroup.StatusMessage),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
