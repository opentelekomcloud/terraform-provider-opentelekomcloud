package vpc

import (
	"context"
	"sort"

	"github.com/hashicorp/go-multierror"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/subnets"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/networkipavailabilities"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceVpcSubnetIdsV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcSubnetIdsV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func dataSourceVpcSubnetIdsV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud NetworkingV1 client: %w", err)
	}

	vpcID := d.Get("vpc_id").(string)
	listOpts := subnets.ListOpts{
		VpcID: vpcID,
	}

	refinedSubnets, err := subnets.List(client, listOpts)
	if err != nil {
		return fmterr.Errorf("unable to retrieve subnets: %w", err)
	}

	if len(refinedSubnets) == 0 {
		return fmterr.Errorf("no matching subnet found for vpc with id %s", vpcID)
	}

	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %w", err)
	}

	sortedSubnets := make([]SubnetIP, 0)
	for _, subnet := range refinedSubnets {
		net, err := networkipavailabilities.Get(networkingClient, subnet.ID).Extract()
		if err != nil {
			return fmterr.Errorf("error retrieving NetworkIP availabilities: %w", err)
		}
		subnetIPAvail := net.SubnetIPAvailabilities[0]
		newSubnet := SubnetIP{
			ID:  subnet.ID,
			IPs: subnetIPAvail.TotalIPs - subnetIPAvail.UsedIPs,
		}
		sortedSubnets = append(sortedSubnets, newSubnet)
	}

	// Returns the Subnet contains most available IPs out of a slice of subnets.
	sort.Sort(sort.Reverse(subnetSort(sortedSubnets)))
	subnetIDs := make([]string, 0)
	for _, subnet := range sortedSubnets {
		subnetIDs = append(subnetIDs, subnet.ID)
	}

	d.SetId(vpcID)
	mErr := multierror.Append(
		d.Set("ids", subnetIDs),
		d.Set("region", config.GetRegion(d)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

type SubnetIP struct {
	ID  string
	IPs int
}

type subnetSort []SubnetIP

func (a subnetSort) Len() int      { return len(a) }
func (a subnetSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a subnetSort) Less(i, j int) bool {
	return a[i].IPs < a[j].IPs
}
