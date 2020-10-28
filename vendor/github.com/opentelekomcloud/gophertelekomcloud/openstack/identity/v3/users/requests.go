package users

import (
	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/groups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/projects"
	"github.com/opentelekomcloud/gophertelekomcloud/pagination"
)

// Option is a specific option defined at the API to enable features
// on a user account.
type Option string

const (
	IgnoreChangePasswordUponFirstUse Option = "ignore_change_password_upon_first_use"
	IgnorePasswordExpiry             Option = "ignore_password_expiry"
	IgnoreLockoutFailureAttempts     Option = "ignore_lockout_failure_attempts"
	MultiFactorAuthRules             Option = "multi_factor_auth_rules"
	MultiFactorAuthEnabled           Option = "multi_factor_auth_enabled"
)

// ListOptsBuilder allows extensions to add additional parameters to
// the List request
type ListOptsBuilder interface {
	ToUserListQuery() (string, error)
}

// ListOpts provides options to filter the List results.
type ListOpts struct {
	// DomainID filters the response by a domain ID.
	DomainID string `q:"domain_id"`

	// Enabled filters the response by enabled users.
	Enabled *bool `q:"enabled"`

	// Name filters the response by username.
	Name string `q:"name"`
}

// ToUserListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToUserListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// List enumerates the Users to which the current token has access.
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToUserListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return UserPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves details on a single user, by ID.
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateOptsBuilder interface {
	ToUserCreateMap() (map[string]interface{}, error)
}

// CreateOpts provides options used to create a user.
type CreateOpts struct {
	// Name is the name of the new user.
	Name string `json:"name" required:"true"`

	// DefaultProjectID is the ID of the default project of the user.
	DefaultProjectID string `json:"default_project_id,omitempty"`

	// DomainID is the ID of the domain the user belongs to.
	DomainID string `json:"domain_id,omitempty"`

	// Enabled sets the user status to enabled or disabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Password is the password of the new user.
	Password string `json:"password,omitempty"`

	// Description is a description of the user.
	Description string `json:"description,omitempty"`
}

// ToUserCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToUserCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "user")
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Create creates a new User.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToUserCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), &b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{201},
	})
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateOptsBuilder interface {
	ToUserUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts provides options for updating a user account.
type UpdateOpts struct {
	// Name is the name of the new user.
	Name string `json:"name,omitempty"`

	// DefaultProjectID is the ID of the default project of the user.
	DefaultProjectID string `json:"default_project_id,omitempty"`

	// DomainID is the ID of the domain the user belongs to.
	DomainID string `json:"domain_id,omitempty"`

	// Enabled sets the user status to enabled or disabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Password is the password of the new user.
	Password string `json:"password,omitempty"`

	// Description is a description of the user.
	Description string `json:"description,omitempty"`
}

// ToUserUpdateMap formats a UpdateOpts into an update request.
func (opts UpdateOpts) ToUserUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "user")
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Update updates an existing User.
func Update(client *golangsdk.ServiceClient, userID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToUserUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Patch(updateURL(client, userID), &b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// ExtendedUpdateOptsBuilder allows extensions to add additional parameters to
// the ExtendedUpdate request.
type ExtendedUpdateOptsBuilder interface {
	ToUserUpdateMap() (map[string]interface{}, error)
}

// ExtendedUpdateOpts allows modifying User information (including e-mail address and mobile number)
type ExtendedUpdateOpts struct {

	// Name is the name of the user.
	Name string `json:"name,omitempty"`

	/*Password of the user. The password must meet the following requirements:
	  - Can contain 6 to 32 characters. The default minimum password length is 6 characters.
	  - Must contain at least two of the following character types: uppercase letters, lowercase letters, digits, and special characters.
	  - Must meet the requirements of the password policy configured on the account settings page.
	  - Must be different from the old password.
	*/
	Password string `json:"password,omitempty"`

	// Enabled is whether or not the user is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Description is a description of the user.
	Description string `json:"description,omitempty"`

	// Email is the email of the user
	Email string `json:"email,omitempty"`

	// AreaCode is country code
	AreaCode string `json:"areacode,omitempty"`

	// Phone is mobile number, which can contain a maximum of 32 digits.
	// The mobile number must be used together with a country code.
	Phone string `json:"phone,omitempty"`

	// Whether password reset is required at first login
	PwdResetRequired *bool `json:"pwd_status,omitempty"`

	// XUserType is Type of the IAM user in the external system.
	XUserType string `json:"xuser_type,omitempty"`

	// XUserID is ID of the IAM user in the external system.
	XUserID string `json:"xuser_id,omitempty"`
}

func (opts ExtendedUpdateOpts) ToUserUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "user")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func ExtendedUpdate(client *golangsdk.ServiceClient, userID string, opts ExtendedUpdateOpts) (r UpdateExtendedResult) {
	b, err := opts.ToUserUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(updateExtendedURL(client, userID), &b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete deletes a user.
func Delete(client *golangsdk.ServiceClient, userID string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, userID), nil)
	return
}

// ListGroups enumerates groups user belongs to.
func ListGroups(client *golangsdk.ServiceClient, userID string) pagination.Pager {
	url := listGroupsURL(client, userID)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return groups.GroupPage{LinkedPageBase: pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListProjects enumerates groups user belongs to.
func ListProjects(client *golangsdk.ServiceClient, userID string) pagination.Pager {
	url := listProjectsURL(client, userID)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return projects.ProjectPage{LinkedPageBase: pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListInGroup enumerates users that belong to a group.
func ListInGroup(client *golangsdk.ServiceClient, groupID string, opts ListOptsBuilder) pagination.Pager {
	url := listInGroupURL(client, groupID)
	if opts != nil {
		query, err := opts.ToUserListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return UserPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Add a user into one group
func AddToGroup(client *golangsdk.ServiceClient, groupID string, userID string) (r AddMembershipResult) {
	_, r.Err = client.Put(membershipURL(client, groupID, userID), nil, nil, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}

// Remove user from group
func RemoveFromGroup(client *golangsdk.ServiceClient, groupID string, userID string) (r DeleteResult) {
	_, r.Err = client.Delete(membershipURL(client, groupID, userID), nil)
	return
}
