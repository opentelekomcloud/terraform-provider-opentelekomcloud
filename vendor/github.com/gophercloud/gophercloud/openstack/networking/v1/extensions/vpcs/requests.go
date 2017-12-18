package vpcs

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the floating IP attributes you want to see returned. SortKey allows you to
// sort by a particular network attribute. SortDir sets the direction, and is
// either `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	// ID is the unique identifier for the vpc.
	ID string `json:"id"`

	// Name is the human readable name for the vpc. It does not have to be
	// unique.
	Name string `json:"name"`

	//Specifies the range of available subnets in the VPC.
	CIDR string `json:"cidr"`

	// Status indicates whether or not a vpc is currently operational.
	Status string `json:"status"`

	// Routes are a collection of static routes that the vpc will host.
	Routes []Route `json:"routes"`

	TenantID     string `q:"tenant_id"`
}

// List returns a Pager which allows you to iterate over a collection of
// vpcs. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those vpcs that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(c *gophercloud.ServiceClient, opts ListOpts) pagination.Pager {
	q, err := gophercloud.BuildQueryString(&opts)
	if err != nil {
		return pagination.Pager{Err: err}
	}
	u := rootURL(c) + q.String()
	return pagination.NewPager(c, u, func(r pagination.PageResult) pagination.Page {
		return VpcPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToVpcCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new vpc. There are
// no required values.
type CreateOpts struct {

	Name string `json:"name,omitempty"`
	CIDR string `json:"cidr,omitempty"`
}

// ToVpcCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToVpcCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "vpc")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// logical vpc. When it is created, the vpc does not have an internal
// interface - it is not associated to any subnet.
//
// You can optionally specify an external gateway for a vpc using the
// GatewayInfo struct. The external gateway for the vpc must be plugged into
// an external network (it is external if its `vpc:external' field is set to
// true).
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToVpcCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}

// Get retrieves a particular vpc based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToVpcUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains the values used when updating a vpc.
type UpdateOpts struct {
	CIDR string `json:"cidr,omitempty"`
	Name string `json:"name,omitempty"`
}

// ToVpcUpdateMap builds an update body based on UpdateOpts.
func (opts UpdateOpts) ToVpcUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "vpc")
}

// Update allows vpcs to be updated. You can update the name, administrative
// state, and the external gateway. For more information about how to set the
// external gateway for a vpc, see Create. This operation does not enable
// the update of vpc interfaces. To do this, use the AddInterface and
// RemoveInterface functions.
func Update(c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToVpcUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
// Delete will permanently delete a particular vpc based on its unique ID.
func Delete(c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}


