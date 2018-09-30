package hosts

import (
	"reflect"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// AllocateOptsBuilder allows extensions to add additional parameters to the
// Allocate request.
type AllocateOptsBuilder interface {
	ToDeHAllocateMap() (map[string]interface{}, error)
}

// AllocateOpts contains all the values needed to allocate a new DeH.
type AllocateOpts struct {
	Name          string `json:"name" required:"true"`
	Az            string `json:"availability_zone" required:"true"`
	AutoPlacement string `json:"auto_placement,omitempty"`
	HostType      string `json:"host_type" required:"true"`
	Quantity      int    `json:"quantity" required:"true"`
}

// ToDeHAllocateMap builds a allocate request body from AllocateOpts.
func (opts AllocateOpts) ToDeHAllocateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Allocate accepts a AllocateOpts struct and uses the values to allocate a new DeH.
func Allocate(c *golangsdk.ServiceClient, opts AllocateOptsBuilder) (r AllocateResult) {
	b, err := opts.ToDeHAllocateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200, 201}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToDeHUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains all the values needed to update a DeH.
type UpdateOpts struct {
	Name          string `json:"name,omitempty"`
	AutoPlacement string `json:"auto_placement,omitempty"`
}

// ToDeHUpdateMap builds a update request body from UpdateOpts.
func (opts UpdateOpts) ToDeHUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "dedicated_host")
}

// Update accepts a UpdateOpts struct and uses the values to update a DeH.The response code from api is 204
func Update(c *golangsdk.ServiceClient, hostID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToDeHUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{204}}
	_, r.Err = c.Put(resourceURL(c, hostID), b, nil, reqOpt)
	return
}

//Deletes the DeH using the specified hostID.
func Delete(c *golangsdk.ServiceClient, hostid string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, hostid), nil)
	return
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API.
type ListOpts struct {
	// Specifies Dedicated Host ID.
	ID string `q:"dedicated_host_id"`
	// Specifies the Dedicated Host name.
	Name string `q:"name"`
	// Specifes the Dedicated Host type.
	HostType string `q:"host_type"`
	// Specifes the Dedicated Host name of type.
	HostTypeName string `q:"host_type_name"`
	// Specifies flavor ID.
	Flavor string `q:"flavor"`
	// Specifies the Dedicated Host status.
	// The value can be available, fault or released.
	State string `q:"state"`
	// Specifies the AZ to which the Dedicated Host belongs.
	Az string `q:"availability_zone"`
	// Specifies the number of entries displayed on each page.
	Limit string `q:"limit"`
	// 	The value is the ID of the last record on the previous page.
	Marker string `q:"marker"`
	// Filters the response by a date and time stamp when the dedicated host last changed status.
	ChangesSince string `q:"changes-since"`
	// Specifies the UUID of the tenant in a multi-tenancy cloud.
	TenantId string `q:"tenant"`
}

// ListOptsBuilder allows extensions to add parameters to the List request.
type ListOptsBuilder interface {
	ToHostListQuery() (string, error)
}

// ToRegionListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToHostListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// dedicated hosts resources. It accepts a ListOpts struct, which allows you to
// filter the returned collection for greater efficiency.
func List(c *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToHostListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return HostPage{pagination.LinkedPageBase{PageResult: r}}
	})

}

// Get retrieves a particular host based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// ListServerOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned. Marker and Limit are used
// for pagination.
type ListServerOpts struct {
	// Specifies the number of entries displayed on each page.
	Limit int `q:"limit"`
	// The value is the ID of the last record on the previous page.
	// If the marker value is invalid, error code 400 will be returned.
	Marker string `q:"marker"`
	// ID uniquely identifies this server amongst all other servers,
	// including those not accessible to the current tenant.
	ID string `json:"id"`
	// Name contains the human-readable name for the server.
	Name string `json:"name"`
	// Status contains the current operational status of the server,
	// such as IN_PROGRESS or ACTIVE.
	Status string `json:"status"`
	// UserID uniquely identifies the user account owning the tenant.
	UserID string `json:"user_id"`
}

// ListServer returns a Pager which allows you to iterate over a collection of
// dedicated hosts Server resources. It accepts a ListServerOpts struct, which allows you to
// filter the returned collection for greater efficiency.
func ListServer(c *golangsdk.ServiceClient, id string, opts ListServerOpts) ([]Server, error) {
	q, err := golangsdk.BuildQueryString(&opts)
	if err != nil {
		return nil, err
	}
	u := listServerURL(c, id) + q.String()
	pages, err := pagination.NewPager(c, u, func(r pagination.PageResult) pagination.Page {
		return ServerPage{pagination.LinkedPageBase{PageResult: r}}
	}).AllPages()

	allservers, err := ExtractServers(pages)
	if err != nil {
		return nil, err
	}

	return FilterServers(allservers, opts)
}

func FilterServers(servers []Server, opts ListServerOpts) ([]Server, error) {

	var refinedServers []Server
	var matched bool
	m := map[string]interface{}{}

	if opts.ID != "" {
		m["ID"] = opts.ID
	}
	if opts.Name != "" {
		m["Name"] = opts.Name
	}
	if opts.Status != "" {
		m["Status"] = opts.Status
	}
	if opts.UserID != "" {
		m["UserID"] = opts.UserID
	}

	if len(m) > 0 && len(servers) > 0 {
		for _, server := range servers {
			matched = true

			for key, value := range m {
				if sVal := getStructServerField(&server, key); !(sVal == value) {
					matched = false
				}
			}

			if matched {
				refinedServers = append(refinedServers, server)
			}
		}

	} else {
		refinedServers = servers
	}

	return refinedServers, nil
}

func getStructServerField(v *Server, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}
