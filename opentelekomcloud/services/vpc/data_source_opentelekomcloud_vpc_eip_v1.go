package vpc

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/eips"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func DataSourceVPCEipV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPCEipV1Read,

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
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
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

func dataSourceVPCEipV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV1 client: %w", err)
	}

	listOpts := eips.ListOpts{
		ID:               d.Get("id").(string),
		Status:           d.Get("status").(string),
		PrivateIPAddress: d.Get("private_ip_address").(string),
		PortID:           d.Get("port_id").(string),
		BandwidthID:      d.Get("bandwidth_id").(string),
		PublicIPAddress:  d.Get("public_ip_address").(string),
	}

	refinedEIPs, err := eips.List(client, listOpts)
	if err != nil {
		return fmt.Errorf("unable to retrieve EIPs: %w", err)
	}

	tagRaw := d.Get("tags").(map[string]interface{})
	var refinedByTags []eips.PublicIp
	if len(tagRaw) > 0 {
		networkingV2Client, err := config.NetworkingV2Client(config.GetRegion(d))
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %w", err)
		}
		tagList := common.ExpandResourceTags(tagRaw)
		for _, eip := range refinedEIPs {
			resourceTagList, err := tags.Get(networkingV2Client, "", eip.ID).Extract()
			if err != nil {
				return err
			}

			if reflect.DeepEqual(tagList, resourceTagList) {
				refinedByTags = append(refinedByTags, eip)
			}
		}
	} else {
		refinedByTags = refinedEIPs
	}

	if len(refinedByTags) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(refinedByTags) > 1 {
		return fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
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
	)

	if mErr.ErrorOrNil() != nil {
		return mErr
	}

	return nil
}
