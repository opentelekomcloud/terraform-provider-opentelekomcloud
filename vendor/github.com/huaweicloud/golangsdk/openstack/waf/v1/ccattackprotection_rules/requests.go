package ccattackprotection_rules

import (
	"github.com/huaweicloud/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToCcAttackCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new cc attack protection rule.
type CreateOpts struct {
	Url          string       `json:"url" required:"true"`
	LimitNum     *int         `json:"limit_num" required:"true"`
	LimitPeriod  *int         `json:"limit_period" required:"true"`
	LockTime     *int         `json:"lock_time,omitempty"`
	TagType      string       `json:"tag_type" required:"true"`
	TagIndex     string       `json:"tag_index,omitempty"`
	TagCondition TagCondition `json:"tag_condition,omitempty"`
	Action       Action       `json:"action" required:"true"`
}

type TagCondition struct {
	Category string   `json:"category" required:"true"`
	Contents []string `json:"contents" required:"true"`
}

type Action struct {
	Category string `json:"category" required:"true"`
	Detail   Detail `json:"detail,omitempty"`
}

type Detail struct {
	Response Response `json:"response" required:"true"`
}

type Response struct {
	ContentType string `json:"content_type" required:"true"`
	Content     string `json:"content" required:"true"`
}

// ToCcAttackCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToCcAttackCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create will create a new cc attack protection rule based on the values in CreateOpts.
func Create(c *golangsdk.ServiceClient, policyID string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToCcAttackCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c, policyID), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular cc attack rule based on its unique ID.
func Get(c *golangsdk.ServiceClient, policyID, ruleID string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, policyID, ruleID), &r.Body, &RequestOpts)
	return
}

// Delete will permanently delete a particular cc attack rule based on its unique ID.
func Delete(c *golangsdk.ServiceClient, policyID, ruleID string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{204},
		MoreHeaders: RequestOpts.MoreHeaders}
	_, r.Err = c.Delete(resourceURL(c, policyID, ruleID), reqOpt)
	return
}
