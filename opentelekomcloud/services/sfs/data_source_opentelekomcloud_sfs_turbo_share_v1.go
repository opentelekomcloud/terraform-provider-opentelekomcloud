package sfs

import (
	"context"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/sfs_turbo/v1/shares"
	"github.com/opentelekomcloud/gophertelekomcloud/pagination"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceSFSTurboShareV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSFSTurboShareV1Read,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: common.ValidateName,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"share_proto": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"crypt_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"export_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_capacity": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSFSTurboShareV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SfsTurboV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud SFSTurboV1 client: %s", err)
	}

	var share *shares.Turbo
	name := d.Get("name")
	err = shares.List(client, shares.ListOpts{}).EachPage(func(p pagination.Page) (bool, error) {
		results, err := shares.ExtractTurbos(p)
		if err != nil {
			return false, err
		}
		for _, turbo := range results {
			if turbo.Name == name {
				share = &turbo
				return false, nil
			}
		}
		return true, nil
	})
	if err != nil {
		return fmterr.Errorf("error listing SFS turbo shares: %w", err)
	}

	if share == nil {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again")
	}

	// n.Size is a string of float64, should convert it to int
	fSize, err := strconv.ParseFloat(share.Size, 64)
	if err != nil {
		return fmterr.Errorf("error parsing SFS Turbo sharing size: %w", err)
	}

	d.SetId(share.ID)
	mErr := multierror.Append(nil,
		d.Set("name", share.Name),
		d.Set("share_proto", share.ShareProto),
		d.Set("share_type", share.ShareType),
		d.Set("vpc_id", share.VpcID),
		d.Set("subnet_id", share.SubnetID),
		d.Set("security_group_id", share.SecurityGroupID),
		d.Set("version", share.Version),
		d.Set("region", config.GetRegion(d)),
		d.Set("availability_zone", share.AvailabilityZone),
		d.Set("available_capacity", share.AvailCapacity),
		d.Set("export_location", share.ExportLocation),
		d.Set("crypt_key_id", share.CryptKeyID),
		d.Set("size", fSize),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}
