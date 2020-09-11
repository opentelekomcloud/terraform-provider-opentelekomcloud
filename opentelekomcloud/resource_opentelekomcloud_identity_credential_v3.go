package opentelekomcloud

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
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
				Required: true,
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

	credential, err := CredentialCreate(client, CredentialCreateOpts{
		UserID:      d.Get("user_id").(string),
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
	credential, err := CredentialGet(client, d.Id()).Extract()
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
	opts := CredentialUpdateOpts{}
	if d.HasChange("status") {
		opts.Status = d.Get("status").(string)
	}
	if d.HasChange("description") {
		opts.Status = d.Get("description").(string)
	}
	_, err = CredentialUpdate(client, d.Id(), opts).Extract()
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
	err = CredentialDelete(client, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("error deleting AK/SK: %s", err)
	}
	d.SetId("")
	return nil
}

// Following to be moved to SDK `openstack/identity/v3/credentials`

// urls.go

const (
	credRootPath    = "OS-CREDENTIAL"
	credentialsPath = "credentials"
)

func credListURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(credRootPath, credentialsPath)
}

func credGetURL(client *golangsdk.ServiceClient, credID string) string {
	return client.ServiceURL(credRootPath, credentialsPath, credID)
}

func credCreateURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(credRootPath, credentialsPath)
}

func credUpdateURL(client *golangsdk.ServiceClient, credID string) string {
	return client.ServiceURL(credRootPath, credentialsPath, credID)
}

func credDeleteURL(client *golangsdk.ServiceClient, credID string) string {
	return client.ServiceURL(credRootPath, credentialsPath, credID)
}

// results.go

type Credential struct {
	// IAM user ID
	UserID string `json:"user_id"`

	// Description of the access key
	Description string `json:"description"`

	// Time when the access key was created
	CreateTime string `json:"create_time"`

	// Time when the access key was last used
	LastUseTime string `json:"last_use_time,omitempty"`

	// AK
	AccessKey string `json:"access"`

	// SK, returned only during creation
	SecretKey string `json:"secret,omitempty"`

	// Status of the access key, active/inactive
	Status CredentialStatus `json:"status"`
}

type credentialResult struct {
	golangsdk.Result
}

// CreateResult is the response of a CredentialGet operations. Call its Extract method to
// interpret it as a Credential.
type CreateResult struct {
	credentialResult
}

// Extract provides access to the individual Flavor returned by the CredentialGet and
// Create functions.
func (r credentialResult) Extract() (*Credential, error) {
	var s struct {
		Credential *Credential `json:"credential"`
	}
	err := r.ExtractInto(&s)
	return s.Credential, err
}

// CredentialGetResult is the response of a CredentialGet operations. Call its Extract method to
// interpret it as a Credential.
type CredentialGetResult struct {
	credentialResult
}

// CredentialUpdateResult is the response from an Update operation. Call its Extract
// method to interpret it as a Role.
type CredentialUpdateResult struct {
	credentialResult
}

type CredentialListResult struct {
	golangsdk.Result
}

func (lr CredentialListResult) Extract() ([]Credential, error) {
	var a struct {
		Instances []Credential `json:"credentials"`
	}
	err := lr.Result.ExtractInto(&a)
	return a.Instances, err
}

// CredentialDeleteResult is the response from a Delete operation. Call its ExtractErr to
// determine if the request succeeded or failed.
type CredentialDeleteResult struct {
	golangsdk.ErrResult
}

// requests.go

type CredentialStatus string

const credentialParentElement = "credential"

type CredentialListOptsBuilder interface {
	ToCredentialListQuery() (string, error)
}

type CredentialListOpts struct {
	UserID string `json:"user_id,omitempty"`
}

func (opts CredentialListOpts) ToCredentialListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func CredentialList(client *golangsdk.ServiceClient, opts CredentialListOptsBuilder) (l CredentialListResult) {
	q, err := opts.ToCredentialListQuery()
	if err != nil {
		l.Err = err
		return
	}
	_, l.Err = client.Get(credListURL(client)+q, &l.Body, nil)
	return
}

func CredentialGet(client *golangsdk.ServiceClient, credentialID string) (r CredentialGetResult) {
	_, r.Err = client.Get(credGetURL(client, credentialID), &r.Body, nil)
	return
}

type CredentialCreateOptsBuilder interface {
	ToCredentialCreateMap() (map[string]interface{}, error)
}

type CredentialCreateOpts struct {
	UserID      string `json:"user_id"`
	Description string `json:"description"`
}

func (opts CredentialCreateOpts) ToCredentialCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, credentialParentElement)
}

func CredentialCreate(client *golangsdk.ServiceClient, opts CredentialCreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToCredentialCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(credCreateURL(client), &b, &r.Body, nil)
	return
}

type CredentialUpdateOptsBuilder interface {
	ToCredentialUpdateMap() (map[string]interface{}, error)
}

type CredentialUpdateOpts struct {
	Status      string `json:"status,omitempty"`
	Description string `json:"description,omitempty"`
}

func (opts CredentialUpdateOpts) ToCredentialUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, credentialParentElement)
}

func CredentialUpdate(client *golangsdk.ServiceClient, credentialID string, opts CredentialUpdateOptsBuilder) (r CreateResult) {
	b, err := opts.ToCredentialUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(credUpdateURL(client, credentialID), &b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200}})
	return
}

func CredentialDelete(client *golangsdk.ServiceClient, credentialID string) (r CredentialDeleteResult) {
	_, r.Err = client.Delete(credDeleteURL(client, credentialID), nil)
	return
}
