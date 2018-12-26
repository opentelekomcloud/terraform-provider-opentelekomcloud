package defsecrules

import (
	"strings"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// List will return a collection of default rules.
func List(client *golangsdk.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, rootURL(client), func(r pagination.PageResult) pagination.Page {
		return DefaultRulePage{pagination.SinglePageBase(r)}
	})
}

// CreateOpts represents the configuration for adding a new default rule.
type CreateOpts struct {
	// The lower bound of the port range that will be opened.
	FromPort int `json:"from_port"`

	// The upper bound of the port range that will be opened.
	ToPort int `json:"to_port"`

	// The protocol type that will be allowed, e.g. TCP.
	IPProtocol string `json:"ip_protocol" required:"true"`

	// ONLY required if FromGroupID is blank. This represents the IP range that
	// will be the source of network traffic to your security group.
	//
	// Use 0.0.0.0/0 to allow all IPv4 addresses.
	// Use ::/0 to allow all IPv6 addresses.
	CIDR string `json:"cidr,omitempty"`
}

// CreateOptsBuilder builds the create rule options into a serializable format.
type CreateOptsBuilder interface {
	ToRuleCreateMap() (map[string]interface{}, error)
}

// ToRuleCreateMap builds the create rule options into a serializable format.
func (opts CreateOpts) ToRuleCreateMap() (map[string]interface{}, error) {
	if opts.FromPort == 0 && strings.ToUpper(opts.IPProtocol) != "ICMP" {
		return nil, golangsdk.ErrMissingInput{Argument: "FromPort"}
	}
	if opts.ToPort == 0 && strings.ToUpper(opts.IPProtocol) != "ICMP" {
		return nil, golangsdk.ErrMissingInput{Argument: "ToPort"}
	}
	return golangsdk.BuildRequestBody(opts, "security_group_default_rule")
}

// Create is the operation responsible for creating a new default rule.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get will return details for a particular default rule.
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, nil)
	return
}

// Delete will permanently delete a rule the project's default security group.
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, id), nil)
	return
}
