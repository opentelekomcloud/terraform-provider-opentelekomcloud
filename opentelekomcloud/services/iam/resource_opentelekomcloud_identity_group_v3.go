package iam

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/groups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceIdentityGroupV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityGroupV3Create,
		ReadContext:   resourceIdentityGroupV3Read,
		UpdateContext: resourceIdentityGroupV3Update,
		DeleteContext: resourceIdentityGroupV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceIdentityGroupV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("Error creating OpenTelekomCloud identity client: %s", err)
	}

	createOpts := groups.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		DomainID:    d.Get("domain_id").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	group, err := groups.Create(identityClient, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("Error creating OpenTelekomCloud Group: %s", err)
	}

	d.SetId(group.ID)

	return resourceIdentityGroupV3Read(ctx, d, meta)
}

func resourceIdentityGroupV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("Error creating OpenTelekomCloud identity client: %s", err)
	}

	group, err := groups.Get(identityClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "group"))
	}

	log.Printf("[DEBUG] Retrieved OpenTelekomCloud Group: %#v", group)

	d.Set("domain_id", group.DomainID)
	d.Set("description", group.Description)
	d.Set("name", group.Name)
	d.Set("region", config.GetRegion(d))

	return nil
}

func resourceIdentityGroupV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("Error creating OpenTelekomCloud identity client: %s", err)
	}

	var hasChange bool
	var updateOpts groups.UpdateOpts

	if d.HasChange("description") {
		hasChange = true
		updateOpts.Description = d.Get("description").(string)
	}

	if d.HasChange("domain_id") {
		hasChange = true
		updateOpts.DomainID = d.Get("domain_id").(string)
	}

	if d.HasChange("name") {
		hasChange = true
		updateOpts.Name = d.Get("name").(string)
	}

	if hasChange {
		log.Printf("[DEBUG] Update Options: %#v", updateOpts)
	}

	if hasChange {
		_, err := groups.Update(identityClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmterr.Errorf("Error updating OpenTelekomCloud group: %s", err)
		}
	}

	return resourceIdentityGroupV3Read(ctx, d, meta)
}

func resourceIdentityGroupV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("Error creating OpenTelekomCloud identity client: %s", err)
	}

	err = groups.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		return fmterr.Errorf("Error deleting OpenTelekomCloud group: %s", err)
	}

	return nil
}
