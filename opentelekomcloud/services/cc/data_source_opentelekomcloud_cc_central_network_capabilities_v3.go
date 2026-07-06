package cc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cc/v3/capability"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCcCentralNetworkCapabilitiesV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCcCentralNetworkCapabilitiesV3Read,

		Schema: map[string]*schema.Schema{
			"capability": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"capabilities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"capability": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"specifications": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_support": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"charge_mode": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"support_regions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"support_ipv6_regions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"support_dscp_regions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"support_sts5_regions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"support_freeze_regions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"size_range": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"min": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"max": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"free_lines": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"local_site_code": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"remote_site_code": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"support_sites": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"region_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"site_code": {
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
				},
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceCcCentralNetworkCapabilitiesV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CcV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	var capabilities []capability.CapabilityEntry
	marker := ""
	for {
		resp, err := capability.List(client, capability.ListOpts{
			Marker:     marker,
			Capability: stringFilter(d.Get("capability").(string)),
		})
		if err != nil {
			return diag.Errorf("error retrieving OpenTelekomCloud CC central network capabilities: %s", err)
		}
		capabilities = append(capabilities, resp.Capabilities...)
		if resp.PageInfo.NextMarker == "" || len(resp.Capabilities) == 0 {
			break
		}
		marker = resp.PageInfo.NextMarker
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("capabilities", flattenCentralNetworkCapabilities(capabilities)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCentralNetworkCapabilities(capabilities []capability.CapabilityEntry) []map[string]interface{} {
	if len(capabilities) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(capabilities))
	for i, c := range capabilities {
		result[i] = map[string]interface{}{
			"id":             c.ID,
			"domain_id":      c.DomainId,
			"capability":     c.Capability,
			"specifications": flattenCapabilitySpecifications(c.Specifications),
		}
	}
	return result
}

func flattenCapabilitySpecifications(spec capability.CapabilitySpecifications) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"is_support":             spec.IsSupport,
			"charge_mode":            spec.ChargeMode,
			"support_regions":        spec.SupportRegions,
			"support_ipv6_regions":   spec.SupportIpv6Regions,
			"support_dscp_regions":   spec.SupportDscpRegions,
			"support_sts5_regions":   spec.SupportSts5Regions,
			"support_freeze_regions": spec.SupportFreezeRegions,
			"size_range": []map[string]interface{}{
				{
					"min": spec.SizeRange.Min,
					"max": spec.SizeRange.Max,
				},
			},
			"free_lines":    flattenCapabilityFreeLines(spec.FreeLines),
			"support_sites": flattenCapabilitySupportSites(spec.SupportSites),
		},
	}
}

func flattenCapabilityFreeLines(lines []capability.ConnectionBandwidthFreeLine) []map[string]interface{} {
	if len(lines) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(lines))
	for i, line := range lines {
		result[i] = map[string]interface{}{
			"local_site_code":  line.LocalSiteCode,
			"remote_site_code": line.RemoteSiteCode,
		}
	}
	return result
}

func flattenCapabilitySupportSites(sites []capability.SiteSpecifications) []map[string]interface{} {
	if len(sites) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(sites))
	for i, site := range sites {
		result[i] = map[string]interface{}{
			"region_id": site.RegionId,
			"site_code": site.SiteCode,
		}
	}
	return result
}
