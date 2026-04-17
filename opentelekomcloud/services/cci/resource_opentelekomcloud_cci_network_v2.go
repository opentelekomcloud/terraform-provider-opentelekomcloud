package cci

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cci/v2/network"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCCINetworkV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCciNetworkV2Create,
		ReadContext:   resourceCciNetworkV2Read,
		UpdateContext: resourceCciNetworkV2Update,
		DeleteContext: resourceCciNetworkV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCciNetworkV2ImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ip_families": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subnets": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"finalizers": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"conditions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"last_transition_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"reason": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"message": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"subnet_attrs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"subnet_v4_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"subnet_v6_id": {
										Type:     schema.TypeString,
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

func resourceCciNetworkV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cciNetworkClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2NetworkClient(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationCciNetworkClient, err)
	}

	createOpts := network.CreateOpts{
		Namespace: d.Get("namespace").(string),
		Metadata: &network.ObjectMeta{
			Name:        d.Get("name").(string),
			Namespace:   d.Get("namespace").(string),
			Annotations: expandStringMap(d.Get("annotations")),
		},
		Spec: &network.NetworkSpec{
			IPFamilies:     expandStringList(d.Get("ip_families").([]interface{})),
			NetworkType:    "underlay_neutron",
			SecurityGroups: expandStringList(d.Get("security_group_ids").([]interface{})),
			Subnets:        expandNetworkSubnets(d.Get("subnets").([]interface{})),
		},
	}

	resp, err := network.Create(client, createOpts)
	if err != nil {
		return diag.Errorf("error creating CCI network: %s", err)
	}

	ns := resp.Metadata.Namespace
	name := resp.Metadata.Name
	if ns == "" || name == "" {
		return diag.Errorf("unable to find CCI network name or namespace from API response")
	}
	d.SetId(ns + "/" + name)

	err = waitForNetworkReady(ctx, client, ns, name, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCciNetworkV2Read(ctx, d, meta)
}

func resourceCciNetworkV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cciNetworkClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2NetworkClient(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationCciNetworkClient, err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)
	resp, err := network.Get(client, ns, name)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying CCI v2 network")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("namespace", resp.Metadata.Namespace),
		d.Set("name", resp.Metadata.Name),
		d.Set("kind", resp.Kind),
		d.Set("api_version", resp.APIVersion),
		d.Set("annotations", resp.Metadata.Annotations),
		d.Set("creation_timestamp", resp.Metadata.CreationTimestamp),
		d.Set("finalizers", resp.Metadata.Finalizers),
		d.Set("resource_version", resp.Metadata.ResourceVersion),
		d.Set("uid", resp.Metadata.UID),
		d.Set("ip_families", resp.Spec.IPFamilies),
		d.Set("security_group_ids", resp.Spec.SecurityGroups),
		d.Set("subnets", flattenNetworkSubnets(resp.Spec.Subnets)),
		d.Set("status", flattenNetworkStatus(resp.Status)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCciNetworkV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cciNetworkClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2NetworkClient(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationCciNetworkClient, err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)

	updateOpts := network.UpdateOpts{
		Namespace:  ns,
		Name:       name,
		APIVersion: d.Get("api_version").(string),
		Kind:       d.Get("kind").(string),
		Metadata: &network.ObjectMeta{
			Name:              name,
			Namespace:         ns,
			UID:               d.Get("uid").(string),
			ResourceVersion:   d.Get("resource_version").(string),
			CreationTimestamp: d.Get("creation_timestamp").(string),
			Annotations:       expandStringMap(d.Get("annotations")),
			Finalizers:        expandStringList(d.Get("finalizers").([]interface{})),
		},
		Spec: &network.NetworkSpec{
			IPFamilies:     expandStringList(d.Get("ip_families").([]interface{})),
			NetworkType:    "underlay_neutron",
			SecurityGroups: expandStringList(d.Get("security_group_ids").([]interface{})),
			Subnets:        expandNetworkSubnets(d.Get("subnets").([]interface{})),
		},
	}

	_, err = network.Update(client, updateOpts)
	if err != nil {
		return diag.Errorf("error updating CCI v2 network: %s", err)
	}

	return resourceCciNetworkV2Read(ctx, d, meta)
}

func resourceCciNetworkV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cciNetworkClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2NetworkClient(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationCciNetworkClient, err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)

	_, err = network.Delete(client, network.DeleteOpts{
		Namespace: ns,
		Name:      name,
		Body:      network.DeleteBody{},
	})
	if err != nil {
		return diag.Errorf("error deleting CCI v2 network: %s", err)
	}

	err = waitForNetworkDeleted(ctx, client, ns, name, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceCciNetworkV2ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<namespace>/<name>', but got '%s'", d.Id())
	}

	mErr := multierror.Append(nil,
		d.Set("namespace", parts[0]),
		d.Set("name", parts[1]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func waitForNetworkReady(ctx context.Context, client *golangsdk.ServiceClient, ns, name string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Ready"},
		Refresh: func() (interface{}, string, error) {
			resp, err := network.Get(client, ns, name)
			if err != nil {
				return nil, "ERROR", err
			}
			if resp.Status.Status != "Ready" {
				return resp, "Pending", nil
			}
			return resp, "Ready", nil
		},
		Timeout:      timeout,
		PollInterval: 10 * time.Second,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CCI network to become ready: %s", err)
	}
	return nil
}

func waitForNetworkDeleted(ctx context.Context, client *golangsdk.ServiceClient, ns, name string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Deleted"},
		Refresh: func() (interface{}, string, error) {
			_, err := network.Get(client, ns, name)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return "", "Deleted", nil
				}
				return nil, "ERROR", err
			}
			return nil, "Pending", nil
		},
		Timeout:      timeout,
		PollInterval: 10 * time.Second,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CCI v2 network to be deleted: %s", err)
	}
	return nil
}

func expandStringList(raw []interface{}) []string {
	if len(raw) == 0 {
		return nil
	}
	result := make([]string, len(raw))
	for i, v := range raw {
		result[i] = v.(string)
	}
	return result
}

func expandStringMap(raw interface{}) map[string]string {
	if raw == nil {
		return nil
	}
	m, ok := raw.(map[string]interface{})
	if !ok || len(m) == 0 {
		return nil
	}
	result := make(map[string]string, len(m))
	for k, v := range m {
		result[k] = v.(string)
	}
	return result
}

func expandNetworkSubnets(raw []interface{}) []network.SubnetConf {
	if len(raw) == 0 {
		return nil
	}
	subnets := make([]network.SubnetConf, len(raw))
	for i, v := range raw {
		m := v.(map[string]interface{})
		subnets[i] = network.SubnetConf{
			SubnetID: m["subnet_id"].(string),
		}
	}
	return subnets
}

func flattenNetworkSubnets(subnets []network.SubnetConf) []map[string]interface{} {
	if len(subnets) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(subnets))
	for i, s := range subnets {
		result[i] = map[string]interface{}{
			"subnet_id": s.SubnetID,
		}
	}
	return result
}

func flattenNetworkStatus(status network.NetworkStatus) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"status":       status.Status,
			"conditions":   flattenNetworkConditions(status.Conditions),
			"subnet_attrs": flattenNetworkSubnetAttrs(status.SubnetAttrs),
		},
	}
}

func flattenNetworkConditions(conditions []network.Condition) []map[string]interface{} {
	if len(conditions) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(conditions))
	for i, c := range conditions {
		result[i] = map[string]interface{}{
			"type":                 c.Type,
			"status":               c.Status,
			"last_transition_time": c.LastTransitionTime,
			"reason":               c.Reason,
			"message":              c.Message,
		}
	}
	return result
}

func flattenNetworkSubnetAttrs(attrs []network.SubnetAttr) []map[string]interface{} {
	if len(attrs) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(attrs))
	for i, a := range attrs {
		result[i] = map[string]interface{}{
			"network_id":   a.NetworkID,
			"subnet_v4_id": a.SubnetV4ID,
			"subnet_v6_id": a.SubnetV6ID,
		}
	}
	return result
}
