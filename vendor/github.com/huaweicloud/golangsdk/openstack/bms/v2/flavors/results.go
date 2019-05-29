package flavors

import (
	"encoding/json"
	"strconv"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// Flavor represent (virtual) hardware configurations for server resources
// in a region.
type Flavor struct {
	// ID is the flavor's unique ID.
	ID string `json:"id"`

	// Disk is the amount of root disk, measured in GB.
	Disk int `json:"disk"`

	// MinDisk and MinRAM, if provided, elides flavors which do not meet your
	// criteria.
	MinDisk int `json:"minDisk"`

	MinRAM int `json:"minRam"`

	// RAM is the amount of memory, measured in MB.
	RAM int `json:"ram"`

	// Name is the name of the flavor.
	Name string `json:"name"`

	// RxTxFactor describes bandwidth alterations of the flavor.
	RxTxFactor float64 `json:"rxtx_factor"`

	// Swap is the amount of swap space, measured in MB.
	Swap int `json:"-"`

	// VCPUs indicates how many (virtual) CPUs are available for this flavor.
	VCPUs int `json:"vcpus"`

	// IsPublic indicates whether the flavor is public.
	IsPublic bool `json:"os-flavor-access:is_public"`

	// Ephemeral is the amount of ephemeral disk space, measured in GB.
	Ephemeral int `json:"OS-FLV-EXT-DATA:ephemeral"`

	// Whether or not the flavor has been administratively disabled
	Disabled bool `json:"OS-FLV-DISABLED:disabled"`

	// Specifies the shortcut link of the BMS flavor.
	Links []golangsdk.Link `json:"links"`

	SortKey string `json:"sort_key"`

	//SortDir sets the direction, and is either `asc' or `desc'
	SortDir string `json:"sort_dir"`
}

func (r *Flavor) UnmarshalJSON(b []byte) error {
	type tmp Flavor
	var s struct {
		tmp
		Swap interface{} `json:"swap"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*r = Flavor(s.tmp)

	switch t := s.Swap.(type) {
	case float64:
		r.Swap = int(t)
	case string:
		switch t {
		case "":
			r.Swap = 0
		default:
			swap, err := strconv.ParseFloat(t, 64)
			if err != nil {
				return err
			}
			r.Swap = int(swap)
		}
	}

	return nil
}

// FlavorPage contains a single page of all flavors from a List call.
type FlavorPage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines if a FlavorPage contains any results.
func (page FlavorPage) IsEmpty() (bool, error) {
	flavors, err := ExtractFlavors(page)
	return len(flavors) == 0, err
}

// NextPageURL uses the response's embedded link reference to navigate to the
// next page of results.
func (page FlavorPage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"flavors_links"`
	}
	err := page.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Links)
}

// ExtractFlavors provides access to the list of flavors in a page acquired
// from the List operation.
func ExtractFlavors(r pagination.Page) ([]Flavor, error) {
	var s struct {
		Flavors []Flavor `json:"flavors"`
	}
	err := (r.(FlavorPage)).ExtractInto(&s)
	return s.Flavors, err
}
