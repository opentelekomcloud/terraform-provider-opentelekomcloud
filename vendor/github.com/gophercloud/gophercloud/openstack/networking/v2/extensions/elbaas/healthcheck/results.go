package healthcheck

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
    "fmt"
)


// Health represents a load balancer health check. A health monitor is used
// to determine whether or not back-end members of the VIP's pool are usable
// for processing a request. A pool can have several health monitors associated
// with it. There are different types of health monitors supported:
//
// PING: used to ping the members using ICMP.
// TCP: used to connect to the members using TCP.
// HTTP: used to send an HTTP request to the member.
// HTTPS: used to send a secure HTTP request to the member.
//
// When a pool has several monitors associated with it, each member of the pool
// is monitored by all these monitors. If any monitor declares the member as
// unhealthy, then the member status is changed to INACTIVE and the member
// won't participate in its pool's load balancing. In other words, ALL monitors
// must declare the member to be healthy for it to stay ACTIVE.
type Health struct {
    // Thealthcheck_interval
    HealthcheckInterval int `json:"healthcheck_interval"`

    // listener_id
    ListenerId string `json:"listener_id"`

    // The unique ID for the health.
	ID string `json:"id"`

	// The healthcheck_ protocol
	HealthcheckProtocol string `json:"healthcheck_ protocol"`

	// unhealthy_threshold
	UnhealthyThreshold int `json:"unhealthy_threshold"`

	// update_time
	UpdateTime string `json:"update_time"`

	// create_time
	CreateTime string `json:"create_time"`

	// healthcheck_connect_port
	HealthcheckConnectPort int `json:"healthcheck_connect_port"`

	// healthcheck_timeout
	HealthcheckTimeout int `json:"healthcheck_timeout"`

	// healthcheck_uri
	HealthcheckUri string `json:"healthcheck_uri" `

	// healthy_threshold
	HealthyThreshold int `json:"healthy_threshold"`

	// The administrative state of the health monitor, which is up (true) or down (false).
	//AdminStateUp bool `json:"admin_state_up"`
}

// MonitorPage is the page returned by a pager when traversing over a
// collection of health monitors.
type HealthPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of monitors has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r HealthPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"healthcheck_links"`
	}

	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}

	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a MonitorPage struct is empty.
func (r HealthPage) IsEmpty() (bool, error) {
	is, err := ExtractHealth(r)
	return len(is) == 0, err
}

// ExtractMonitors accepts a Page struct, specifically a MonitorPage struct,
// and extracts the elements into a slice of Monitor structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractHealth(r pagination.Page) ([]Health, error) {
	var s struct {
		Healths []Health `json:"healthcheck"`
	}
	err := (r.(HealthPage)).ExtractInto(&s)
	return s.Healths, err
}

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a monitor.
func (r commonResult) Extract() (*Health, error) {
	fmt.Printf("Extracting Health...\n")
	l := new(Health)
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

// UpdateResult represents the result of an update operation.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}
