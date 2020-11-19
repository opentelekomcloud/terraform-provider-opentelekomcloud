package policies

import (
	"reflect"

	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/pagination"
)

type ListOpts struct {
	ID         string `json:"id"`
	Name       string `q:"name"`
	Sort       string `q:"sort"`
	Status     string `json:"status"`
	Limit      int    `q:"limit"`
	Marker     string `q:"marker"`
	Offset     int    `q:"offset"`
	AllTenants string `q:"all_tenants"`
}

// List returns a Pager which allows you to iterate over a collection of
// backup policies. It accepts a ListOpts struct, which allows you to
// filter the returned collection for greater efficiency.
func List(client *golangsdk.ServiceClient, opts ListOpts) ([]BackupPolicy, error) {
	query, err := golangsdk.BuildQueryString(&opts)
	if err != nil {
		return nil, err
	}
	url := rootURL(client) + query.String()
	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return BackupPolicyPage{pagination.LinkedPageBase{PageResult: r}}
	}).AllPages()
	if err != nil {
		return nil, err
	}
	policies, err := ExtractBackupPolicies(pages)
	if err != nil {
		return nil, err
	}

	return FilterPolicies(policies, opts)
}

func FilterPolicies(policies []BackupPolicy, opts ListOpts) ([]BackupPolicy, error) {
	var refinedPolicies []BackupPolicy
	var matched bool
	m := map[string]interface{}{}

	if opts.ID != "" {
		m["ID"] = opts.ID
	}
	if opts.Status != "" {
		m["Status"] = opts.Status
	}
	if len(m) > 0 && len(policies) > 0 {
		for _, policy := range policies {
			matched = true

			for key, value := range m {
				if sVal := getStructPolicyField(&policy, key); !(sVal == value) {
					matched = false
				}
			}

			if matched {
				refinedPolicies = append(refinedPolicies, policy)
			}
		}
	} else {
		refinedPolicies = policies
	}
	return refinedPolicies, nil
}

func getStructPolicyField(v *BackupPolicy, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToBackupPolicyCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains the options for create a Backup Policy. This object is
// passed to policies.Create().
type CreateOpts struct {
	Description         string               `json:"description,omitempty"`
	Name                string               `json:"name" required:"true"`
	Parameters          PolicyParam          `json:"parameters" required:"true"`
	ProviderId          string               `json:"provider_id" required:"true"`
	Resources           []Resource           `json:"resources" required:"true"`
	ScheduledOperations []ScheduledOperation `json:"scheduled_operations" required:"true"`
	Tags                []ResourceTag        `json:"tags,omitempty"`
}

type PolicyParam struct {
	Common interface{} `json:"common,omitempty"`
}

type Resource struct {
	Id        string      `json:"id" required:"true"`
	Type      string      `json:"type" required:"true"`
	Name      string      `json:"name" required:"true"`
	ExtraInfo interface{} `json:"extra_info,omitempty"`
}

type ScheduledOperation struct {
	Description         string              `json:"description,omitempty"`
	Enabled             bool                `json:"enabled"`
	Name                string              `json:"name,omitempty"`
	OperationType       string              `json:"operation_type" required:"true"`
	OperationDefinition OperationDefinition `json:"operation_definition" required:"true"`
	Trigger             Trigger             `json:"trigger" required:"true"`
}

type OperationDefinition struct {
	MaxBackups            int    `json:"max_backups,omitempty"`
	RetentionDurationDays int    `json:"retention_duration_days,omitempty"`
	Permanent             bool   `json:"permanent"`
	PlanId                string `json:"plan_id,omitempty"`
	ProviderId            string `json:"provider_id,omitempty"`
	DayBackups            int    `json:"day_backups,omitempty"`
	WeekBackups           int    `json:"week_backups,omitempty"`
	MonthBackups          int    `json:"month_backups,omitempty"`
	YearBackups           int    `json:"year_backups,omitempty"`
	TimeZone              string `json:"timezone,omitempty"`
}

type Trigger struct {
	Properties TriggerProperties `json:"properties" required:"true"`
}

type TriggerProperties struct {
	Pattern string `json:"pattern" required:"true"`
}

type ResourceTag struct {
	Key   string `json:"key" required:"true"`
	Value string `json:"value" required:"true"`
}

// ToBackupPolicyCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToBackupPolicyCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "policy")
}

// Create will create a new backup policy based on the values in CreateOpts. To extract
// the Backup object from the response, call the Extract method on the
// CreateResult.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToBackupPolicyCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get will get a single backup policy with specific ID.
// call the Extract method on the GetResult.
func Get(client *golangsdk.ServiceClient, policyId string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, policyId), &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})

	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToPoliciesUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains the values used when updating a backup policy.
type UpdateOpts struct {
	Description         string                       `json:"description,omitempty"`
	Name                string                       `json:"name,omitempty"`
	Parameters          PolicyParam                  `json:"parameters,omitempty"`
	Resources           []Resource                   `json:"resources,omitempty"`
	ScheduledOperations []ScheduledOperationToUpdate `json:"scheduled_operations,omitempty"`
}

type ScheduledOperationToUpdate struct {
	Description         string              `json:"description,omitempty"`
	Enabled             bool                `json:"enabled"`
	TriggerId           string              `json:"trigger_id,omitempty"`
	Name                string              `json:"name,omitempty"`
	OperationDefinition OperationDefinition `json:"operation_definition,omitempty"`
	Trigger             Trigger             `json:"trigger,omitempty"`
	Id                  string              `json:"id" required:"true"`
}

// ToPoliciesUpdateMap builds an update body based on UpdateOpts.
func (opts UpdateOpts) ToPoliciesUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "policy")
}

// Update allows backup policies to be updated.
// call the Extract method on the UpdateResult.
func Update(c *golangsdk.ServiceClient, policyId string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToPoliciesUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, policyId), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will delete an existing backup policy.
func Delete(client *golangsdk.ServiceClient, policyId string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, policyId), &golangsdk.RequestOpts{
		OkCodes:      []int{200},
		JSONResponse: nil,
	})
	return
}
