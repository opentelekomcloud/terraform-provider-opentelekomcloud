package cci

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cci/v2/namespace"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCCINamespacesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCINamespacesV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
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
						"finalizers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"resource_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCCINamespacesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)
	client, err := common.ClientFromCtx(ctx, cciClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2Client(region)
	})
	if err != nil {
		return fmterr.Errorf(errCreationCciClient, err)
	}

	nsList := make([]namespace.Namespace, 0)
	if name, ok := d.GetOk("name"); ok {
		resp, err := namespace.Get(client, name.(string))
		if err != nil {
			return diag.Errorf("error getting the namespace (%s) from the server: %s", name.(string), err)
		}
		nsList = append(nsList, *resp)
	} else {
		resp, err := namespace.List(client, namespace.ListOpts{})
		if err != nil {
			return diag.Errorf("error finding the namespace list from the server: %s", err)
		}
		nsList = resp
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("namespaces", flattenNamespaces(nsList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenNamespaces(nsList []namespace.Namespace) []map[string]interface{} {
	if len(nsList) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(nsList))
	for _, ns := range nsList {
		result = append(result, map[string]interface{}{
			"name":               ns.Metadata.Name,
			"api_version":        ns.APIVersion,
			"kind":               ns.Kind,
			"annotations":        ns.Metadata.Annotations,
			"labels":             ns.Metadata.Labels,
			"creation_timestamp": ns.Metadata.CreationTimestamp,
			"resource_version":   ns.Metadata.ResourceVersion,
			"uid":                ns.Metadata.UID,
			"finalizers":         ns.Spec.Finalizers,
			"status":             ns.Status.Phase,
		})
	}
	return result
}
