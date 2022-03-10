package cce

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/clusters"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCCEClusterKubeConfigV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCEClusterKubeConfigV3Read,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"expiration": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},
			"kubeconfig": {
				Type:     schema.TypeString,
				Computed: true,
				StateFunc: func(v interface{}) string {
					jsonString, _ := common.NormalizeJsonString(v)
					return jsonString
				},
			},
		},
	}
}

func dataSourceCCEClusterKubeConfigV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CceV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(cceClientError, err)
	}

	clusterID := d.Get("cluster_id").(string)
	expiration := d.Get("expiration").(int)
	expiryOpts := clusters.ExpirationOpts{
		Duration: expiration,
	}

	kubeconfig, err := clusters.GetCertWithExpiration(client, clusterID, expiryOpts).ExtractMap()
	if err != nil {
		return fmterr.Errorf("unable to retrieve cluster kubeconfig: %w", err)
	}

	d.SetId(clusterID)

	mErr := multierror.Append(nil,
		d.Set("cluster_id", clusterID),
		d.Set("expiration", expiration),
		d.Set("kubeconfig", kubeconfig),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
