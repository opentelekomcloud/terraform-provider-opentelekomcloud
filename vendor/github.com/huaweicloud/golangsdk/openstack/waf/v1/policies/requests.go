package policies

import (
	"github.com/huaweicloud/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToPolicyCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new policy.
type CreateOpts struct {
	//Policy name
	Name string `json:"name" required:"true"`
}

// ToPolicyCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToPolicyCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create will create a new policy based on the values in CreateOpts.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToPolicyCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToPolicyUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains all the values needed to update a policy.
type UpdateOpts struct {
	//Policy name
	Name string `json:"name,omitempty"`
	//Protective Action
	Action *Action `json:"action,omitempty"`
	//Protection Switches
	Options *Options `json:"options,omitempty"`
	//Protection Level
	Level int `json:"level,omitempty"`
	//Detection Mode
	FullDetection *bool `json:"full_detection,omitempty"`
}

// ToPolicyUpdateMap builds a update request body from UpdateOpts.
func (opts UpdateOpts) ToPolicyUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and uses the values to update a policy.The response code from api is 200
func Update(c *golangsdk.ServiceClient, policyID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToPolicyUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Put(resourceURL(c, policyID), b, nil, reqOpt)
	return
}

// UpdateHostsOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateHostsOptsBuilder interface {
	ToPolicyHostsUpdateMap() (map[string]interface{}, error)
}

// UpdateHostsOpts contains all the values needed to update a policy hosts.
type UpdateHostsOpts struct {
	//Domain IDs
	Hosts []string `json:"hosts" required:"true"`
}

// ToPolicyHostsUpdateMap builds a update request body from UpdateHostsOpts.
func (opts UpdateHostsOpts) ToPolicyHostsUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update accepts a UpdateHostsOpts struct and uses the values to update a policy hosts.The response code from api is 200
func UpdateHosts(c *golangsdk.ServiceClient, policyID string, opts UpdateHostsOptsBuilder) (r UpdateResult) {
	b, err := opts.ToPolicyHostsUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Put(hostsURL(c, policyID), b, nil, reqOpt)
	return
}

// Get retrieves a particular policy based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, &RequestOpts)
	return
}

// Delete will permanently delete a particular policy based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{204},
		MoreHeaders: RequestOpts.MoreHeaders}
	_, r.Err = c.Delete(resourceURL(c, id), reqOpt)
	return
}
