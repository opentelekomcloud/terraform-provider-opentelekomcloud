package opentelekomcloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3/users"
)

func dataSourceIdentityUserV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIdentityUserV3Read,

		Schema: map[string]*schema.Schema{
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

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

// dataSourceIdentityUserV3Read performs the user lookup.
func dataSourceIdentityUserV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.identityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	enabled := d.Get("enabled").(bool)
	listOpts := users.ListOpts{
		DomainID: d.Get("domain_id").(string),
		Enabled:  &enabled,
		Name:     d.Get("name").(string),
	}

	log.Printf("[DEBUG] List Options: %#v", listOpts)

	var user users.User
	allPages, err := users.List(identityClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to query users: %s", err)
	}

	allUsers, err := users.ExtractUsers(allPages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve users: %s", err)
	}

	if len(allUsers) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allUsers) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allUsers)
		return fmt.Errorf("Your query returned more than one result")
	}
	user = allUsers[0]

	log.Printf("[DEBUG] Single user found: %s", user.ID)
	return dataSourceIdentityUserV3Attributes(d, &user)
}

// dataSourceIdentityUserV3Attributes populates the fields of an User resource.
func dataSourceIdentityUserV3Attributes(d *schema.ResourceData, user *users.User) error {
	log.Printf("[DEBUG] opentelekomcloud_identity_user_v3 details: %#v", user)

	d.SetId(user.ID)
	d.Set("default_project_id", user.DefaultProjectID)
	d.Set("domain_id", user.DomainID)
	d.Set("enabled", user.Enabled)
	d.Set("name", user.Name)
	d.Set("password_expires_at", user.PasswordExpiresAt.Format(time.RFC3339))

	return nil
}

// Ensure that password_expires_at query matches format explained in
// https://developer.openstack.org/api-ref/identity/v3/#list-users
func validatePasswordExpiresAtQuery(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	values := strings.SplitN(value, ":", 2)
	if len(values) != 2 {
		err := fmt.Errorf("%s '%s' does not match expected format: {operator}:{timestamp}", k, value)
		errors = append(errors, err)
	}
	operator, timestamp := values[0], values[1]

	validOperators := map[string]bool{
		"lt":  true,
		"lte": true,
		"gt":  true,
		"gte": true,
		"eq":  true,
		"neq": true,
	}
	if !validOperators[operator] {
		err := fmt.Errorf("'%s' is not a valid operator for %s. Choose one of 'lt', 'lte', 'gt', 'gte', 'eq', 'neq'", operator, k)
		errors = append(errors, err)
	}

	_, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		err = fmt.Errorf("'%s' is not a valid timestamp for %s. It should be in the form 'YYYY-MM-DDTHH:mm:ssZ'", timestamp, k)
		errors = append(errors, err)
	}

	return
}
