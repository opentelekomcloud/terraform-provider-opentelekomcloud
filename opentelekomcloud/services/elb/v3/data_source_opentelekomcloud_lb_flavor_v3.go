package v3

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/flavors"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceLBFlavorV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLBFlavorV3Read,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shared": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"max_connections": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cps": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"qps": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceLBFlavorV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	id := d.Get("id").(string)
	if id != "" {
		flavor, err := flavors.Get(client, id).Extract()
		if err != nil {
			return fmterr.Errorf("error getting ELBv3 flavor: %w", err)
		}
		return setFlavorFields(d, flavor)
	}

	listOpts := flavors.ListOpts{}
	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = []string{v.(string)}
	}

	pages, err := flavors.List(client, listOpts).AllPages()
	if err != nil {
		return fmterr.Errorf("error listing ELBv3 flavors: %w", err)
	}
	flavorList, err := flavors.ExtractFlavors(pages)
	if err != nil {
		return fmterr.Errorf("error extracting ELBv3 flavors: %w", err)
	}

	if len(flavorList) > 1 {
		return common.DataSourceTooManyDiag
	}
	if len(flavorList) < 1 {
		return common.DataSourceTooFewDiag
	}

	flavor := &flavorList[0]
	return setFlavorFields(d, flavor)
}

func setFlavorFields(d *schema.ResourceData, flavor *flavors.Flavor) diag.Diagnostics {
	d.SetId(flavor.ID)
	mErr := multierror.Append(
		d.Set("name", flavor.Name),
		d.Set("shared", flavor.Shared),
		d.Set("type", flavor.Type),
		d.Set("max_connections", flavor.Info.Connection),
		d.Set("cps", flavor.Info.Cps),
		d.Set("qps", flavor.Info.Qps),
		d.Set("bandwidth", flavor.Info.Bandwidth),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting ELBv3 flavor fields: %w", err)
	}

	return nil
}
