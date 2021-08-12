package kms

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/kms/v1/keys"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

const (
	WaitingForEnableState = "1"
	EnabledState          = "2"
	DisabledState         = "3"
	PendingDeletionState  = "4"
	WaitingImportState    = "5"
)

func ResourceKmsKeyV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKmsKeyV1Create,
		ReadContext:   resourceKmsKeyV1Read,
		UpdateContext: resourceKmsKeyV1Update,
		DeleteContext: resourceKmsKeyV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"key_alias": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"realm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scheduled_deletion_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"default_key_flag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiration_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"origin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pending_days": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "7",
			},
			"tags": common.TagsSchema(),
		},
	}
}

func resourceKmsKeyV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.KmsKeyV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud KMSv1 client: %s", err)
	}

	createOpts := &keys.CreateOpts{
		KeyAlias:       d.Get("key_alias").(string),
		KeyDescription: d.Get("key_description").(string),
		Realm:          d.Get("realm").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	key, err := keys.Create(client, createOpts).ExtractKeyInfo()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud key: %s", err)
	}
	log.Printf("[INFO] Key ID: %s", key.KeyID)

	// Wait for the key to become enabled.
	log.Printf("[DEBUG] Waiting for key (%s) to become enabled", key.KeyID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{WaitingForEnableState, DisabledState},
		Target:     []string{EnabledState},
		Refresh:    keyV1StateRefreshFunc(client, key.KeyID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for key (%s) to become ready: %s", key.KeyID, err)
	}

	if !d.Get("is_enabled").(bool) {
		disableKey, err := keys.DisableKey(client, key.KeyID).ExtractKeyInfo()
		if err != nil {
			return fmterr.Errorf("error disabling key: %s", err)
		}

		if disableKey.KeyState != DisabledState {
			return fmterr.Errorf("error disabling key, the key state is: %s", disableKey.KeyState)
		}
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		tagList := common.ExpandResourceTags(tagRaw)
		if err := tags.Create(client, "kms", key.KeyID, tagList).ExtractErr(); err != nil {
			return fmterr.Errorf("error setting tags of KMS: %s", err)
		}
	}

	// Store the key ID now
	d.SetId(key.KeyID)

	return resourceKmsKeyV1Read(ctx, d, meta)
}

func resourceKmsKeyV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.KmsKeyV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud KMSv1 client: %s", err)
	}

	key, err := keys.Get(client, d.Id()).ExtractKeyInfo()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Kms key %s: %+v", d.Id(), key)
	if key.KeyState == PendingDeletionState {
		log.Printf("[WARN] Removing KMS key %s because it's already gone", d.Id())
		d.SetId("")
		return nil
	}

	mErr := multierror.Append(nil,
		d.Set("domain_id", key.DomainID),
		d.Set("key_alias", key.KeyAlias),
		d.Set("realm", key.Realm),
		d.Set("key_description", key.KeyDescription),
		d.Set("creation_date", key.CreationDate),
		d.Set("scheduled_deletion_date", key.ScheduledDeletionDate),
		d.Set("is_enabled", key.KeyState == EnabledState),
		d.Set("default_key_flag", key.DefaultKeyFlag),
		d.Set("expiration_time", key.ExpirationTime),
		d.Set("origin", key.Origin),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	// save tags
	resourceTags, err := tags.Get(client, "kms", d.Id()).Extract()
	if err != nil {
		return fmterr.Errorf("error fetching OpenTelekomCloud KMS tags: %s", err)
	}
	tagMap := common.TagsToMap(resourceTags)
	if err := d.Set("tags", tagMap); err != nil {
		return fmterr.Errorf("error saving tags for OpenTelekomCloud KMS: %s", err)
	}

	return nil
}

func resourceKmsKeyV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.KmsKeyV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud KMSv1 client: %s", err)
	}

	if d.HasChange("key_alias") {
		updateAliasOpts := keys.UpdateAliasOpts{
			KeyID:    d.Id(),
			KeyAlias: d.Get("key_alias").(string),
		}
		_, err = keys.UpdateAlias(client, updateAliasOpts).ExtractKeyInfo()
		if err != nil {
			return fmterr.Errorf("error updating OpenTelekomCloud key: %s", err)
		}
	}

	if d.HasChange("key_description") {
		updateDesOpts := keys.UpdateDesOpts{
			KeyID:          d.Id(),
			KeyDescription: d.Get("key_description").(string),
		}
		_, err = keys.UpdateDes(client, updateDesOpts).ExtractKeyInfo()
		if err != nil {
			return fmterr.Errorf("error updating OpenTelekomCloud key: %s", err)
		}
	}

	if d.HasChange("is_enabled") {
		key, err := keys.Get(client, d.Id()).ExtractKeyInfo()
		if err != nil {
			return fmterr.Errorf("describeKey got an error: %s", err)
		}

		if d.Get("is_enabled").(bool) && key.KeyState == DisabledState {
			key, err := keys.EnableKey(client, d.Id()).ExtractKeyInfo()
			if err != nil {
				return fmterr.Errorf("error enabling key: %s", err)
			}
			if key.KeyState != EnabledState {
				return fmterr.Errorf("error enabling key, the key state is: %s", key.KeyState)
			}
		}

		if !d.Get("is_enabled").(bool) && key.KeyState == EnabledState {
			key, err := keys.DisableKey(client, d.Id()).ExtractKeyInfo()
			if err != nil {
				return fmterr.Errorf("error disabling key: %s", err)
			}
			if key.KeyState != DisabledState {
				return fmterr.Errorf("error disabling key, the key state is: %s", key.KeyState)
			}
		}
	}

	// update tags
	if d.HasChange("tags") {
		if err := common.UpdateResourceTags(client, d, "kms", d.Id()); err != nil {
			return fmterr.Errorf("error updating tags of KMS %s: %s", d.Id(), err)
		}
	}

	return resourceKmsKeyV1Read(ctx, d, meta)
}

func resourceKmsKeyV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.KmsKeyV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud KMSv1 client: %s", err)
	}

	key, err := keys.Get(client, d.Id()).ExtractKeyInfo()
	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "key"))
	}

	deleteOpts := &keys.DeleteOpts{
		KeyID: d.Id(),
	}
	if v, ok := d.GetOk("pending_days"); ok {
		deleteOpts.PendingDays = v.(string)
	}

	// It's possible that this key was used as a boot device and is currently
	// in a pending deletion state from when the instance was terminated.
	// If this is true, just move on. It'll eventually delete.
	if key.KeyState != PendingDeletionState {
		key, err = keys.Delete(client, deleteOpts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		if key.KeyState != PendingDeletionState {
			return fmterr.Errorf("failed to delete key")
		}
	}

	log.Printf("[DEBUG] KMS Key %s deactivated.", d.Id())
	d.SetId("")
	return nil
}

func keyV1StateRefreshFunc(client *golangsdk.ServiceClient, keyID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := keys.Get(client, keyID).ExtractKeyInfo()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, PendingDeletionState, nil
			}
			return nil, "", err
		}

		return v, v.KeyState, nil
	}
}
