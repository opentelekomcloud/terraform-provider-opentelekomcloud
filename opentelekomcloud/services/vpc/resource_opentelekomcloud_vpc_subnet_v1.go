package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/subnets"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceVpcSubnetV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcSubnetV1Create,
		ReadContext:   resourceVpcSubnetV1Read,
		UpdateContext: resourceVpcSubnetV1Update,
		DeleteContext: resourceVpcSubnetV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: common.ValidateName,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cidr": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
			},
			"dns_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPv4Address,
				},
			},
			"gateway_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"dhcp_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"primary_dns": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"secondary_dns": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
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
			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVpcSubnetV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	primaryDNS := d.Get("primary_dns").(string)
	secondaryDNS := d.Get("secondary_dns").(string)
	dnsList := common.ExpandToStringSlice(d.Get("dns_list").([]interface{}))
	if primaryDNS == "" && secondaryDNS == "" && len(dnsList) == 0 {
		primaryDNS = defaultDNS[0]
		secondaryDNS = defaultDNS[1]
	}

	enableDHCP := d.Get("dhcp_enable").(bool)
	createOpts := subnets.CreateOpts{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		CIDR:             d.Get("cidr").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
		GatewayIP:        d.Get("gateway_ip").(string),
		EnableDHCP:       &enableDHCP,
		VpcID:            d.Get("vpc_id").(string),
		PrimaryDNS:       primaryDNS,
		SecondaryDNS:     secondaryDNS,
		DNSList:          dnsList,
	}

	if common.HasFilledOpt(d, "ntp_addresses") {
		var extraDhcpRequests []subnets.ExtraDHCPOpt
		extraDhcpReq := subnets.ExtraDHCPOpt{
			OptName:  "ntp",
			OptValue: d.Get("ntp_addresses").(string),
		}
		extraDhcpRequests = append(extraDhcpRequests, extraDhcpReq)
		createOpts.ExtraDHCPOpts = extraDhcpRequests
	}

	subnet, err := subnets.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC subnet: %w", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForVpcSubnetActive(client, subnet.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for Subnet (%s) to become ACTIVE: %w", subnet.ID, err)
	}

	d.SetId(subnet.ID)

	if err := addNetworkingTags(d, config, "subnets"); err != nil {
		return diag.FromErr(err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceVpcSubnetV1Read(clientCtx, d, config)
}

func resourceVpcSubnetV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	subnet, err := subnets.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "vpc subnet")
	}

	mErr := multierror.Append(
		d.Set("name", subnet.Name),
		d.Set("description", subnet.Description),
		d.Set("cidr", subnet.CIDR),
		d.Set("dns_list", subnet.DNSList),
		d.Set("gateway_ip", subnet.GatewayIP),
		d.Set("dhcp_enable", subnet.EnableDHCP),
		d.Set("primary_dns", subnet.PrimaryDNS),
		d.Set("secondary_dns", subnet.SecondaryDNS),
		d.Set("availability_zone", subnet.AvailabilityZone),
		d.Set("vpc_id", subnet.VpcID),
		d.Set("subnet_id", subnet.SubnetID),
		d.Set("network_id", subnet.NetworkID),
		d.Set("status", subnet.Status),
		d.Set("region", config.GetRegion(d)),
	)

	for _, opt := range subnet.ExtraDHCPOpts {
		if opt.OptName == "ntp" {
			mErr = multierror.Append(mErr, d.Set("ntp_addresses", opt.OptValue))
			break
		}
	}

	if mErr.ErrorOrNil() != nil {
		return fmterr.Errorf("error setting subnet fields: %w", mErr)
	}

	if err := readNetworkingTags(d, config, "subnets"); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVpcSubnetV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	var updateOpts subnets.UpdateOpts

	// as name is mandatory while updating subnet
	updateOpts.Name = d.Get("name").(string)

	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("primary_dns") {
		updateOpts.PrimaryDNS = d.Get("primary_dns").(string)
	}
	if d.HasChange("secondary_dns") {
		updateOpts.SecondaryDNS = d.Get("secondary_dns").(string)
	}
	if d.HasChange("dns_list") {
		updateOpts.DNSList = common.ExpandToStringSlice(d.Get("dns_list").([]interface{}))
	}
	if d.HasChange("dhcp_enable") {
		enableDHCP := d.Get("dhcp_enable").(bool)
		updateOpts.EnableDHCP = &enableDHCP
	}
	if d.HasChange("ntp_addresses") {
		var extraDhcpRequests []subnets.ExtraDHCPOpt
		extraDhcpReq := subnets.ExtraDHCPOpt{
			OptName:  "ntp",
			OptValue: d.Get("ntp_addresses").(string),
		}
		extraDhcpRequests = append(extraDhcpRequests, extraDhcpReq)
		updateOpts.ExtraDhcpOpts = extraDhcpRequests
	}

	vpcID := d.Get("vpc_id").(string)

	_, err = subnets.Update(client, vpcID, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud VPC Subnet: %w", err)
	}

	// update tags
	if d.HasChange("tags") {
		networkingV2Client, err := config.NetworkingV2Client(config.GetRegion(d))
		if err != nil {
			return fmterr.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
		}

		if err := common.UpdateResourceTags(networkingV2Client, d, "subnets", d.Id()); err != nil {
			return fmterr.Errorf("error updating tags of VPC subnet %s: %w", d.Id(), err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceVpcSubnetV1Read(clientCtx, d, meta)
}

func resourceVpcSubnetV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	vpcID := d.Get("vpc_id").(string)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForVpcSubnetDelete(client, vpcID, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud Subnet: %w", err)
	}

	d.SetId("")
	return nil
}

func waitForVpcSubnetActive(client *golangsdk.ServiceClient, subnetID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		subnet, err := subnets.Get(client, subnetID).Extract()
		if err != nil {
			return nil, "", err
		}

		if subnet.Status == "ACTIVE" {
			return subnet, "ACTIVE", nil
		}

		// If subnet status is other than Active, send error
		if subnet.Status == "DOWN" || subnet.Status == "error" {
			return nil, "", fmt.Errorf("subnet status: %s", subnet.Status)
		}

		return subnet, "CREATING", nil
	}
}

func waitForVpcSubnetDelete(client *golangsdk.ServiceClient, vpcID string, subnetID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		subnet, err := subnets.Get(client, subnetID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud subnet %s", subnetID)
				return subnet, "DELETED", nil
			}
			return subnet, "ACTIVE", err
		}

		if err := subnets.Delete(client, vpcID, subnetID).ExtractErr(); err != nil {
			switch err.(type) {
			case golangsdk.ErrDefault404, golangsdk.ErrDefault400:
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud subnet %s", subnetID)
				return subnet, "DELETED", nil
			case golangsdk.ErrDefault409, golangsdk.ErrDefault500:
				return subnet, "ACTIVE", nil
			default:
				return subnet, "ACTIVE", err
			}
		}

		return subnet, "ACTIVE", nil
	}
}
