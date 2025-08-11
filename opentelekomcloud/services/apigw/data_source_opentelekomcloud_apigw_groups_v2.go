package apigw

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/group"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func DataSourceApigwGroupsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApigwGroupsV2Read,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sl_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"on_sell_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"url_domains": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cname_status": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ssl_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ssl_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"min_ssl_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"verified_client_certificate_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"is_has_trusted_root_ca": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"sl_domains": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"environment": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"variable": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"id": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"environment_id": {
										Type:     schema.TypeString,
										Computed: true,
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

func flattenGroupsUrlDomain(urlDomains []group.UrlDomains) []map[string]interface{} {
	if len(urlDomains) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(urlDomains))
	for i, v := range urlDomains {
		result[i] = map[string]interface{}{
			"id":                                  v.DomainId,
			"name":                                v.DomainName,
			"cname_status":                        v.CnameStatus,
			"ssl_id":                              v.SslID,
			"ssl_name":                            v.SslName,
			"min_ssl_version":                     v.MinSslVersion,
			"verified_client_certificate_enabled": v.VfClientCertEnabled,
			"is_has_trusted_root_ca":              v.HasTrustedCa,
		}
	}

	return result
}

func flattenGroups(group group.GroupResp, variables []map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"id":             group.ID,
		"name":           group.Name,
		"status":         group.Status,
		"sl_domain":      group.SlDomain,
		"created_at":     group.RegisterTime,
		"updated_at":     group.UpdateTime,
		"on_sell_status": group.OnSellStatus,
		"url_domains":    flattenGroupsUrlDomain(group.UrlDomains),
		"sl_domains":     group.SlDomains,
		"description":    group.Description,
		"is_default":     group.IsDefault,
		"environment":    variables,
	}

	return result
}

func dataSourceApigwGroupsV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.APIGWV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	opts := group.ListOpts{
		GatewayID: d.Get("instance_id").(string),
		ID:        d.Get("group_id").(string),
		Name:      d.Get("name").(string),
	}
	groups, err := group.List(client, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	result := make([]map[string]interface{}, 0)
	for _, gr := range groups {
		variables, err := queryEnvironmentVariables(client, d.Get("instance_id").(string), gr.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		result = append(result, flattenGroups(gr, flattenEnvironmentVariables(variables)))
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("groups", result),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving dedicated groups data source fields: %s", mErr)
	}
	return nil
}
