package vpc

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/subnets"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func resourceSubnetDNSListV1(d *schema.ResourceData) []string {
	rawDNS := d.Get("dns_list").([]interface{})
	dns := make([]string, len(rawDNS))
	for i, raw := range rawDNS {
		dns[i] = raw.(string)
	}
	return dns
}

func ResourceVpcSubnetV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcSubnetV1Create,
		Read:   resourceVpcSubnetV1Read,
		Update: resourceVpcSubnetV1Update,
		Delete: resourceVpcSubnetV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ // request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: common.ValidateName,
			},
			"cidr": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: common.ValidateCIDR,
			},
			"dns_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: common.ValidateIP,
				},
			},
			"gateway_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: common.ValidateIP,
			},
			"dhcp_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: false,
			},
			"primary_dns": {
				Type:         schema.TypeString,
				ForceNew:     false,
				Optional:     true,
				ValidateFunc: common.ValidateIP,
				Default:      defaultDNS[0],
			},
			"secondary_dns": {
				Type:         schema.TypeString,
				ForceNew:     false,
				Optional:     true,
				ValidateFunc: common.ValidateIP,
				Default:      defaultDNS[1],
			},
			"availability_zone": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"tags": common.TagsSchema(),
			"ntp_addresses": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVpcSubnetV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	subnetClient, err := config.NetworkingV1Client(config.GetRegion(d))

	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	createOpts := subnets.CreateOpts{
		Name:             d.Get("name").(string),
		CIDR:             d.Get("cidr").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
		GatewayIP:        d.Get("gateway_ip").(string),
		EnableDHCP:       d.Get("dhcp_enable").(bool),
		VPC_ID:           d.Get("vpc_id").(string),
		PRIMARY_DNS:      d.Get("primary_dns").(string),
		SECONDARY_DNS:    d.Get("secondary_dns").(string),
		DnsList:          resourceSubnetDNSListV1(d),
	}

	if common.HasFilledOpt(d, "ntp_addresses") {
		var extraDhcpRequests []subnets.ExtraDhcpOpt
		extraDhcpReq := subnets.ExtraDhcpOpt{
			OptName:  "ntp",
			OptValue: d.Get("ntp_addresses").(string),
		}
		extraDhcpRequests = append(extraDhcpRequests, extraDhcpReq)
		createOpts.ExtraDhcpOpts = extraDhcpRequests
	}

	n, err := subnets.Create(subnetClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud VPC subnet: %s", err)
	}

	d.SetId(n.ID)
	log.Printf("[INFO] Vpc Subnet ID: %s", n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForVpcSubnetActive(subnetClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf(
			"error waiting for Subnet (%s) to become ACTIVE: %s",
			n.ID, stateErr)
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		vpcSubnetV2Client, err := config.NetworkingV2Client(config.GetRegion(d))
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
		}

		taglist := common.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(vpcSubnetV2Client, "subnets", n.ID, taglist).ExtractErr(); tagErr != nil {
			return fmt.Errorf("error setting tags of VpcSubnet %s: %s", n.ID, tagErr)
		}
	}

	return resourceVpcSubnetV1Read(d, config)

}

func resourceVpcSubnetV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	subnetClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	n, err := subnets.Get(subnetClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("error retrieving OpenTelekomCloud Subnets: %s", err)
	}

	mErr := multierror.Append(
		d.Set("name", n.Name),
		d.Set("cidr", n.CIDR),
		d.Set("dns_list", n.DnsList),
		d.Set("gateway_ip", n.GatewayIP),
		d.Set("dhcp_enable", n.EnableDHCP),
		d.Set("primary_dns", n.PRIMARY_DNS),
		d.Set("secondary_dns", n.SECONDARY_DNS),
		d.Set("availability_zone", n.AvailabilityZone),
		d.Set("vpc_id", n.VPC_ID),
		d.Set("subnet_id", n.SubnetId),
		d.Set("region", config.GetRegion(d)),
	)

	for _, opt := range n.ExtraDhcpOpts {
		if opt.OptName == "ntp" {
			mErr = multierror.Append(mErr, d.Set("ntp_addresses", opt.OptValue))
			break
		}
	}

	if mErr.ErrorOrNil() != nil {
		return fmt.Errorf("error setting subnet fields: %w", mErr)
	}

	// save VpcSubnet tags
	vpcSubnetV2Client, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	resourceTags, err := tags.Get(vpcSubnetV2Client, "subnets", d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("error fetching OpenTelekomCloud VpcSubnet tags: %s", err)
	}

	tagmap := common.TagsToMap(resourceTags)
	if err := d.Set("tags", tagmap); err != nil {
		return fmt.Errorf("error saving tags for OpenTelekomCloud VpcSubnet %s: %s", d.Id(), err)
	}

	return nil
}

func resourceVpcSubnetV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	subnetClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	var updateOpts subnets.UpdateOpts

	// as name is mandatory while updating subnet
	updateOpts.Name = d.Get("name").(string)

	if d.HasChange("primary_dns") {
		updateOpts.PRIMARY_DNS = d.Get("primary_dns").(string)
	}
	if d.HasChange("secondary_dns") {
		updateOpts.SECONDARY_DNS = d.Get("secondary_dns").(string)
	}
	if d.HasChange("dns_list") {
		updateOpts.DnsList = resourceSubnetDNSListV1(d)
	}
	if d.HasChange("dhcp_enable") {
		updateOpts.EnableDHCP = d.Get("dhcp_enable").(bool)

	} else if d.Get("dhcp_enable").(bool) { // maintaining dhcp to be true if it was true earlier as default update option for dhcp bool is always going to be false in golangsdk
		updateOpts.EnableDHCP = true
	}
	if d.HasChange("ntp_addresses") {
		var extraDhcpRequests []subnets.ExtraDhcpOpt
		extraDhcpReq := subnets.ExtraDhcpOpt{
			OptName:  "ntp",
			OptValue: d.Get("ntp_addresses").(string),
		}
		extraDhcpRequests = append(extraDhcpRequests, extraDhcpReq)
		updateOpts.ExtraDhcpOpts = extraDhcpRequests
	}

	vpc_id := d.Get("vpc_id").(string)

	_, err = subnets.Update(subnetClient, vpc_id, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("error updating OpenTelekomCloud VPC Subnet: %s", err)
	}

	// update tags
	if d.HasChange("tags") {
		vpcSubnetV2Client, err := config.NetworkingV2Client(config.GetRegion(d))
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
		}

		tagErr := common.UpdateResourceTags(vpcSubnetV2Client, d, "subnets", d.Id())
		if tagErr != nil {
			return fmt.Errorf("error updating tags of VPC subnet %s: %s", d.Id(), tagErr)
		}
	}

	return resourceVpcSubnetV1Read(d, meta)
}

func resourceVpcSubnetV1Delete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*cfg.Config)
	subnetClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}
	vpc_id := d.Get("vpc_id").(string)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForVpcSubnetDelete(subnetClient, vpc_id, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error deleting OpenTelekomCloud Subnet: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForVpcSubnetActive(subnetClient *golangsdk.ServiceClient, vpcId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := subnets.Get(subnetClient, vpcId).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.Status == "ACTIVE" {
			return n, "ACTIVE", nil
		}

		// If subnet status is other than Active, send error
		if n.Status == "DOWN" || n.Status == "error" {
			return nil, "", fmt.Errorf("Subnet status: '%s'", n.Status)
		}

		return n, "CREATING", nil
	}
}

func waitForVpcSubnetDelete(subnetClient *golangsdk.ServiceClient, vpcId string, subnetId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := subnets.Get(subnetClient, subnetId).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud subnet %s", subnetId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}
		err = subnets.Delete(subnetClient, vpcId, subnetId).ExtractErr()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud subnet %s", subnetId)
				return r, "DELETED", nil
			}
			if _, ok := err.(golangsdk.ErrDefault400); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud subnet %s", subnetId)
				return r, "DELETED", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 || errCode.Actual == 500 {
					return r, "ACTIVE", nil
				}
			}
			return r, "ACTIVE", err
		}

		return r, "ACTIVE", nil
	}
}
