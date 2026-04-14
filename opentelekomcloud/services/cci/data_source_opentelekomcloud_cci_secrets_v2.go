package cci

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cci/v2/secret"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCCISecretsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCISecretsV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secrets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"annotations": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"labels": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"creation_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"string_data": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"data": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"immutable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCCISecretsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)
	client, err := common.ClientFromCtx(ctx, cciClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2Client(region)
	})
	if err != nil {
		return fmterr.Errorf(errCreationCciClient, err)
	}

	ns := d.Get("namespace").(string)

	secretList := make([]secret.Secret, 0)
	if name, ok := d.GetOk("name"); ok {
		resp, err := secret.Get(client, ns, name.(string))
		if err != nil {
			return diag.Errorf("error getting the CCI v2 Secret (%s/%s) from the server: %s", ns, name.(string), err)
		}
		secretList = append(secretList, *resp)
	} else {
		resp, err := secret.List(client, ns, secret.ListOpts{})
		if err != nil {
			return diag.Errorf("error querying CCI v2 Secrets under namespace %s: %s", ns, err)
		}
		secretList = resp
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("secrets", flattenSecrets(secretList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSecrets(secretList []secret.Secret) []map[string]interface{} {
	if len(secretList) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(secretList))
	for _, s := range secretList {
		result = append(result, map[string]interface{}{
			"name":               s.Metadata.Name,
			"namespace":          s.Metadata.Namespace,
			"api_version":        s.APIVersion,
			"kind":               s.Kind,
			"annotations":        s.Metadata.Annotations,
			"labels":             s.Metadata.Labels,
			"creation_timestamp": s.Metadata.CreationTimestamp,
			"resource_version":   s.Metadata.ResourceVersion,
			"uid":                s.Metadata.UID,
			"string_data":        s.StringData,
			"data":               s.Data,
			"type":               s.Type,
			"immutable":          s.Immutable,
		})
	}
	return result
}
