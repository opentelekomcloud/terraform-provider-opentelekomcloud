package iam

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3.0/acl"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func ResourceIdentityAclV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityACLV3Create,
		ReadContext:   resourceIdentityACLV3Read,
		UpdateContext: resourceIdentityACLV3Update,
		DeleteContext: resourceIdentityACLV3Delete,

		CustomizeDiff: resourceIdentityACLV3CustomizeDiff,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"console", "api",
				}, true),
			},
			"ip_cidrs": {
				Type:         schema.TypeSet,
				Optional:     true,
				MaxItems:     200,
				AtLeastOneOf: []string{"ip_ranges"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: common.ValidateCIDR,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: resourceACLPolicyCIDRHash,
			},
			"ip_ranges": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 200,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"range": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: common.ValidateIPRange,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: resourceACLPolicyRangeHash,
			},
			"ipv6_cidrs": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 200,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: common.ValidateIPv6CIDR,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: resourceACLPolicyCIDRHash,
			},
			"ipv6_ranges": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 200,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"range": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: common.ValidateIPv6Range,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: resourceACLPolicyRangeHash,
			},
		},
	}
}

func resourceIdentityACLV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	iamClient, err := config.IdentityV30AdminClient()
	if err != nil {
		return fmterr.Errorf("error creating client: %s", err)
	}
	domainID, err := getDomainID(config, iamClient)
	if err != nil {
		return fmterr.Errorf("error getting the domain id, err=%s", err)
	}
	d.SetId(domainID)

	if err := updateACLPolicy(d, meta); err != nil {
		return diag.Errorf("error creating identity ACL: %s", err)
	}

	return resourceIdentityACLV3Read(ctx, d, meta)
}

func resourceIdentityACLV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mErr := &multierror.Error{}
	config := meta.(*cfg.Config)
	iamClient, err := config.IdentityV30AdminClient()
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	var res *acl.ACLPolicy
	switch d.Get("type").(string) {
	case "console":
		res, err = acl.ConsoleACLPolicyGet(iamClient, d.Id())
		if err != nil {
			return diag.Errorf("error fetching identity ACL for console access")
		}
	default:
		res, err = acl.APIACLPolicyGet(iamClient, d.Id())
		if err != nil {
			return diag.Errorf("error fetching identity ACL for API access")
		}
	}

	log.Printf("[DEBUG] Retrieved identity ACL: %#v", res)
	if len(res.AllowAddressNetmasks) > 0 {
		addressNetmasks := make([]map[string]string, 0, len(res.AllowAddressNetmasks))
		for _, v := range res.AllowAddressNetmasks {
			addressNetmask := map[string]string{
				"cidr":        v.AddressNetmask,
				"description": v.Description,
			}
			addressNetmasks = append(addressNetmasks, addressNetmask)
		}
		mErr = multierror.Append(mErr, d.Set("ip_cidrs", addressNetmasks))
	}
	if len(res.AllowIPRanges) > 0 {
		ipRanges := make([]map[string]string, 0, len(res.AllowIPRanges))
		for _, v := range res.AllowIPRanges {
			ipRange := map[string]string{
				"range":       v.IPRange,
				"description": v.Description,
			}
			ipRanges = append(ipRanges, ipRange)
		}
		mErr = multierror.Append(mErr, d.Set("ip_ranges", ipRanges))
	}
	if len(res.AllowAddressNetmasksIPv6) > 0 {
		addressNetmasks := make([]map[string]string, 0, len(res.AllowAddressNetmasksIPv6))
		for _, v := range res.AllowAddressNetmasksIPv6 {
			addressNetmask := map[string]string{
				"cidr":        v.AddressNetmask,
				"description": v.Description,
			}
			addressNetmasks = append(addressNetmasks, addressNetmask)
		}
		mErr = multierror.Append(mErr, d.Set("ipv6_cidrs", addressNetmasks))
	}
	if len(res.AllowIPRangesIPv6) > 0 {
		ipRanges := make([]map[string]string, 0, len(res.AllowIPRangesIPv6))
		for _, v := range res.AllowIPRangesIPv6 {
			ipRange := map[string]string{
				"range":       v.IPRange,
				"description": v.Description,
			}
			ipRanges = append(ipRanges, ipRange)
		}
		mErr = multierror.Append(mErr, d.Set("ipv6_ranges", ipRanges))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting identity ACL fields: %s", err)
	}

	return nil
}

func resourceIdentityACLV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChanges("ip_cidrs", "ip_ranges", "ipv6_cidrs", "ipv6_ranges") {
		if err := updateACLPolicy(d, meta); err != nil {
			return diag.Errorf("error updating identity ACL: %s", err)
		}
	}

	return resourceIdentityACLV3Read(ctx, d, meta)
}

func resourceIdentityACLV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	iamClient, err := config.IdentityV30AdminClient()
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	netmasksList := make([]acl.AllowAddressNetmasks, 0, 1)
	netmask := acl.AllowAddressNetmasks{
		AddressNetmask: "0.0.0.0-255.255.255.255",
	}
	netmasksList = append(netmasksList, netmask)

	deleteOpts := acl.ACLPolicy{
		AllowAddressNetmasks: netmasksList,
		DomainId:             d.Id(),
	}
	if d.Get("type").(string) == "console" {
		deleteOpts.AllowAddressNetmasksIPv6 = []acl.AllowAddressNetmasks{
			{
				AddressNetmask: "0000:0000:0000:0000:0000:0000:0000:0000-FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF",
			},
		}
	}

	switch d.Get("type").(string) {
	case "console":
		_, err = acl.ConsoleACLPolicyUpdate(iamClient, deleteOpts)
	default:
		_, err = acl.APIACLPolicyUpdate(iamClient, deleteOpts)
	}

	if err != nil {
		return diag.Errorf("error resetting identity ACL: %s", err)
	}

	return nil
}

func updateACLPolicy(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	iamClient, err := config.IdentityV30AdminClient()
	if err != nil {
		return fmt.Errorf("error creating IAM client: %s", err)
	}

	t := d.Get("type").(string)
	if t != "console" {
		if d.Get("ipv6_cidrs").(*schema.Set).Len() > 0 || d.Get("ipv6_ranges").(*schema.Set).Len() > 0 {
			return fmt.Errorf(`ipv6_cidrs and ipv6_ranges can only be set when type = "console"`)
		}
	}

	updateOpts := acl.ACLPolicy{
		DomainId: d.Id(),
	}
	if addressNetmasks, ok := d.GetOk("ip_cidrs"); ok {
		netmasksList := make([]acl.AllowAddressNetmasks, 0, addressNetmasks.(*schema.Set).Len())
		for _, v := range addressNetmasks.(*schema.Set).List() {
			netmask := acl.AllowAddressNetmasks{
				AddressNetmask: v.(map[string]interface{})["cidr"].(string),
				Description:    v.(map[string]interface{})["description"].(string),
			}
			netmasksList = append(netmasksList, netmask)
		}
		updateOpts.AllowAddressNetmasks = netmasksList
	}
	if ipRanges, ok := d.GetOk("ip_ranges"); ok {
		rangeList := make([]acl.AllowIPRanges, 0, ipRanges.(*schema.Set).Len())
		for _, v := range ipRanges.(*schema.Set).List() {
			ipRange := acl.AllowIPRanges{
				IPRange:     v.(map[string]interface{})["range"].(string),
				Description: v.(map[string]interface{})["description"].(string),
			}
			rangeList = append(rangeList, ipRange)
		}
		updateOpts.AllowIPRanges = rangeList
	}
	if addressNetmasksv6, ok := d.GetOk("ipv6_cidrs"); ok {
		netmasksList := make([]acl.AllowAddressNetmasks, 0, addressNetmasksv6.(*schema.Set).Len())
		for _, v := range addressNetmasksv6.(*schema.Set).List() {
			netmask := acl.AllowAddressNetmasks{
				AddressNetmask: v.(map[string]interface{})["cidr"].(string),
				Description:    v.(map[string]interface{})["description"].(string),
			}
			netmasksList = append(netmasksList, netmask)
		}
		updateOpts.AllowAddressNetmasksIPv6 = netmasksList
	}
	if ipRangesv6, ok := d.GetOk("ipv6_ranges"); ok {
		rangeList := make([]acl.AllowIPRanges, 0, ipRangesv6.(*schema.Set).Len())
		for _, v := range ipRangesv6.(*schema.Set).List() {
			ipRange := acl.AllowIPRanges{
				IPRange:     v.(map[string]interface{})["range"].(string),
				Description: v.(map[string]interface{})["description"].(string),
			}
			rangeList = append(rangeList, ipRange)
		}
		updateOpts.AllowIPRangesIPv6 = rangeList
	}

	switch d.Get("type").(string) {
	case "console":
		_, err = acl.ConsoleACLPolicyUpdate(iamClient, updateOpts)
	case "api":
		_, err = acl.APIACLPolicyUpdate(iamClient, updateOpts)
	}

	return err
}

func resourceACLPolicyCIDRHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if m["cidr"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["cidr"].(string)))
	}

	return hashcode.String(buf.String())
}

func resourceACLPolicyRangeHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if m["range"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["range"].(string)))
	}

	return hashcode.String(buf.String())
}

func resourceIdentityACLV3CustomizeDiff(
	ctx context.Context,
	d *schema.ResourceDiff,
	meta interface{},
) error {
	if d.Get("type").(string) != "console" {
		if d.Get("ipv6_cidrs").(*schema.Set).Len() > 0 {
			return fmt.Errorf(`ipv6_cidrs can only be specified when type = "console"`)
		}

		if d.Get("ipv6_ranges").(*schema.Set).Len() > 0 {
			return fmt.Errorf(`ipv6_ranges can only be specified when type = "console"`)
		}
	}

	return nil
}
