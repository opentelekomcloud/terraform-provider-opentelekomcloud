package iam

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/users"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func ResourceIdentityUserV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityUserV3Create,
		Read:   resourceIdentityUserV3Read,
		Update: resourceIdentityUserV3Update,
		Delete: resourceIdentityUserV3Delete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"default_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"email": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: common.SuppressCaseInsensitive,
			},

			"send_welcome_email": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceIdentityUserV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf(clientCreationFail, err)
	}

	enabled := d.Get("enabled").(bool)
	createOpts := users.CreateOpts{
		DefaultProjectID: d.Get("default_project_id").(string),
		DomainID:         d.Get("domain_id").(string),
		Enabled:          &enabled,
		Name:             d.Get("name").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	user, err := users.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud user: %s", err)
	}

	d.SetId(user.ID)

	return setExtendedOpts(d, meta)
}

func resourceIdentityUserV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf(clientCreationFail, err)
	}

	user, err := users.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "user")
	}

	log.Printf("[DEBUG] Retrieved OpenStack user: %#v", user)

	mErr := multierror.Append(nil,
		d.Set("default_project_id", user.DefaultProjectID),
		d.Set("domain_id", user.DomainID),
		d.Set("enabled", user.Enabled),
		d.Set("name", user.Name),
		d.Set("region", config.GetRegion(d)),
	)

	// Read extended options
	user, err = users.ExtendedUpdate(client, d.Id(), users.ExtendedUpdateOpts{}).Extract()
	if err != nil {
		return err
	}
	mErr = multierror.Append(mErr,
		d.Set("email", user.Email),
	)

	if mErr.ErrorOrNil() != nil {
		return mErr
	}

	return nil
}

func setExtendedOpts(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf(clientCreationFail, err)
	}

	var hasChange bool
	var updateOpts users.ExtendedUpdateOpts

	if d.HasChange("email") {
		hasChange = true
		updateOpts.Email = d.Get("email").(string)
	}

	if hasChange {
		_, err := users.ExtendedUpdate(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("error updating OpenTelekomCloud user: %w", err)
		}

		sendWelcomeEmail := d.Get("send_welcome_email").(bool)
		if sendWelcomeEmail {
			if err := users.SendWelcomeEmail(client, d.Id()).ExtractErr(); err != nil {
				return fmt.Errorf("error sending a welcome email: %w", err)
			}
		}
	}

	return resourceIdentityUserV3Read(d, meta)
}

func resourceIdentityUserV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf(clientCreationFail, err)
	}

	var hasChange bool
	var updateOpts users.UpdateOpts

	if d.HasChange("default_project_id") {
		hasChange = true
		updateOpts.DefaultProjectID = d.Get("default_project_id").(string)
	}

	if d.HasChange("domain_id") {
		hasChange = true
		updateOpts.DomainID = d.Get("domain_id").(string)
	}

	if d.HasChange("enabled") {
		hasChange = true
		enabled := d.Get("enabled").(bool)
		updateOpts.Enabled = &enabled
	}

	if d.HasChange("name") {
		hasChange = true
		updateOpts.Name = d.Get("name").(string)
	}

	if hasChange {
		log.Printf("[DEBUG] Update Options: %#v", updateOpts)
	}

	if d.HasChange("password") {
		hasChange = true
		updateOpts.Password = d.Get("password").(string)
	}

	if hasChange {
		_, err := users.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("error updating OpenTelekomCloud user: %w", err)
		}
	}

	return setExtendedOpts(d, meta)
}

func resourceIdentityUserV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf(clientCreationFail, err)
	}

	err = users.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("error deleting OpenTelekomCloud user: %w", err)
	}

	return nil
}
