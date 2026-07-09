package vpc

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v3/security/rules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceVpcSecGroupRuleV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcSecGroupRuleV3Create,
		ReadContext:   resourceVpcSecGroupRuleV3Read,
		DeleteContext: resourceVpcSecGroupRuleV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ether_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"multi_port": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			},
			"remote_address_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVpcSecGroupRuleV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	opts := rules.CreateOpts{
		SecurityGroupRule: rules.SecurityGroupRuleOptions{
			SecurityGroupID:      d.Get("security_group_id").(string),
			Description:          d.Get("description").(string),
			Direction:            d.Get("direction").(string),
			Ethertype:            d.Get("ether_type").(string),
			Protocol:             d.Get("protocol").(string),
			Multiport:            d.Get("multi_port").(string),
			RemoteGroupID:        d.Get("remote_group_id").(string),
			RemoteIPPrefix:       d.Get("remote_ip_prefix").(string),
			RemoteAddressGroupID: d.Get("remote_address_group_id").(string),
			Action:               d.Get("action").(string),
			Priority:             d.Get("priority").(int),
		},
	}

	securityGroupRule, err := rules.Create(client, opts)
	if err != nil {
		return fmterr.Errorf("error creating security group rule: %w", err)
	}
	d.SetId(securityGroupRule.ID)

	log.Printf("[DEBUG] OpenTelekomCloud Security Group Rule `%s` created: %#v", d.Id(), securityGroupRule)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceVpcSecGroupRuleV3Read(clientCtx, d, meta)
}

func resourceVpcSecGroupRuleV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	securityGroupRule, err := rules.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error fetching OpenTelekomCloud VPC Security Group Rule V3")
	}

	mErr := multierror.Append(
		d.Set("security_group_id", securityGroupRule.SecurityGroupID),
		d.Set("description", securityGroupRule.Description),
		d.Set("direction", securityGroupRule.Direction),
		d.Set("protocol", securityGroupRule.Protocol),
		d.Set("ether_type", securityGroupRule.Ethertype),
		d.Set("multi_port", securityGroupRule.Multiport),
		d.Set("action", securityGroupRule.Action),
		d.Set("priority", securityGroupRule.Priority),
		d.Set("created_at", securityGroupRule.CreatedAt),
		d.Set("updated_at", securityGroupRule.UpdatedAt),
		d.Set("project_id", securityGroupRule.ProjectID),
		d.Set("remote_group_id", securityGroupRule.RemoteGroupID),
		d.Set("remote_ip_prefix", securityGroupRule.RemoteIPPrefix),
		d.Set("remote_address_group_id", securityGroupRule.RemoteAddressGroupID),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVpcSecGroupRuleV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	err = rules.Delete(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting security group rule `%s` : %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}
