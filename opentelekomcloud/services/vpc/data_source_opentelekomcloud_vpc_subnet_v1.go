package vpc

import (
	"context"
	"log"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/subnets"

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
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 255),
					validation.StringMatch(regexp.MustCompile("^[^<>]*$"),
						"description cannot contain angle brackets"),
				),
			},
			"cidr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"gateway_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dhcp_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"primary_dns": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"secondary_dns": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"subnet_id_v6": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ntp_addresses": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tenant_id": {
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
		},
	}
}

func dataSourceVpcSubnetV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.VpcV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v1 client: %w", err)
	}

	// The API can repeat a page when vpc_id and marker are combined, so apply all filters client-side.
	allSubnets, err := subnets.List(client, subnets.ListOpts{})
	if err != nil {
		return fmterr.Errorf("unable to retrieve subnets: %w", err)
	}
	refinedSubnets := filterSubnets(allSubnets, subnetFilters{
		ID:               d.Get("id").(string),
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		CIDR:             d.Get("cidr").(string),
		Status:           d.Get("status").(string),
		GatewayIP:        d.Get("gateway_ip").(string),
		PrimaryDNS:       d.Get("primary_dns").(string),
		SecondaryDNS:     d.Get("secondary_dns").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
		NtpAddresses:     d.Get("ntp_addresses").(string),
		Scope:            d.Get("scope").(string),
		VpcID:            d.Get("vpc_id").(string),
	})

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
		d.Set("description", subnet.Description),
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
		d.Set("subnet_id_v6", subnet.SubnetIDV6),
		d.Set("network_id", subnet.NetworkID),
		d.Set("ipv6_enable", subnet.EnableIpv6),
		d.Set("cidr_ipv6", subnet.CidrV6),
		d.Set("gateway_ipv6", subnet.GatewayIpV6),
		d.Set("ntp_addresses", subnetNtpAddresses(subnet.ExtraDHCPOpts)),
		d.Set("scope", subnet.Scope),
		d.Set("tenant_id", subnet.TenantID),
		d.Set("created_at", subnet.CreatedAt),
		d.Set("updated_at", subnet.UpdatedAt),
		d.Set("region", config.GetRegion(d)),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

type subnetFilters struct {
	ID               string
	Name             string
	Description      string
	CIDR             string
	Status           string
	GatewayIP        string
	PrimaryDNS       string
	SecondaryDNS     string
	AvailabilityZone string
	NtpAddresses     string
	Scope            string
	VpcID            string
}

func filterSubnets(allSubnets []subnets.Subnet, filters subnetFilters) []subnets.Subnet {
	refined := make([]subnets.Subnet, 0, len(allSubnets))
	for _, subnet := range allSubnets {
		if filters.ID != "" && subnet.ID != filters.ID {
			continue
		}
		if filters.Name != "" && subnet.Name != filters.Name {
			continue
		}
		if filters.Description != "" && subnet.Description != filters.Description {
			continue
		}
		if filters.CIDR != "" && subnet.CIDR != filters.CIDR {
			continue
		}
		if filters.Status != "" && subnet.Status != filters.Status {
			continue
		}
		if filters.GatewayIP != "" && subnet.GatewayIP != filters.GatewayIP {
			continue
		}
		if filters.PrimaryDNS != "" && subnet.PrimaryDNS != filters.PrimaryDNS {
			continue
		}
		if filters.SecondaryDNS != "" && subnet.SecondaryDNS != filters.SecondaryDNS {
			continue
		}
		if filters.AvailabilityZone != "" && subnet.AvailabilityZone != filters.AvailabilityZone {
			continue
		}
		if filters.NtpAddresses != "" && subnetNtpAddresses(subnet.ExtraDHCPOpts) != filters.NtpAddresses {
			continue
		}
		if filters.Scope != "" && subnet.Scope != filters.Scope {
			continue
		}
		if filters.VpcID != "" && subnet.VpcID != filters.VpcID {
			continue
		}
		refined = append(refined, subnet)
	}
	return refined
}
