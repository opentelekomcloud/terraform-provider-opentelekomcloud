package privatenat

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v3/privatenat/snatrules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourcePrivateNatSnatRuleV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateNatSnatRuleV3Create,
		ReadContext:   resourcePrivateNatSnatRuleV3Read,
		UpdateContext: resourcePrivateNatSnatRuleV3Update,
		DeleteContext: resourcePrivateNatSnatRuleV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"virsubnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transit_ip_ids": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 20,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transit_ip_associations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"transit_ip_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePrivateNatSnatRuleV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NatV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	cidr := d.Get("cidr").(string)
	virSubnetId := d.Get("virsubnet_id").(string)
	if (cidr == "" && virSubnetId == "") || (cidr != "" && virSubnetId != "") {
		return fmterr.Errorf("Error: (Only) one of cidr or virsubnet_id must be specified.")
	}
	createOpts := snatrules.CreatePrivateSnatOpts{
		Description:  d.Get("description").(string),
		Cidr:         cidr,
		GatewayId:    d.Get("gateway_id").(string),
		VirSubnetId:  virSubnetId,
		TransitIpIds: common.ExpandToStringList(d.Get("transit_ip_ids").([]interface{})),
	}

	createResp, err := snatrules.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating Private NAT SNAT rule: %w", err)
	}
	d.SetId(createResp.SnatRule.Id)
	log.Printf("Created Private NAT SNAT rule %s: %#v", d.Id(), createResp.SnatRule)

	return resourcePrivateNatSnatRuleV3Read(ctx, d, meta)
}

func resourcePrivateNatSnatRuleV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NatV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	getResp, err := snatrules.Get(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error fetching private NAT SNAT rule : %w", err)
	}

	natSnatRule := getResp.SnatRule

	transitIpIds, transitIpAssociations := processTransitIpAssociations(natSnatRule.TransitIpAssociations)
	mErr := multierror.Append(
		d.Set("gateway_id", natSnatRule.GatewayId),
		d.Set("cidr", natSnatRule.Cidr),
		d.Set("virsubnet_id", natSnatRule.VirSubnetId),
		d.Set("description", natSnatRule.Description),
		d.Set("transit_ip_ids", transitIpIds),
		d.Set("project_id", natSnatRule.ProjectId),
		d.Set("transit_ip_associations", transitIpAssociations),
		d.Set("created_at", natSnatRule.CreatedAt),
		d.Set("updated_at", natSnatRule.UpdatedAt),
		d.Set("enterprise_project_id", natSnatRule.EnterpriseProjectId),
		d.Set("status", natSnatRule.Status),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePrivateNatSnatRuleV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NatV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	var updateOpts snatrules.UpdatePrivateSnatOpts
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}
	if d.HasChange("transit_ip_ids") {
		updateOpts.TransitIpIds = common.ExpandToStringList(d.Get("transit_ip_ids").([]interface{}))
	}

	log.Printf("[DEBUG] Update Options: %#v", updateOpts)

	_, err = snatrules.Update(client, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating Private NAT SNAT rule: %w", err)
	}

	return resourcePrivateNatSnatRuleV3Read(ctx, d, meta)
}

func resourcePrivateNatSnatRuleV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NatV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	err = snatrules.Delete(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting Private NAT SNAT rule: %w", err)
	}

	d.SetId("")
	return nil
}

func processTransitIpAssociations(transitIpAssociationsInResp []snatrules.TransitIPAssociation) ([]string, []map[string]interface{}) {
	var transitIpAssociations []map[string]interface{}
	var transitIpIds []string
	for _, transitIpAssociationInResp := range transitIpAssociationsInResp {
		transitIpAssociation := map[string]interface{}{
			"transit_ip_id":      transitIpAssociationInResp.TransitIpId,
			"transit_ip_address": transitIpAssociationInResp.TransitIpAddress,
		}
		transitIpIds = append(transitIpIds, transitIpAssociationInResp.TransitIpId)
		transitIpAssociations = append(transitIpAssociations, transitIpAssociation)
	}
	return transitIpIds, transitIpAssociations
}
