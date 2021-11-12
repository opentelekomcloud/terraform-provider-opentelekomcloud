package v3

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/certificates"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCertificateV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCertificateV3Read,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"server", "client",
				}, false),
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCertificateV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	id := d.Get("id").(string)
	if id != "" {
		cert, err := certificates.Get(client, id).Extract()
		if err != nil {
			return fmterr.Errorf("error getting ELBv3 certificate: %w", err)
		}
		return setCertificateFields(d, cert)
	}

	listOpts := certificates.ListOpts{
		Name:   common.StrSlice(d.Get("name").(string)),
		Type:   common.StrSlice(d.Get("type").(string)),
		Domain: common.StrSlice(d.Get("domain").(string)),
	}

	pages, err := certificates.List(client, listOpts).AllPages()
	if err != nil {
		return fmterr.Errorf("error listing ELBv3 certificates: %w", err)
	}
	certList, err := certificates.ExtractCertificates(pages)
	if err != nil {
		return fmterr.Errorf("error extracting ELBv3 certificates: %w", err)
	}

	if len(certList) > 1 {
		return common.DataSourceTooManyDiag
	}
	if len(certList) < 1 {
		return common.DataSourceTooFewDiag
	}

	cert := &certList[0]
	return setCertificateFields(d, cert)
}

func setCertificateFields(d *schema.ResourceData, cert *certificates.Certificate) diag.Diagnostics {
	d.SetId(cert.ID)
	mErr := multierror.Append(nil,
		d.Set("name", cert.Name),
		d.Set("description", cert.Description),
		d.Set("domain", cert.Domain),
		d.Set("certificate", cert.Certificate),
		d.Set("private_key", cert.PrivateKey),
		d.Set("type", cert.Type),
		d.Set("created_at", cert.CreatedAt),
		d.Set("updated_at", cert.UpdatedAt),
		d.Set("expire_time", cert.ExpireTime),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
