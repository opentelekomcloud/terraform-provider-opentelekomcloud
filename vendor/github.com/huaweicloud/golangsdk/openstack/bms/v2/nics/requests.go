package nics

import (
	"reflect"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the nic attributes you want to see returned.
type ListOpts struct {
	// ID is the unique identifier for the nic.
	ID string `json:"port_id"`

	// Status indicates whether or not a nic is currently operational.
	Status string `json:"port_state"`
}

// List returns collection of
// nics. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those nics that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(c *golangsdk.ServiceClient, serverId string, opts ListOpts) ([]Nic, error) {
	u := listURL(c, serverId)
	pages, err := pagination.NewPager(c, u, func(r pagination.PageResult) pagination.Page {
		return NicPage{pagination.LinkedPageBase{PageResult: r}}
	}).AllPages()

	allNICs, err := ExtractNics(pages)
	if err != nil {
		return nil, err
	}

	return FilterNICs(allNICs, opts)
}

func FilterNICs(nics []Nic, opts ListOpts) ([]Nic, error) {

	var refinedNICs []Nic
	var matched bool
	m := map[string]interface{}{}

	if opts.ID != "" {
		m["ID"] = opts.ID
	}
	if opts.Status != "" {
		m["Status"] = opts.Status
	}
	if len(m) > 0 && len(nics) > 0 {
		for _, nic := range nics {
			matched = true

			for key, value := range m {
				if sVal := getStructField(&nic, key); !(sVal == value) {
					matched = false
				}
			}

			if matched {
				refinedNICs = append(refinedNICs, nic)
			}
		}

	} else {
		refinedNICs = nics
	}

	return refinedNICs, nil
}

func getStructField(v *Nic, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

// Get retrieves a particular nic based on its unique ID.
func Get(c *golangsdk.ServiceClient, serverId string, id string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, serverId, id), &r.Body, nil)
	return
}
