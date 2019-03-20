package hosts

import (
	"time"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type Host struct {
	// ID is the unique identifier for the dedicated host .
	ID string `json:"dedicated_host_id"`
	// Specifies the Dedicated Host name.
	Name string `json:"name"`
	// Specifies whether to allow a VM to be placed on this available host
	// if its Dedicated Host ID is not specified during its creation.
	AutoPlacement string `json:"auto_placement"`
	// Specifies the AZ to which the Dedicated Host belongs.
	Az string `json:"availability_zone"`
	// Specifies the tenant who owns the Dedicated Host.
	TenantId string `json:"project_id"`
	// Specifies the host status.
	State string `json:"state"`
	// Specifies the number of available vCPUs for the Dedicated Host.
	AvailableVcpus int `json:"available_vcpus"`
	// 	Specifies the size of available memory for the Dedicated Host.
	AvailableMemory int `json:"available_memory"`
	// Time at which the dedicated host has been allocated.
	AllocatedAt string `json:"allocated_at"`
	// Time at which the dedicated host has been released.
	ReleasedAt string `json:"released_at"`
	// Specifies the number of the placed VMs.
	InstanceTotal int `json:"instance_total"`
	// Specifies the VMs started on the Dedicated Host.
	InstanceUuids []string `json:"instance_uuids"`
	// Specifies the property of host.
	HostProperties HostPropertiesOpts `json:"host_properties"`
}
type HostPropertiesOpts struct {
	// Specifies the property of host.
	HostType           string               `json:"host_type"`
	HostTypeName       string               `json:"host_type_name"`
	Vcpus              int                  `json:"vcpus"`
	Cores              int                  `json:"cores"`
	Sockets            int                  `json:"sockets"`
	Memory             int                  `json:"memory"`
	InstanceCapacities []InstanceCapacities `json:"available_instance_capacities"`
}
type InstanceCapacities struct {
	// Specifies the number of supported flavors.
	Flavor string `json:"flavor"`
}

// HostPage is the page returned by a pager when traversing over a
// collection of Hosts.
type HostPage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if a ListResult contains no Dedicated Hosts.
func (r HostPage) IsEmpty() (bool, error) {
	stacks, err := ExtractHosts(r)
	return len(stacks) == 0, err
}

// ExtractHosts accepts a Page struct, specifically a HostPage struct,
// and extracts the elements into a slice of hosts structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractHosts(r pagination.Page) ([]Host, error) {
	var s struct {
		ListedStacks []Host `json:"dedicated_hosts"`
	}
	err := (r.(HostPage)).ExtractInto(&s)
	return s.ListedStacks, err
}

// NextPageURL is invoked when a paginated collection of hosts has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r HostPage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"dedicated_hostslinks"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Links)
}

type commonResult struct {
	golangsdk.Result
}

// AllocateResult represents the result of a allocate operation. Call its Extract
// method to interpret it as a host.
type AllocateResult struct {
	commonResult
}

// Extract is a function that accepts a result and extracts Allocated Hosts.
func (r AllocateResult) ExtractHost() (*AllocatedHosts, error) {
	var response AllocatedHosts
	err := r.ExtractInto(&response)
	return &response, err
}

//AllocatedHosts is the response structure of the allocated DeH
type AllocatedHosts struct {
	AllocatedHostIds []string `json:"dedicated_host_ids"`
}

// AllocateResult represents the result of a allocate operation. Call its Extract
// method to interpret it as a host.
type UpdateResult struct {
	commonResult
}

type DeleteResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a host.
type GetResult struct {
	commonResult
}

// Extract is a function that accepts a result and extracts a host.
func (r commonResult) Extract() (*Host, error) {
	var s struct {
		Host *Host `json:"dedicated_host"`
	}
	err := r.ExtractInto(&s)
	return s.Host, err
}

// Server represents a server/instance in the OpenStack cloud.
type Server struct {
	// ID uniquely identifies this server amongst all other servers,
	// including those not accessible to the current tenant.
	ID string `json:"id"`
	// TenantID identifies the tenant owning this server resource.
	TenantID string `json:"tenant_id"`
	// UserID uniquely identifies the user account owning the tenant.
	UserID string `json:"user_id"`
	// Name contains the human-readable name for the server.
	Name string `json:"name"`
	// Updated and Created contain ISO-8601 timestamps of when the state of the
	// server last changed, and when it was created.
	Updated time.Time `json:"updated"`
	Created time.Time `json:"created"`
	// Status contains the current operational status of the server,
	// such as IN_PROGRESS or ACTIVE.
	Status string `json:"status"`
	// Image refers to a JSON object, which itself indicates the OS image used to
	// deploy the server.
	Image map[string]interface{} `json:"-"`
	// Flavor refers to a JSON object, which itself indicates the hardware
	// configuration of the deployed server.
	Flavor map[string]interface{} `json:"flavor"`
	// Addresses includes a list of all IP addresses assigned to the server,
	// keyed by pool.
	Addresses map[string]interface{} `json:"addresses"`
	// Metadata includes a list of all user-specified key-value pairs attached
	// to the server.
	Metadata map[string]string `json:"metadata"`
}

type ServerPage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if a page contains no Server results.
func (r ServerPage) IsEmpty() (bool, error) {
	s, err := ExtractServers(r)
	return len(s) == 0, err
}

// NextPageURL uses the response's embedded link reference to navigate to the
// next page of results.
func (r ServerPage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"servers_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Links)
}

// ExtractServers accepts a Page struct, specifically a ServerPage struct,
// and extracts the elements into a slice of Server structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractServers(r pagination.Page) ([]Server, error) {
	var s struct {
		ListedStacks []Server `json:"servers"`
	}
	err := (r.(ServerPage)).ExtractInto(&s)
	return s.ListedStacks, err
}
