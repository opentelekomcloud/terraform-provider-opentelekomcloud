package vpn

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/vpnaas/ikepolicies"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func ResourceVpnIKEPolicyV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpnIKEPolicyV2Create,
		Read:   resourceVpnIKEPolicyV2Read,
		Update: resourceVpnIKEPolicyV2Update,
		Delete: resourceVpnIKEPolicyV2Delete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auth_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "sha1",
				ValidateFunc: validation.StringInSlice([]string{
					"md5", "sha1", "sha2-256", "sha2-384", "sha2-512",
				}, false),
			},
			"encryption_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "aes-128",
				ValidateFunc: validation.StringInSlice([]string{
					"3des", "aes-128", "aes-192", "aes-256",
				}, false),
			},
			"pfs": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "group5",
				ValidateFunc: validation.StringInSlice([]string{
					"group1", "group2", "group5", "group14", "group15", "group16", "group19", "group20", "group21",
				}, false),
			},
			"phase1_negotiation_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "main",
				ValidateFunc: validation.StringInSlice([]string{
					"main", "aggressive",
				}, false),
			},
			"ike_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "v1",
				ValidateFunc: validation.StringInSlice([]string{
					"v1", "v2",
				}, false),
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
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVpnIKEPolicyV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
	}

	lifetimeRaw := d.Get("lifetime").(*schema.Set).List()
	var lifetime *ikepolicies.LifetimeCreateOpts
	if len(lifetimeRaw) == 1 {
		lifetimeInfo := lifetimeRaw[0].(map[string]interface{})
		lifetime = &ikepolicies.LifetimeCreateOpts{
			Units: ikepolicies.Unit(lifetimeInfo["units"].(string)),
			Value: lifetimeInfo["value"].(int),
		}
	}

	opts := VpnIKEPolicyCreateOpts{
		ikepolicies.CreateOpts{
			Name:                  d.Get("name").(string),
			Description:           d.Get("description").(string),
			TenantID:              d.Get("tenant_id").(string),
			Lifetime:              lifetime,
			AuthAlgorithm:         ikepolicies.AuthAlgorithm(d.Get("auth_algorithm").(string)),
			EncryptionAlgorithm:   ikepolicies.EncryptionAlgorithm(d.Get("encryption_algorithm").(string)),
			PFS:                   ikepolicies.PFS(d.Get("pfs").(string)),
			IKEVersion:            ikepolicies.IKEVersion(d.Get("ike_version").(string)),
			Phase1NegotiationMode: ikepolicies.Phase1NegotiationMode(d.Get("phase1_negotiation_mode").(string)),
		},
		common.MapValueSpecs(d),
	}
	log.Printf("[DEBUG] Create IKE policy: %#v", opts)

	policy, err := ikepolicies.Create(networkingClient, opts).Extract()
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_CREATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForIKEPolicyCreate(networkingClient, policy.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForState()

	log.Printf("[DEBUG] IKE policy created: %#v", policy)

	d.SetId(policy.ID)

	return resourceVpnIKEPolicyV2Read(d, meta)
}

func resourceVpnIKEPolicyV2Read(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Retrieve information about IKE policy: %s", d.Id())

	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
	}

	policy, err := ikepolicies.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "IKE policy")
	}

	log.Printf("[DEBUG] Read OpenTelekomCloud IKE Policy %s: %#v", d.Id(), policy)

	mErr := multierror.Append(nil,
		d.Set("name", policy.Name),
		d.Set("description", policy.Description),
		d.Set("auth_algorithm", policy.AuthAlgorithm),
		d.Set("encryption_algorithm", policy.EncryptionAlgorithm),
		d.Set("tenant_id", policy.TenantID),
		d.Set("pfs", policy.PFS),
		d.Set("phase1_negotiation_mode", policy.Phase1NegotiationMode),
		d.Set("ike_version", policy.IKEVersion),
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

	if err := mErr.ErrorOrNil(); err != nil {
		return fmt.Errorf("error setting IKE policy fields: %s", err)
	}

	return nil
}

func resourceVpnIKEPolicyV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
	}

	opts := ikepolicies.UpdateOpts{}
	var hasChange bool

	if d.HasChange("name") {
		name := d.Get("name").(string)
		opts.Name = name
		hasChange = true
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		opts.Description = description
		hasChange = true
	}
	if d.HasChange("pfs") {
		opts.PFS = ikepolicies.PFS(d.Get("pfs").(string))
		hasChange = true
	}
	if d.HasChange("auth_algorithm") {
		opts.AuthAlgorithm = ikepolicies.AuthAlgorithm(d.Get("auth_algorithm").(string))
		hasChange = true
	}
	if d.HasChange("encryption_algorithm") {
		opts.EncryptionAlgorithm = ikepolicies.EncryptionAlgorithm(d.Get("encryption_algorithm").(string))
		hasChange = true
	}
	if d.HasChange("phase_1_negotiation_mode") {
		opts.Phase1NegotiationMode = ikepolicies.Phase1NegotiationMode(d.Get("phase_1_negotiation_mode").(string))
		hasChange = true
	}
	if d.HasChange("ike_version") {
		opts.IKEVersion = ikepolicies.IKEVersion(d.Get("ike_version").(string))
		hasChange = true
	}

	if d.HasChange("lifetime") {
		lifetimeRaw := d.Get("lifetime").(*schema.Set).List()
		var lifetime *ikepolicies.LifetimeUpdateOpts
		if len(lifetimeRaw) == 1 {
			lifetimeInfo := lifetimeRaw[0].(map[string]interface{})
			lifetime = &ikepolicies.LifetimeUpdateOpts{
				Units: ikepolicies.Unit(lifetimeInfo["units"].(string)),
				Value: lifetimeInfo["value"].(int),
			}
		}
		opts.Lifetime = lifetime
		hasChange = true
	}

	log.Printf("[DEBUG] Updating IKE policy with id %s: %#v", d.Id(), opts)

	if hasChange {
		err = ikepolicies.Update(networkingClient, d.Id(), opts).Err
		if err != nil {
			return err
		}
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"PENDING_UPDATE"},
			Target:     []string{"ACTIVE"},
			Refresh:    waitForIKEPolicyUpdate(networkingClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			MinTimeout: 2 * time.Second,
		}
		if _, err = stateConf.WaitForState(); err != nil {
			return err
		}
	}

	return resourceVpnIKEPolicyV2Read(d, meta)
}

func resourceVpnIKEPolicyV2Delete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Destroy IKE policy: %s", d.Id())

	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForIKEPolicyDelete(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 2 * time.Second,
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return err
	}

	return nil
}

func waitForIKEPolicyDelete(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		err := ikepolicies.Delete(networkingClient, id).Err
		if err == nil {
			return "", "DELETED", nil
		}

		return nil, "ACTIVE", err
	}
}

func waitForIKEPolicyCreate(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := ikepolicies.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_CREATE", nil
		}
		return policy, "ACTIVE", nil
	}
}

func waitForIKEPolicyUpdate(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := ikepolicies.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_UPDATE", nil
		}
		return policy, "ACTIVE", nil
	}
}
