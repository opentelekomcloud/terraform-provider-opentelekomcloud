package v3

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/flavors"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceLBFlavorsV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLBFlavorsV3Read,

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
			"flavors": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceLBFlavorsV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	listOpts := flavors.ListOpts{}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = []string{v.(string)}
	}
	if v, ok := d.GetOk("id"); ok {
		listOpts.ID = []string{v.(string)}
	}

	pages, err := flavors.List(client, listOpts).AllPages()
	if err != nil {
		return fmterr.Errorf("error listing ELBv3 flavors: %w", err)
	}
	flavorList, err := flavors.ExtractFlavors(pages)
	if err != nil {
		return fmterr.Errorf("error extracting ELBv3 flavors: %w", err)
	}

	if len(flavorList) < 1 {
		return fmterr.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	var allFlavors []string
	for _, v := range flavorList {
		allFlavors = append(allFlavors, v.Name)
	}
	d.SetId("flavors")
	mErr := multierror.Append(
		d.Set("flavors", allFlavors),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
