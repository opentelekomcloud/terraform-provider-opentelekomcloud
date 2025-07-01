package ces

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v1/quotas"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCesQuotasV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCesQuotasRead,

		Schema: map[string]*schema.Schema{
			"quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"used": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"unit": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"quota": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceCesQuotasRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CesV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	listQuotas, err := quotas.ShowQuotas(client).Extract()
	if err != nil {
		return fmterr.Errorf("error getting quotas : %w", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Retrieved quotas list %s: %#v", d.Id(), listQuotas)

	quotas := []map[string]interface{}{
		{
			"resources": setResources(listQuotas.Quotas.Resources),
		},
	}
	mErr := multierror.Append(
		d.Set("quotas", quotas),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setResources(resourcesListInResp []quotas.Resource) []map[string]interface{} {
	var resourcesList []map[string]interface{}
	for _, resourcesInResp := range resourcesListInResp {
		resources := map[string]interface{}{
			"type":  resourcesInResp.Type,
			"used":  resourcesInResp.Used,
			"unit":  resourcesInResp.Unit,
			"quota": resourcesInResp.Quota,
		}
		resourcesList = append(resourcesList, resources)
	}
	return resourcesList
}
