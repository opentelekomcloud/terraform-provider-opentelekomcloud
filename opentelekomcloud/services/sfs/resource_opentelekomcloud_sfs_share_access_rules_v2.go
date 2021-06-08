package sfs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/sfs/v2/shares"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceSFSShareAccessRulesV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSFSShareAccessRulesV2Create,
		ReadContext:   resourceSFSShareAccessRulesV2Read,
		UpdateContext: resourceSFSShareAccessRulesV2Update,
		DeleteContext: resourceSFSShareAccessRulesV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"share_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_rule": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 20,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_level": {
							Type:     schema.TypeString,
							Required: true,
						},
						"access_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "cert",
						},
						"access_to": {
							Type:     schema.TypeString,
							Required: true,
						},
						"access_rule_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"share_access_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceSFSShareAccessRulesV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SfsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud File Share client: %w", err)
	}

	shareID := d.Get("share_id").(string)
	accessRules := d.Get("access_rule").([]interface{})
	for _, rule := range accessRules {
		accessRuleMap := rule.(map[string]interface{})
		grantAccessOpts := shares.GrantAccessOpts{
			AccessLevel: accessRuleMap["access_level"].(string),
			AccessType:  accessRuleMap["access_type"].(string),
			AccessTo:    accessRuleMap["access_to"].(string),
		}
		_, err = shares.GrantAccess(client, shareID, grantAccessOpts).ExtractAccess()
		if err != nil {
			return fmterr.Errorf("error applying access rule for OpenTelekomCloud File Share: %w", err)
		}
	}

	d.SetId(shareID)

	return resourceSFSShareAccessRulesV2Read(ctx, d, meta)
}

func resourceSFSShareAccessRulesV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SfsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud File Share client: %w", err)
	}

	rules, err := shares.ListAccessRights(client, d.Id()).ExtractAccessRights()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving rules of OpenTelekomCloud File Share: %s", err)
	}

	var accessRules []map[string]interface{}
	for _, v := range rules {
		accessRule := make(map[string]interface{})
		accessRule["access_level"] = v.AccessLevel
		accessRule["access_to"] = v.AccessTo
		accessRule["access_type"] = v.AccessType
		accessRule["access_rule_status"] = v.State
		accessRule["share_access_id"] = v.ID

		accessRules = append(accessRules, accessRule)
	}

	if err := d.Set("access_rule", accessRules); err != nil {
		return fmterr.Errorf("error saving access_rule to state for OpenTelekomCloud File Share: %w", err)
	}

	if err := d.Set("share_id", d.Id()); err != nil {
		return fmterr.Errorf("error saving share_id to state for OpenTelekomCloud File Share: %w", err)
	}

	return nil
}

func resourceSFSShareAccessRulesV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SfsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud File Share client: %w", err)
	}

	accessRules := d.Get("access_rule").([]interface{})
	for _, rule := range accessRules {
		accessRuleMap := rule.(map[string]interface{})
		deleteAccessOpts := shares.DeleteAccessOpts{
			AccessID: accessRuleMap["share_access_id"].(string),
		}
		if err := shares.DeleteAccess(client, d.Id(), deleteAccessOpts).Err; err != nil {
			return fmterr.Errorf("error deleting access rule for OpenTelekomCloud File Share: %w", err)
		}
	}

	d.SetId("")
	return nil
}

func resourceSFSShareAccessRulesV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SfsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud File Share client: %w", err)
	}

	if d.HasChange("access_rule") {
		oldMapRaw, newMapRaw := d.GetChange("access_rule")
		oldMap := oldMapRaw.([]interface{})
		newMap := newMapRaw.([]interface{})

		for _, oldRule := range oldMap {
			oldAccessRuleMap := oldRule.(map[string]interface{})
			deleteAccessOpts := shares.DeleteAccessOpts{
				AccessID: oldAccessRuleMap["share_access_id"].(string),
			}
			if err := shares.DeleteAccess(client, d.Id(), deleteAccessOpts).Err; err != nil {
				return fmterr.Errorf("error deleting access rule for OpenTelekomCloud File Share: %w", err)
			}
		}

		for _, newRule := range newMap {
			newAccessRuleMap := newRule.(map[string]interface{})
			grantAccessOpts := shares.GrantAccessOpts{
				AccessLevel: newAccessRuleMap["access_level"].(string),
				AccessType:  newAccessRuleMap["access_type"].(string),
				AccessTo:    newAccessRuleMap["access_to"].(string),
			}
			_, err = shares.GrantAccess(client, d.Id(), grantAccessOpts).ExtractAccess()
			if err != nil {
				return fmterr.Errorf("error applying access rule for OpenTelekomCloud File Share: %w", err)
			}
		}
	}

	return resourceSFSShareAccessRulesV2Read(ctx, d, meta)
}
