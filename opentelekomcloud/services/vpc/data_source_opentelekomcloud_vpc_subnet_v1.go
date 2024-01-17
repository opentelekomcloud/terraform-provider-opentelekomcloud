package vpc

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/subnets"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceVpcSubnetV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcSubnetV1Read,

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
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_list": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateway_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dhcp_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"primary_dns": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secondary_dns": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cidr_ipv6": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"gateway_ipv6": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVpcSubnetV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	listOpts := subnets.ListOpts{
		ID:               d.Get("id").(string),
		Name:             d.Get("name").(string),
		CIDR:             d.Get("cidr").(string),
		Status:           d.Get("status").(string),
		GatewayIP:        d.Get("gateway_ip").(string),
		PrimaryDNS:       d.Get("primary_dns").(string),
		SecondaryDNS:     d.Get("secondary_dns").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
		VpcID:            d.Get("vpc_id").(string),
	}

	refinedSubnets, err := subnets.List(client, listOpts)
	if err != nil {
		return fmterr.Errorf("unable to retrieve subnets: %w", err)
	}

	if len(refinedSubnets) == 0 {
		return fmterr.Errorf("no matching subnet found. Please change your search criteria and try again")
	}

	if len(refinedSubnets) > 1 {
		return fmterr.Errorf("multiple subnets matched; use additional constraints to reduce matches to a single subnet")
	}

	subnet := refinedSubnets[0]

	log.Printf("[INFO] Retrieved Subnet using given filter %s: %+v", subnet.ID, subnet)
	d.SetId(subnet.ID)

	mErr := multierror.Append(
		d.Set("name", subnet.Name),
		d.Set("cidr", subnet.CIDR),
		d.Set("dns_list", subnet.DNSList),
		d.Set("status", subnet.Status),
		d.Set("gateway_ip", subnet.GatewayIP),
		d.Set("dhcp_enable", subnet.EnableDHCP),
		d.Set("primary_dns", subnet.PrimaryDNS),
		d.Set("secondary_dns", subnet.SecondaryDNS),
		d.Set("availability_zone", subnet.AvailabilityZone),
		d.Set("vpc_id", subnet.VpcID),
		d.Set("subnet_id", subnet.SubnetID),
		d.Set("network_id", subnet.NetworkID),
		d.Set("ipv6_enable", subnet.EnableIpv6),
		d.Set("cidr_ipv6", subnet.CidrV6),
		d.Set("gateway_ipv6", subnet.GatewayIpV6),
		d.Set("region", config.GetRegion(d)),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}
