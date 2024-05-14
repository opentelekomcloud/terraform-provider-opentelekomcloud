package vpc

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/ports"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceNetworkingPortSecGroupAssociateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingPortSecGroupAssociateV2Create,
		ReadContext:   resourceNetworkingPortSecGroupAssociateV2Read,
		UpdateContext: resourceNetworkingPortSecGroupAssociateV2Update,
		DeleteContext: resourceNetworkingPortSecGroupAssociateV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"port_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"all_security_group_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceNetworkingPortSecGroupAssociateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	securityGroups := common.ExpandToStringSlice(d.Get("security_group_ids").(*schema.Set).List())
	portID := d.Get("port_id").(string)

	port, err := ports.Get(client, portID).Extract()
	if err != nil {
		return diag.Errorf("unable to get OpenTelekomCloud %s Port: %s", portID, err)
	}

	log.Printf("[DEBUG] Retrieved Port %s: %+v", portID, port)

	var updateOpts ports.UpdateOpts
	var force bool
	if v, ok := d.GetOk("force"); ok {
		force = v.(bool)
	}

	if force {
		updateOpts.SecurityGroups = &securityGroups
	} else {
		sg := common.SliceUnion(port.SecurityGroups, securityGroups)
		updateOpts.SecurityGroups = &sg
	}

	log.Printf("[DEBUG] Port Security Group Associate Options: %#v", updateOpts.SecurityGroups)

	_, err = ports.Update(client, portID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error associating OpenTelekomCloud %s port with '%s' security groups: %s", portID, strings.Join(securityGroups, ","), err)
	}

	d.SetId(portID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceNetworkingPortSecGroupAssociateV2Read(clientCtx, d, meta)
}

func resourceNetworkingPortSecGroupAssociateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	port, err := ports.Get(client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "error fetching OpenTelekomCloud port security groups"))
	}

	force := false
	if v, ok := d.GetOk("force"); ok {
		force = v.(bool)
	}
	mErr := multierror.Append(nil,
		d.Set("all_security_group_ids", port.SecurityGroups),
		d.Set("force", force),
		d.Set("port_id", d.Id()),
		d.Set("region", config.GetRegion(d)),
	)

	if force {
		mErr = multierror.Append(mErr, d.Set("security_group_ids", port.SecurityGroups))
	} else {
		allSet := d.Get("all_security_group_ids").(*schema.Set)
		desiredSet := d.Get("security_group_ids").(*schema.Set)
		actualSet := allSet.Intersection(desiredSet)
		if !actualSet.Equal(desiredSet) {
			mErr = multierror.Append(mErr, d.Set("security_group_ids", common.ExpandToStringSlice(actualSet.List())))
		}
	}
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving port resource fields: %s", mErr)
	}

	return nil
}

func resourceNetworkingPortSecGroupAssociateV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	var updateOpts ports.UpdateOpts
	var force bool
	if v, ok := d.GetOk("force"); ok {
		force = v.(bool)
	}

	if force {
		securityGroups := ResourcePortSecurityGroupsV2(d)
		updateOpts.SecurityGroups = &securityGroups
	} else {
		allSet := d.Get("all_security_group_ids").(*schema.Set)
		oldIDs, newIDs := d.GetChange("security_group_ids")
		oldSet, newSet := oldIDs.(*schema.Set), newIDs.(*schema.Set)
		allWithoutOld := allSet.Difference(oldSet)
		newSecurityGroups := common.ExpandToStringSlice(allWithoutOld.Union(newSet).List())
		updateOpts.SecurityGroups = &newSecurityGroups
	}

	if d.HasChange("security_group_ids") || d.HasChange("force") {
		log.Printf("[DEBUG] Port Security Group Update Options: %#v", updateOpts.SecurityGroups)
		_, err = ports.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating OpenTelekomCloud Port: %s", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceNetworkingPortSecGroupAssociateV2Read(clientCtx, d, meta)
}

func resourceNetworkingPortSecGroupAssociateV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	var updateOpts ports.UpdateOpts
	var force bool
	if v, ok := d.GetOk("force"); ok {
		force = v.(bool)
	}

	if force {
		updateOpts.SecurityGroups = &[]string{}
	} else {
		allSet := d.Get("all_security_group_ids").(*schema.Set)
		oldSet := d.Get("security_group_ids").(*schema.Set)
		allWithoutOld := allSet.Difference(oldSet)
		newSecurityGroups := common.ExpandToStringSlice(allWithoutOld.List())
		updateOpts.SecurityGroups = &newSecurityGroups
	}
	log.Printf("[DEBUG] Port security groups disassociation options: %#v", updateOpts.SecurityGroups)

	_, err = ports.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "error disassociating OpenTelekomCloud port security groups"))
	}

	return nil
}
