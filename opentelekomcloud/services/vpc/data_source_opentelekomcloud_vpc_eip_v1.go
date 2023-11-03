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
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/eips"
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
			"tags": common.TagsSchema(),
		},
	}
}

func dataSourceVPCEipV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud NetworkingV1 client: %w", err)
	}

	listOpts := eips.ListOpts{
		ID:             d.Get("id").(string),
		Status:         d.Get("status").(string),
		PrivateAddress: d.Get("private_ip_address").(string),
		PortID:         d.Get("port_id").(string),
		BandwidthID:    d.Get("bandwidth_id").(string),
		PublicAddress:  d.Get("public_ip_address").(string),
	}

	refinedEIPs, err := eips.List(client, listOpts)
	if err != nil {
		return fmterr.Errorf("unable to retrieve EIPs: %w", err)
	}

	var filteredEips []eips.PublicIp
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		r := regexp.MustCompile(nameRegex.(string))
		for _, ip := range refinedEIPs {
			if r.MatchString(ip.Name) {
				filteredEips = append(filteredEips, ip)
			}
		}
		refinedEIPs = filteredEips
	}

	tagRaw := d.Get("tags").(map[string]interface{})
	var refinedByTags []eips.PublicIp
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
		d.Set("bandwidth_id", elasticIP.BandwidthID),
		d.Set("bandwidth_share_type", elasticIP.BandwidthShareType),
		d.Set("bandwidth_size", elasticIP.BandwidthSize),
		d.Set("create_time", elasticIP.CreateTime),
		d.Set("ip_version", elasticIP.IpVersion),
		d.Set("port_id", elasticIP.PortID),
		d.Set("private_ip_address", elasticIP.PrivateAddress),
		d.Set("public_ip_address", elasticIP.PublicAddress),
		d.Set("tenant_id", elasticIP.TenantID),
		d.Set("region", config.GetRegion(d)),
		d.Set("name", elasticIP.Name),
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
