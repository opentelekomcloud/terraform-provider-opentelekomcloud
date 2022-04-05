package fw

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/fwaas_v2/policies"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceFWPolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFWPolicyV2Create,
		ReadContext:   resourceFWPolicyV2Read,
		UpdateContext: resourceFWPolicyV2Update,
		DeleteContext: resourceFWPolicyV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"audited": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceFWPolicyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	v := d.Get("rules").([]interface{})

	log.Printf("[DEBUG] Rules found : %#v", v)
	log.Printf("[DEBUG] Rules count : %d", len(v))

	rules := make([]string, len(v))
	for i, v := range v {
		rules[i] = v.(string)
	}

	audited := d.Get("audited").(bool)

	opts := PolicyCreateOpts{
		policies.CreateOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Audited:     &audited,
			TenantID:    d.Get("tenant_id").(string),
			Rules:       rules,
		},
		common.MapValueSpecs(d),
	}

	if r, ok := d.GetOk("shared"); ok {
		shared := r.(bool)
		opts.Shared = &shared
	}

	log.Printf("[DEBUG] Create firewall policy: %#v", opts)

	policy, err := policies.Create(networkingClient, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Firewall policy created: %#v", policy)

	d.SetId(policy.ID)

	return resourceFWPolicyV2Read(ctx, d, meta)
}

func resourceFWPolicyV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Retrieve information about firewall policy: %s", d.Id())

	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	policy, err := policies.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "FW policy")
	}

	log.Printf("[DEBUG] Read OpenTelekomCloud Firewall Policy %s: %#v", d.Id(), policy)

	mErr := multierror.Append(
		d.Set("name", policy.Name),
		d.Set("description", policy.Description),
		d.Set("shared", policy.Shared),
		d.Set("audited", policy.Audited),
		d.Set("tenant_id", policy.TenantID),
		d.Set("rules", policy.Rules),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceFWPolicyV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	opts := policies.UpdateOpts{}

	if d.HasChange("name") {
		opts.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		opts.Description = d.Get("description").(string)
	}

	if d.HasChange("rules") {
		v := d.Get("rules").([]interface{})

		log.Printf("[DEBUG] Rules found : %#v", v)
		log.Printf("[DEBUG] Rules count : %d", len(v))

		rules := make([]string, len(v))
		for i, v := range v {
			rules[i] = v.(string)
		}
		opts.Rules = rules
	}

	log.Printf("[DEBUG] Updating firewall policy with id %s: %#v", d.Id(), opts)

	err = policies.Update(networkingClient, d.Id(), opts).Err
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceFWPolicyV2Read(ctx, d, meta)
}

func resourceFWPolicyV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Destroy firewall policy: %s", d.Id())

	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForFirewallPolicyDeletion(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 2 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmterr.Errorf("error waiting for FW policy to be deleted: %w", err)
	}

	return nil
}

func waitForFirewallPolicyDeletion(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		err := policies.Delete(networkingClient, id).ExtractErr()
		if err == nil {
			return "", "DELETED", nil
		}
		if _, ok := err.(golangsdk.ErrDefault409); ok {
			return nil, "ACTIVE", nil
		}
		return nil, "ACTIVE", err
	}
}
