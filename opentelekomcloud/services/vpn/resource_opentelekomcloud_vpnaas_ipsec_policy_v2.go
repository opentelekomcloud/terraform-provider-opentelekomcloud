package vpn

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceVpnIPSecPolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpnIPSecPolicyV2Create,
		ReadContext:   resourceVpnIPSecPolicyV2Read,
		UpdateContext: resourceVpnIPSecPolicyV2Update,
		DeleteContext: resourceVpnIPSecPolicyV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
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
				Optional:     true,
				ValidateFunc: common.ValidateName,
			},
			"auth_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"md5", "sha1", "sha2-256", "sha2-384", "sha2-512",
				}, false),
			},
			"encapsulation_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pfs": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"group1", "group2", "group5", "group14", "group15", "group16", "group19", "group20", "group21", "disable",
				}, false),
			},
			"encryption_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"3des", "aes-128", "aes-192", "aes-256",
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transform_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"esp", "ah", "ah-esp",
				}, false),
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"lifetime": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"units": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"value": {
							Type:         schema.TypeInt,
							Computed:     true,
							Optional:     true,
							ValidateFunc: validation.IntBetween(60, 604800),
						},
					},
				},
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVpnIPSecPolicyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
	}

	lifetimeRaw := d.Get("lifetime").(*schema.Set).List()
	var lifetime *ipsecpolicies.LifetimeCreateOpts
	if len(lifetimeRaw) == 1 {
		lifetimeInfo := lifetimeRaw[0].(map[string]interface{})
		lifetime = &ipsecpolicies.LifetimeCreateOpts{
			Units: ipsecpolicies.Unit(lifetimeInfo["units"].(string)),
			Value: lifetimeInfo["value"].(int),
		}
	}

	opts := VpnIPSecPolicyCreateOpts{
		ipsecpolicies.CreateOpts{
			TenantID:            d.Get("tenant_id").(string),
			Description:         d.Get("description").(string),
			Name:                d.Get("name").(string),
			AuthAlgorithm:       ipsecpolicies.AuthAlgorithm(d.Get("auth_algorithm").(string)),
			EncapsulationMode:   ipsecpolicies.EncapsulationMode(d.Get("encapsulation_mode").(string)),
			EncryptionAlgorithm: ipsecpolicies.EncryptionAlgorithm(d.Get("encryption_algorithm").(string)),
			PFS:                 ipsecpolicies.PFS(d.Get("pfs").(string)),
			TransformProtocol:   ipsecpolicies.TransformProtocol(d.Get("transform_protocol").(string)),
			Lifetime:            lifetime,
		},
		common.MapValueSpecs(d),
	}

	log.Printf("[DEBUG] Create IPSec policy: %#v", opts)

	policy, err := ipsecpolicies.Create(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_CREATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForIPSecPolicyCreate(client, policy.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for policy to become active: %w", err)
	}

	log.Printf("[DEBUG] IPSec policy created: %#v", policy)

	d.SetId(policy.ID)

	return resourceVpnIPSecPolicyV2Read(ctx, d, meta)
}

func resourceVpnIPSecPolicyV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
	}

	policy, err := ipsecpolicies.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "IPSec policy")
	}

	log.Printf("[DEBUG] Read OpenTelekomCloud IPSec policy %s: %#v", d.Id(), policy)

	mErr := multierror.Append(nil,
		d.Set("name", policy.Name),
		d.Set("description", policy.Description),
		d.Set("tenant_id", policy.TenantID),
		d.Set("encapsulation_mode", policy.EncapsulationMode),
		d.Set("encryption_algorithm", policy.EncryptionAlgorithm),
		d.Set("transform_protocol", policy.TransformProtocol),
		d.Set("pfs", policy.PFS),
		d.Set("auth_algorithm", policy.AuthAlgorithm),
		d.Set("region", config.GetRegion(d)),
	)

	// Set the lifetime
	lifetime := []map[string]interface{}{
		{
			"units": policy.Lifetime.Units,
			"value": policy.Lifetime.Value,
		},
	}
	mErr = multierror.Append(mErr, d.Set("lifetime", lifetime))

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func resourceVpnIPSecPolicyV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
	}

	var hasChange bool
	opts := ipsecpolicies.UpdateOpts{}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		opts.Name = &name
		hasChange = true
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		opts.Description = &description
		hasChange = true
	}

	if d.HasChange("auth_algorithm") {
		opts.AuthAlgorithm = ipsecpolicies.AuthAlgorithm(d.Get("auth_algorithm").(string))
		hasChange = true
	}

	if d.HasChange("encryption_algorithm") {
		opts.EncryptionAlgorithm = ipsecpolicies.EncryptionAlgorithm(d.Get("encryption_algorithm").(string))
		hasChange = true
	}

	if d.HasChange("transform_protocol") {
		opts.TransformProtocol = ipsecpolicies.TransformProtocol(d.Get("transform_protocol").(string))
		hasChange = true
	}

	if d.HasChange("pfs") {
		opts.PFS = ipsecpolicies.PFS(d.Get("pfs").(string))
		hasChange = true
	}

	if d.HasChange("encapsulation_mode") {
		opts.EncapsulationMode = ipsecpolicies.EncapsulationMode(d.Get("encapsulation_mode").(string))
		hasChange = true
	}

	if d.HasChange("lifetime") {
		lifetimeRaw := d.Get("lifetime").(*schema.Set).List()
		var lifetime *ipsecpolicies.LifetimeUpdateOpts
		if len(lifetimeRaw) == 1 {
			lifetimeInfo := lifetimeRaw[0].(map[string]interface{})
			lifetime = &ipsecpolicies.LifetimeUpdateOpts{
				Units: ipsecpolicies.Unit(lifetimeInfo["units"].(string)),
				Value: lifetimeInfo["value"].(int),
			}
		}
		opts.Lifetime = lifetime
		hasChange = true
	}

	log.Printf("[DEBUG] Updating IPSec policy with id %s: %#v", d.Id(), opts)

	if hasChange {
		_, err = ipsecpolicies.Update(client, d.Id(), opts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"PENDING_UPDATE"},
			Target:     []string{"ACTIVE"},
			Refresh:    waitForIPSecPolicyUpdate(client, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			MinTimeout: 2 * time.Second,
		}
		if _, err = stateConf.WaitForStateContext(ctx); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceVpnIPSecPolicyV2Read(ctx, d, meta)
}

func resourceVpnIPSecPolicyV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
	}

	if err := ipsecpolicies.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmterr.Errorf("error deleting IPSec Poilicy: %s", err)
	}

	return nil
}

func waitForIPSecPolicyCreate(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := ipsecpolicies.Get(client, id).Extract()
		if err != nil {
			return "", "PENDING_CREATE", nil
		}
		return policy, "ACTIVE", nil
	}
}

func waitForIPSecPolicyUpdate(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := ipsecpolicies.Get(client, id).Extract()
		if err != nil {
			return "", "PENDING_UPDATE", nil
		}
		return policy, "ACTIVE", nil
	}
}
