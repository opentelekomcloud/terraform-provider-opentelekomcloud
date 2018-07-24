package servers

import (
	"time"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type commonResult struct {
	golangsdk.Result
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

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a server.
type GetResult struct {
	commonResult
}

// Extract is a function that accepts a result and extracts a server.
func (r commonResult) Extract() (*Server, error) {
	var s struct {
		Server *Server `json:"server"`
	}
	err := r.ExtractInto(&s)
	return s.Server, err
}

// ExtractServers accepts a Page struct, specifically a ServerPage struct,
// and extracts the elements into a slice of Server structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractServers(r pagination.Page) ([]Server, error) {
	var s struct {
		Servers []Server `json:"servers"`
	}
	err := (r.(ServerPage)).ExtractInto(&s)
	return s.Servers, err
}

// Server exposes fields corresponding to a given server on the user's account.
type Server struct {
	// ID uniquely identifies this server amongst all other servers, including those not accessible to the current tenant.
	ID string `json:"id"`
	// TenantID identifies the tenant owning this server resource.
	TenantID string `json:"tenant_id"`
	// UserID uniquely identifies the user account owning the tenant.
	UserID string `json:"user_id"`
	// Name contains the human-readable name for the server.
	Name string `json:"name"`
	// Status contains the current operational status of the server, such as IN_PROGRESS or ACTIVE.
	Status string `json:"status"`
	// Updated and Created contain ISO-8601 timestamps of when the state of the server last changed, and when it was created.
	Updated time.Time `json:"updated"`
	Created time.Time `json:"created"`
	//Specifies the nova-compute status.
	HostStatus string `json:"host_status"`
	//Specifies the host ID of the BMS.
	HostID string `json:"hostid"`
	// Progress ranges from 0..100.
	// A request made against the server completes only once Progress reaches 100.
	Progress int `json:"progress"`
	// AccessIPv4 and AccessIPv6 contain the IP addresses of the server, suitable for remote access for administration.
	AccessIPv4 string `json:"accessIPv4"`
	AccessIPv6 string `json:"accessIPv6"`
	// Image refers to a JSON object, which itself indicates the OS image used to deploy the server.
	Image Images `json:"image"`
	// Flavor refers to a JSON object, which itself indicates the hardware configuration of the deployed server.
	Flavor Flavor `json:"flavor"`
	// Addresses includes a list of all IP addresses assigned to the server, keyed by pool.
	Addresses map[string]interface{} `json:"addresses"`
	// Metadata includes a list of all user-specified key-value pairs attached to the server.
	Metadata map[string]string `json:"metadata"`
	// Links includes HTTP references to the itself, useful for passing along to other APIs that might want a server reference.
	Links []Links `json:"links"`
	// KeyName indicates which public key was injected into the server on launch.
	KeyName string `json:"key_name"`
	// AdminPass will generally be empty.  However, it will contain the administrative password chosen when provisioning a new server without a set AdminPass setting in the first place.
	// Note that this is the ONLY time this field will be valid.
	AdminPass string `json:"adminPass"`
	// SecurityGroups includes the security groups that this instance has applied to it
	SecurityGroups []SecurityGroups `json:"security_groups"`
	//Specifies the BMS tag.
	//Added in micro version 2.26.
	Tags []string `json:"tags"`
	//Specifies whether a BMS is locked
	Locked      bool   `json:"locked"`
	ConfigDrive string `json:"config_drive"`
	//Specifies the AZ ID. This is an extended attribute.
	AvailabilityZone string `json:"OS-EXT-AZ:availability_zone"`
	//Specifies the disk configuration mode. This is an extended attribute.
	DiskConfig string `json:"OS-DCF:diskConfig"`
	//Specifies the name of a host on the hypervisor.
	// It is an extended attribute provided by the Nova driver
	HostName string `json:"OS-EXT-SRV-ATTR:hostname"`
	//Specifies the server description.
	Description string `json:"description"`
	//Specifies the job status of the BMS. This is an extended attribute.
	TaskState string `json:"OS-EXT-STS:task_state"`
	//Specifies the power status of the BMS. This is an extended attribute
	PowerState int `json:"OS-EXT-STS:power_state"`
	//Specifies the UUID of the kernel image when the AMI image is used
	KernelId string `json:"OS-EXT-SRV-ATTR:kernel_id"`
	//Specifies the host name of the BMS. This is an extended attribute
	Host string `json:"OS-EXT-SRV-ATTR:host"`
	//Specifies the UUID of the Ramdisk image when the AMI image is used.
	RamdiskId string `json:"OS-EXT-SRV-ATTR:ramdisk_id"`
	//Specifies the BMS startup sequence in the batch BMS creation scenario.
	Launch_index int `json:"OS-EXT-SRV-ATTR:launch_index"`
	//Specifies the user data specified during BMS creation.
	UserData string `json:"OS-EXT-SRV-ATTR:user_data"`
	//Specifies the reserved BMS IDs in the batch BMS creation scenario.
	ReservationID string `json:"OS-EXT-SRV-ATTR:reservation_id"`
	//Specifies the device name of the BMS system disk
	RootDevicName string `json:"OS-EXT-SRV-ATTR:root_device_name"`
	//Specifies the name of a host on the hypervisor.
	HypervisorHostName string `json:"OS-EXT-SRV-ATTR:hypervisor_hostname"`
	//Specifies the BMS status. This is an extended attribute.
	VMState string `json:"OS-EXT-STS:vm_state"`
	//Specifies the BMS ID. This is an extended attribute.
	InstanceName string `json:"OS-EXT-SRV-ATTR:instance_name"`
}

type SecurityGroups struct {
	Name string `json:"name"`
}

type Flavor struct {
	ID    string  `json:"id"`
	Links []Links `json:"links"`
}

type Links struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type Images struct {
	ID    string  `json:"id"`
	Links []Links `json:"links"`
}
