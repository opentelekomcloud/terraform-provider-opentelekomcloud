package flavors

import (
	"reflect"
	"strings"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToFlavorListQuery() (string, error)
}

//The AccessType arguement is optional, and if it is not supplied, OpenStack
//returns the PublicAccess flavors.
type AccessType string

const (
	// PublicAccess returns public flavors and private flavors associated with
	// that project.
	PublicAccess AccessType = "true"

	// PrivateAccess (admin only) returns private flavors, across all projects.
	PrivateAccess AccessType = "false"

	// AllAccess (admin only) returns public and private flavors across all
	// projects.
	AllAccess AccessType = "None"
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the flavor attributes you want to see returned.
type ListOpts struct {
	//Specifies the name of the BMS flavor
	Name string

	//Specifies the ID of the BMS flavor
	ID string

	// MinDisk and MinRAM, if provided, elides flavors which do not meet your
	// criteria.
	MinDisk int `q:"minDisk"`

	MinRAM int `q:"minRam"`

	// AccessType, if provided, instructs List which set of flavors to return.
	// If IsPublic not provided, flavors for the current project are returned.
	AccessType AccessType `q:"is_public"`

	//SortKey allows you to sort by a particular attribute
	SortKey string `q:"sort_key"`

	//SortDir sets the direction, and is either `asc' or `desc'
	SortDir string `q:"sort_dir"`
}

func List(c *golangsdk.ServiceClient, opts ListOpts) ([]Flavor, error) {
	q, err := golangsdk.BuildQueryString(&opts)
	if err != nil {
		return nil, err
	}
	u := listURL(c) + q.String()
	pages, err := pagination.NewPager(c, u, func(r pagination.PageResult) pagination.Page {
		return FlavorPage{pagination.LinkedPageBase{PageResult: r}}
	}).AllPages()

	allFlavors, err := ExtractFlavors(pages)
	if err != nil {
		return nil, err
	}

	return FilterFlavors(allFlavors, opts)
}

//FilterFlavors used to filter flavors using Id and Name
func FilterFlavors(flavors []Flavor, opts ListOpts) ([]Flavor, error) {

	var refinedFlavors []Flavor
	var matched bool
	m := map[string]interface{}{}

	if opts.ID != "" {
		m["ID"] = opts.ID
	}
	if opts.Name != "" {
		m["Name"] = opts.Name
	}
	if len(m) > 0 && len(flavors) > 0 {
		for _, flavor := range flavors {
			matched = true

			for key, value := range m {
				if sVal := getStructField(&flavor, key); !(sVal == value) {
					matched = false
				}
			}
			if matched {
				refinedFlavors = append(refinedFlavors, flavor)
			}
		}
	} else {
		refinedFlavors = flavors
	}
	var flavorList []Flavor

	for i := 0; i < len(refinedFlavors); i++ {
		if strings.Contains(refinedFlavors[i].Name, "physical") {
			flavorList = append(flavorList, refinedFlavors[i])
		}

	}

	return flavorList, nil
}

func getStructField(v *Flavor, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}
