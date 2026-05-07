package waf

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/waf-premium/v1/cloud"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

const postPaidDomainResourceType = "hws.resource.type.waf.payperusedomain"
const chargingModePostPaid = "postPaid"

func ResourceWafCloudInstanceV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudInstanceCreate,
		ReadContext:   resourceCloudInstanceRead,
		DeleteContext: resourceCloudInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  chargingModePostPaid,
				ValidateFunc: validation.StringInSlice([]string{
					chargingModePostPaid,
				}, false),
				Description: "Specifies the charging mode of the cloud WAF. Only `postPaid` is currently supported.",
			},
			"website": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the website to which the account belongs.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the ID of the enterprise project to which the cloud WAF belongs.",
			},
			"resource_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The specification code returned by the cloud WAF API.",
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The current status of the cloud WAF.",
			},
		},
	}
}

func resourceCloudInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := validateCloudInstanceChargingMode(d); err != nil {
		return diag.FromErr(err)
	}

	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	resp, err := cloud.Enable(client, cloud.EnableOpts{
		ConsoleArea:         d.Get("website").(string),
		EnterpriseProjectID: config.GetEnterpriseProjectID(d),
	})
	if err != nil {
		return diag.Errorf("error creating postpaid cloud WAF: %s", err)
	}

	resourceID, err := flattenPostPaidResourceID(resp.Resources)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resourceID)

	return resourceCloudInstanceRead(ctx, d, meta)
}

func resourceCloudInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	subscription, err := cloud.Get(client)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "cloud WAF")
	}

	res, ok := findPostPaidResource(subscription.Resources, d.Id())
	if !ok {
		d.SetId("")
		return nil
	}

	mErr := multierror.Append(nil,
		d.Set("resource_spec_code", res.ResourceSpecCode),
		d.Set("status", res.Status),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving cloud WAF fields: %s", err)
	}

	return nil
}

func resourceCloudInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := validateCloudInstanceChargingMode(d); err != nil {
		return diag.FromErr(err)
	}

	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.WafDedicatedV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1DedicatedClient, err)
	}

	err = cloud.Disable(client, cloud.DeleteOpts{
		EnterpriseProjectID: config.GetEnterpriseProjectID(d),
	})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "cloud WAF")
	}

	return nil
}

func validateCloudInstanceChargingMode(d *schema.ResourceData) error {
	if d.Get("charging_mode").(string) != chargingModePostPaid {
		return fmt.Errorf("unsupported charging_mode %q: only %q is currently supported", d.Get("charging_mode").(string), chargingModePostPaid)
	}

	return nil
}

func flattenPostPaidResourceID(resources []cloud.ResourceResponse) (string, error) {
	for _, resource := range resources {
		if resource.ResourceType == postPaidDomainResourceType {
			return resource.ID, nil
		}
	}

	return "", fmt.Errorf("cannot find target resource type (%s) from response", postPaidDomainResourceType)
}

func findPostPaidResource(resources []cloud.ResourceResponse, id string) (*cloud.ResourceResponse, bool) {
	for i := range resources {
		resource := &resources[i]
		if resource.ResourceType == postPaidDomainResourceType && resource.ID == id {
			return resource, true
		}
	}

	return nil, false
}
