package backendmember

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"fmt"
)

type LoadBalancerID struct {
	ID string `json:"id"`
}

// Backend is the primary load balancing configuration object that specifies
// the loadbalancer and port on which client traffic is received, as well
// as other details such as the load balancing method to be use, protocol, etc.
type Backend struct {
	// uri
	Uri string `json:"uri"`
	// Owner of the Listener. Only an admin user can specify a tenant ID other than its own.
	JobId string `json:"job_id"`
}

// ListenerPage is the page returned by a pager when traversing over a
// collection of routers.
type BackendPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of routers has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r BackendPage) NextPageURL() (string, error) {
	return "", nil
}

// IsEmpty checks whether a RouterPage struct is empty.
func (r BackendPage) IsEmpty() (bool, error) {
	is, err := ExtractBackend(r)
	return len(is) == 0, err
}

// ExtractBackend accepts a Page struct, specifically a ListenerPage struct,
// and extracts the elements into a slice of Listener structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractBackend(r pagination.Page) ([]Backend, error) {
	var Backends []Backend 
	err := (r.(BackendPage)).ExtractInto(&Backends)
	return Backends, err
}

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a router.
func (r commonResult) Extract() (*Backend, error) {
	fmt.Printf("Extracting Backend...\n")
	l := new(Backend)
	err := r.ExtractInto(l)
	if err != nil {
		fmt.Printf("Error: %s.\n", err.Error())
		return nil, err
	} else {
		fmt.Printf("Returning extract: %+v.\n", l)
		return l, nil
	}
}


// CreateResult represents the result of a create operation.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}
