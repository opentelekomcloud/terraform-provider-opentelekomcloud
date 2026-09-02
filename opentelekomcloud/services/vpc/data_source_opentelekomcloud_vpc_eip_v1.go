package vpc

import (
	"context"
	"log"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/publicips"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceVPCEipV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVPCEipV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bandwidth_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bandwidth_share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_border_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allow_share_bandwidth_types": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": common.TagsSchema(),
		},
	}
}

func dataSourceVPCEipV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.VpcV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v1 client: %w", err)
	}

	allEIPs, err := publicips.List(client, publicips.ListOpts{})
	if err != nil {
		return fmterr.Errorf("unable to retrieve EIPs: %w", err)
	}
	refinedEIPs := filterPublicIPs(allEIPs, publicIPFilters{
		ID:               d.Get("id").(string),
		Status:           d.Get("status").(string),
		PrivateAddress:   d.Get("private_ip_address").(string),
		PortID:           d.Get("port_id").(string),
		BandwidthID:      d.Get("bandwidth_id").(string),
		PublicAddress:    d.Get("public_ip_address").(string),
		NameRegexPattern: d.Get("name_regex").(string),
	})

	tagRaw := d.Get("tags").(map[string]interface{})
	var refinedByTags []publicips.PublicIP
	networkingV2Client, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %w", err)
	}
	if len(tagRaw) > 0 {
		tagList := common.ExpandResourceTags(tagRaw)
		for _, eip := range refinedEIPs {
			resourceTagList, err := tags.Get(networkingV2Client, "publicips", eip.ID).Extract()
			if err != nil {
				return fmterr.Errorf("error fetching OpenTelekomCloud VPC EIP tags: %w", err)
			}

			var flag bool
			for _, v := range tagList {
				if common.Contains(resourceTagList, v) {
					flag = true
					continue
				}
				flag = false
				break
			}
			if flag {
				refinedByTags = append(refinedByTags, eip)
			}
		}
	} else {
		refinedByTags = refinedEIPs
	}

	if len(refinedByTags) < 1 {
		return fmterr.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(refinedByTags) > 1 {
		return fmterr.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	elasticIP := refinedByTags[0]

	log.Printf("[INFO] Retrieved ElasticIP using given filter %s: %+v", elasticIP.ID, elasticIP)
	d.SetId(elasticIP.ID)

	mErr := multierror.Append(
		d.Set("status", elasticIP.Status),
		d.Set("id", elasticIP.ID),
		d.Set("type", elasticIP.Type),
		d.Set("bandwidth_id", elasticIP.BandwidthId),
		d.Set("bandwidth_share_type", elasticIP.BandwidthShareType),
		d.Set("bandwidth_size", elasticIP.BandwidthSize),
		d.Set("create_time", elasticIP.CreateTime),
		d.Set("ip_version", elasticIP.IPVersion),
		d.Set("port_id", elasticIP.PortId),
		d.Set("private_ip_address", elasticIP.PrivateIpAddress),
		d.Set("public_ip_address", elasticIP.PublicIpAddress),
		d.Set("tenant_id", elasticIP.TenantId),
		d.Set("region", config.GetRegion(d)),
		d.Set("name", elasticIP.Alias),
		d.Set("enterprise_project_id", elasticIP.EnterpriseProjectId),
		d.Set("public_border_group", elasticIP.PublicBorderGroup),
		d.Set("allow_share_bandwidth_types", elasticIP.AllowShareBandwidthTypes),
	)

	// save tags
	resourceTags, err := tags.Get(networkingV2Client, "publicips", d.Id()).Extract()
	if err != nil {
		return fmterr.Errorf("error fetching OpenTelekomCloud VPC EIP tags: %w", err)
	}
	tagMap := common.TagsToMap(resourceTags)
	mErr = multierror.Append(mErr,
		d.Set("tags", tagMap),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

type publicIPFilters struct {
	ID               string
	Status           string
	PrivateAddress   string
	PortID           string
	BandwidthID      string
	PublicAddress    string
	NameRegexPattern string
}

func filterPublicIPs(allEIPs []publicips.PublicIP, filters publicIPFilters) []publicips.PublicIP {
	var nameRegex *regexp.Regexp
	if filters.NameRegexPattern != "" {
		nameRegex = regexp.MustCompile(filters.NameRegexPattern)
	}

	refined := make([]publicips.PublicIP, 0, len(allEIPs))
	for _, eip := range allEIPs {
		if filters.ID != "" && eip.ID != filters.ID {
			continue
		}
		if filters.Status != "" && eip.Status != filters.Status {
			continue
		}
		if filters.PrivateAddress != "" && eip.PrivateIpAddress != filters.PrivateAddress {
			continue
		}
		if filters.PortID != "" && eip.PortId != filters.PortID {
			continue
		}
		if filters.BandwidthID != "" && eip.BandwidthId != filters.BandwidthID {
			continue
		}
		if filters.PublicAddress != "" && eip.PublicIpAddress != filters.PublicAddress {
			continue
		}
		if nameRegex != nil && !nameRegex.MatchString(eip.Alias) {
			continue
		}
		refined = append(refined, eip)
	}
	return refined
}
