package vpc

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/security/rules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceNetworkingSecGroupRuleV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingSecGroupRuleV2Create,
		ReadContext:   resourceNetworkingSecGroupRuleV2Read,
		DeleteContext: resourceNetworkingSecGroupRuleV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: false,
				ForceNew: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ethertype": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port_range_min": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"port_range_max": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"remote_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"remote_ip_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				StateFunc: func(v interface{}) string {
					return strings.ToLower(v.(string))
				},
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceNetworkingSecGroupRuleV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	portRangeMin := d.Get("port_range_min").(int)
	portRangeMax := d.Get("port_range_max").(int)
	protocol := d.Get("protocol").(string)

	if protocol == "" {
		if portRangeMin != 0 || portRangeMax != 0 {
			return fmterr.Errorf("A protocol must be specified when using port_range_min and port_range_max")
		}
	}

	opts := rules.CreateOpts{
		Description:    d.Get("description").(string),
		SecGroupID:     d.Get("security_group_id").(string),
		PortRangeMin:   &portRangeMin,
		PortRangeMax:   &portRangeMax,
		RemoteGroupID:  d.Get("remote_group_id").(string),
		RemoteIPPrefix: d.Get("remote_ip_prefix").(string),
		TenantID:       d.Get("tenant_id").(string),
	}

	if v, ok := d.GetOk("direction"); ok {
		direction := resourceNetworkingSecGroupRuleV2DetermineDirection(v.(string))
		opts.Direction = direction
	}

	if v, ok := d.GetOk("ethertype"); ok {
		etherType := resourceNetworkingSecGroupRuleV2DetermineEtherType(v.(string))
		opts.EtherType = etherType
	}

	if v, ok := d.GetOk("protocol"); ok {
		protocol := resourceNetworkingSecGroupRuleV2DetermineProtocol(v.(string))
		opts.Protocol = protocol
	}

	log.Printf("[DEBUG] Create OpenTelekomCloud Neutron security group: %#v", opts)

	securityGroupRule, err := rules.Create(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] OpenTelekomCloud Neutron Security Group Rule created: %#v", securityGroupRule)

	d.SetId(securityGroupRule.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceNetworkingSecGroupRuleV2Read(clientCtx, d, meta)
}

func resourceNetworkingSecGroupRuleV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	securityGroupRule, err := rules.Get(client, d.Id()).Extract()

	if err != nil {
		return common.CheckDeletedDiag(d, err, "OpenTelekomCloud Security Group Rule")
	}

	mErr := multierror.Append(
		d.Set("description", securityGroupRule.Description),
		d.Set("direction", securityGroupRule.Direction),
		d.Set("ethertype", securityGroupRule.EtherType),
		d.Set("protocol", securityGroupRule.Protocol),
		d.Set("port_range_min", securityGroupRule.PortRangeMin),
		d.Set("port_range_max", securityGroupRule.PortRangeMax),
		d.Set("remote_group_id", securityGroupRule.RemoteGroupID),
		d.Set("remote_ip_prefix", securityGroupRule.RemoteIPPrefix),
		d.Set("security_group_id", securityGroupRule.SecGroupID),
		d.Set("tenant_id", securityGroupRule.TenantID),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetworkingSecGroupRuleV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForSecGroupRuleDelete(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud Neutron Security Group Rule: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceNetworkingSecGroupRuleV2DetermineDirection(v string) rules.RuleDirection {
	var direction rules.RuleDirection
	switch v {
	case "ingress":
		direction = rules.DirIngress
	case "egress":
		direction = rules.DirEgress
	}

	return direction
}

func resourceNetworkingSecGroupRuleV2DetermineEtherType(v string) rules.RuleEtherType {
	var etherType rules.RuleEtherType
	switch v {
	case "IPv4":
		etherType = rules.EtherType4
	case "IPv6":
		etherType = rules.EtherType6
	}

	return etherType
}

func resourceNetworkingSecGroupRuleV2DetermineProtocol(v string) rules.RuleProtocol {
	var protocol rules.RuleProtocol

	// Check and see if the requested protocol matched a list of known protocol names.
	switch v {
	case "tcp":
		protocol = rules.ProtocolTCP
	case "udp":
		protocol = rules.ProtocolUDP
	case "icmp":
		protocol = rules.ProtocolICMP
	case "ah":
		protocol = rules.ProtocolAH
	case "dccp":
		protocol = rules.ProtocolDCCP
	case "egp":
		protocol = rules.ProtocolEGP
	case "esp":
		protocol = rules.ProtocolESP
	case "gre":
		protocol = rules.ProtocolGRE
	case "igmp":
		protocol = rules.ProtocolIGMP
	case "ipv6-encap":
		protocol = rules.ProtocolIPv6Encap
	case "ipv6-frag":
		protocol = rules.ProtocolIPv6Frag
	case "ipv6-icmp":
		protocol = rules.ProtocolIPv6ICMP
	case "ipv6-nonxt":
		protocol = rules.ProtocolIPv6NoNxt
	case "ipv6-opts":
		protocol = rules.ProtocolIPv6Opts
	case "ipv6-route":
		protocol = rules.ProtocolIPv6Route
	case "ospf":
		protocol = rules.ProtocolOSPF
	case "pgm":
		protocol = rules.ProtocolPGM
	case "rsvp":
		protocol = rules.ProtocolRSVP
	case "sctp":
		protocol = rules.ProtocolSCTP
	case "udplite":
		protocol = rules.ProtocolUDPLite
	case "vrrp":
		protocol = rules.ProtocolVRRP
	}

	// If the protocol wasn't matched above, see if it's an integer.
	if protocol == "" {
		_, err := strconv.Atoi(v)
		if err == nil {
			protocol = rules.RuleProtocol(v)
		}
	}

	return protocol
}

func waitForSecGroupRuleDelete(client *golangsdk.ServiceClient, secGroupRuleId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete OpenTelekomCloud Security Group Rule %s.\n", secGroupRuleId)

		r, err := rules.Get(client, secGroupRuleId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud Neutron Security Group Rule %s", secGroupRuleId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = rules.Delete(client, secGroupRuleId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud Neutron Security Group Rule %s", secGroupRuleId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud Neutron Security Group Rule %s still active.\n", secGroupRuleId)
		return r, "ACTIVE", nil
	}
}
