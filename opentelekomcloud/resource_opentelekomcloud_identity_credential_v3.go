package opentelekomcloud

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/credentials"
)

func resourceIdentityCredentialV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityCredentialV3Create,
		Read:   resourceIdentityCredentialV3Read,
		Update: resourceIdentityCredentialV3Update,
		Delete: resourceIdentityCredentialV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"access": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_use_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityCredentialV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.identityV30Client()
	if err != nil {
		return fmt.Errorf("error creating OpenStack identity client: %s", err)
	}

	userID, ok := d.GetOk("user_id")
	if !ok {
		userID = client.UserID
	}

	if userID == "" {
		return fmt.Errorf("error defining current user ID, please either provide " +
			"`user_id` or authenticate with token auth (not using AK/SK)")
	}

	credential, err := credentials.Create(client, credentials.CreateOpts{
		UserID:      userID.(string),
		Description: d.Get("description").(string),
	}).Extract()
	if err != nil {
		return fmt.Errorf("error creating AK/SK: %s", err)
	}

	d.SetId(credential.AccessKey)
	_ = d.Set("secret", credential.SecretKey) // secret key returned only once

	return resourceIdentityCredentialV3Read(d, meta)
}

func resourceIdentityCredentialV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.identityV30Client()
	if err != nil {
		return fmt.Errorf("error creating OpenStack identity client: %s", err)
	}
	credential, err := credentials.Get(client, d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("error retrieving AK/SK information: %s", err)
	}
	return multierror.Append(nil,
		d.Set("user_id", credential.UserID),
		d.Set("access", credential.AccessKey),
		d.Set("status", credential.Status),
		d.Set("create_time", credential.CreateTime),
		d.Set("last_use_time", credential.LastUseTime),
		d.Set("description", credential.Description),
	).ErrorOrNil()
}

func resourceIdentityCredentialV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.identityV30Client()
	if err != nil {
		return fmt.Errorf("error creating OpenStack identity client: %s", err)
	}
	opts := credentials.UpdateOpts{}
	if d.HasChange("status") {
		opts.Status = d.Get("status").(string)
	}
	if d.HasChange("description") {
		opts.Description = d.Get("description").(string)
	}
	_, err = credentials.Update(client, d.Id(), opts).Extract()
	if err != nil {
		return fmt.Errorf("error updating AK/SK: %s", err)
	}
	return resourceIdentityCredentialV3Read(d, meta)
}

func resourceIdentityCredentialV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.identityV30Client()
	if err != nil {
		return fmt.Errorf("error creating OpenStack identity client: %s", err)
	}
	err = credentials.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("error deleting AK/SK: %s", err)
	}
	d.SetId("")
	return nil
}
