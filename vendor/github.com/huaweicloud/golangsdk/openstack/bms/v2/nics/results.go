package nics

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type FixedIP struct {
	SubnetID  string `json:"subnet_id"`
	IPAddress string `json:"ip_address"`
}

//Nic Manage and perform other operations on Nic,
//including querying Nics as well as
//querying Nic.
type Nic struct {
	// ID is the unique identifier for the nic.
	ID string `json:"port_id"`

	//Specifies the ID of the network to which the NIC port belongs.
	NetworkID string `json:"net_id"`

	// Status indicates whether or not a nic is currently operational.
	Status string `json:"port_state"`

	//Specifies the NIC private IP address.
	FixedIP []FixedIP `json:"fixed_ips"`

	//Specifies the MAC address of the NIC.
	MACAddress string `json:"mac_addr"`
}

// NicPage is the page returned by a pager when traversing over a
// collection of nics.
type NicPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of nics has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r NicPage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"interfaceAttachments_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a NicPage struct is empty.
func (r NicPage) IsEmpty() (bool, error) {
	is, err := ExtractNics(r)
	return len(is) == 0, err
}

// ExtractNics accepts a Page struct, specifically a NicPage struct,
// and extracts the elements into a slice of Nic structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractNics(r pagination.Page) ([]Nic, error) {
	var s struct {
		Nics []Nic `json:"interfaceAttachments"`
	}
	err := (r.(NicPage)).ExtractInto(&s)
	return s.Nics, err
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a nic.
func (r commonResult) Extract() (*Nic, error) {
	var s struct {
		Nic *Nic `json:"interfaceAttachment"`
	}
	err := r.ExtractInto(&s)
	return s.Nic, err
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Nic.
type GetResult struct {
	commonResult
}
