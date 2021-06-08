package iam

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/federation/mappings"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

const mappingError = "error %s identity mapping v3: %w"

func ResourceIdentityMappingV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityMappingV3Create,
		ReadContext:   resourceIdentityMappingV3Read,
		UpdateContext: resourceIdentityMappingV3Update,
		DeleteContext: resourceIdentityMappingV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"mapping_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rules": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: common.ValidateJsonString,
				StateFunc: func(v interface{}) string {
					jsonString, _ := common.NormalizeJsonString(v)
					return jsonString
				},
			},
			"links": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceIdentityMappingV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	rulesRaw := d.Get("rules").(string)
	rulesBytes := []byte(rulesRaw)
	rules := make([]mappings.RuleOpts, 1)
	if err := json.Unmarshal(rulesBytes, &rules); err != nil {
		return diag.FromErr(err)
	}

	createOpts := mappings.CreateOpts{
		Rules: rules,
	}
	mappingID := d.Get("mapping_id").(string)
	mapping, err := mappings.Create(client, mappingID, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf(mappingError, "creating", err)
	}

	d.SetId(mapping.ID)

	return resourceIdentityMappingV3Read(ctx, d, meta)
}

func resourceIdentityMappingV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	mapping, err := mappings.Get(client, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}
		return fmterr.Errorf(mappingError, "reading", err)
	}

	rules, err := json.Marshal(mapping.Rules)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rules", string(rules)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("mapping_id", mapping.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("links", mapping.Links); err != nil {
		return fmterr.Errorf("error setting identity mapping links: %w", err)
	}

	return nil
}

func resourceIdentityMappingV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}
	changes := false
	updateOpts := mappings.UpdateOpts{}

	if d.HasChange("rules") {
		changes = true
		rulesRaw := d.Get("rules").(string)
		rulesBytes := []byte(rulesRaw)
		rules := make([]mappings.RuleOpts, 1)
		if err := json.Unmarshal(rulesBytes, &rules); err != nil {
			return diag.FromErr(err)
		}
		updateOpts.Rules = rules
	}
	if changes {
		_, err := mappings.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmterr.Errorf(mappingError, "updating", err)
		}
	}

	return resourceIdentityMappingV3Read(ctx, d, meta)
}

func resourceIdentityMappingV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	if err := mappings.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmterr.Errorf(mappingError, "deleting", err)
	}

	return nil
}
