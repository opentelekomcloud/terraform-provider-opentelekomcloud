package dcs

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dcs/v2/ssl"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceDcsCertificateV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsCertificateV2Read,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"file_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsCertificateV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DcsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating dcs key client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)

	v, err := ssl.DownloadCert(client, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instanceID)

	cert, err := downloadCert(v.Link)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("file_name", v.FileName),
		d.Set("link", v.Link),
		d.Set("bucket_name", v.BucketName),
		d.Set("certificate", cert),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func downloadCert(link string) (string, error) {
	resp, err := http.Get(link)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return "", err
	}

	f, err := zipReader.File[0].Open()
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	fileContent, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(fileContent), nil
}
