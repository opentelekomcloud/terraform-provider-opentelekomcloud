package servers

import (
	"reflect"
	"strings"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// SortDir is a type for specifying in which direction to sort a list of servers.
type SortDir string

// SortKey is a type for specifying by which key to sort a list of servers.
type SortKey string

var (
	// SortAsc is used to sort a list of servers in ascending order.
	SortAsc SortDir = "asc"
	// SortDesc is used to sort a list of servers in descending order.
	SortDesc SortDir = "desc"
	// SortUUID is used to sort a list of servers by uuid.
	SortUUID SortKey = "uuid"
	// SortVMState is used to sort a list of servers by vm_state.
	SortVMState SortKey = "vm_state"
	// SortDisplayName is used to sort a list of servers by display_name.
	SortDisplayName SortKey = "display_name"
	// SortTaskState is used to sort a list of servers by task_state.
	SortTaskState SortKey = "task_state"
	// SortPowerState is used to sort a list of servers by power_state.
	SortPowerState SortKey = "power_state"
	// SortAvailabilityZone is used to sort a list of servers by availability_zone.
	SortAvailabilityZone SortKey = "availability_zone"
)

// ListServerOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned.
type ListOpts struct {
	// ID uniquely identifies this server amongst all other servers,
	// including those not accessible to the current tenant.
	ID string
	//ID of the user to which the BMS belongs.
	UserID string
	//Contains the nova-compute status
	HostStatus string
	//Contains the host ID of the BMS.
	HostID string
	// KeyName indicates which public key was injected into the server on launch.
	KeyName string
	// Specifies the BMS name, not added in query since returns like results.
	Name string
	// Specifies the BMS image ID.
	ImageID string `q:"image"`
	// Specifies flavor ID.
	FlavorID string `q:"flavor"`
	// Specifies the BMS status.
	Status string `q:"status"`
	//Filters out the BMSs that have been updated since the changes-since time.
	// The parameter is in ISO 8601 time format, for example, 2013-06-09T06:42:18Z.
	ChangesSince string `q:"changes-since"`
	//Specifies whether to query the BMSs of all tenants. This parameter is available only to administrators.
	// The value can be 0 (do not query the BMSs of all tenants) or 1 (query the BMSs of all tenants).
	AllTenants int `q:"all_tenants"`
	//Specifies the IP address. This parameter supports fuzzy matching.
	IP string `q:"ip"`
	//Specifies the tag list. Returns BMSs that match all tags. Use commas (,) to separate multiple tags
	Tags string `q:"tags"`
	//Specifies the tag list. Returns BMSs that match any tag
	TagsAny string `q:"tags-any"`
	//Specifies the tag list. Returns BMSs that do not match all tags.
	NotTags string `q:"not-tags"`
	//Specifies the tag list. Returns BMSs that do not match any of the tags.
	NotTagsAny string `q:"not-tags-any"`
	//Specifies the BMS sorting attribute, which can be the BMS UUID (uuid), BMS status (vm_state),
	// BMS name (display_name), BMS task status (task_state), power status (power_state),
	// creation time (created_at), last time when the BMS is updated (updated_at), and availability zone
	// (availability_zone). You can specify multiple sort_key and sort_dir pairs.
	SortKey SortKey `q:"sort_key"`
	//Specifies the sorting direction, i.e. asc or desc.
	SortDir SortDir `q:"sort_dir"`
}

// ListServer returns a Pager which allows you to iterate over a collection of
// BMS Server resources. It accepts a ListServerOpts struct, which allows you to
// filter the returned collection for greater efficiency.
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]Server, error) {
	c.Microversion = "2.26"
	q, err := golangsdk.BuildQueryString(&opts)
	if err != nil {
		return nil, err
	}
	u := listURL(c) + q.String()
	pages, err := pagination.NewPager(c, u, func(r pagination.PageResult) pagination.Page {
		return ServerPage{pagination.LinkedPageBase{PageResult: r}}
	}).AllPages()

	allservers, err := ExtractServers(pages)
	if err != nil {
		return nil, err
	}
	return FilterServers(allservers, opts)
}

func FilterServers(servers []Server, opts ListOpts) ([]Server, error) {
	var refinedServers []Server
	var matched bool
	m := map[string]interface{}{}

	if opts.ID != "" {
		m["ID"] = opts.ID
	}
	if opts.Name != "" {
		m["Name"] = opts.Name
	}
	if opts.UserID != "" {
		m["UserID"] = opts.UserID
	}
	if opts.HostStatus != "" {
		m["HostStatus"] = opts.HostStatus
	}
	if opts.HostID != "" {
		m["HostID"] = opts.HostID
	}
	if opts.KeyName != "" {
		m["KeyName"] = opts.KeyName
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
	var serverList []Server

	for i := 0; i < len(refinedServers); i++ {
		if strings.Contains(refinedServers[i].Flavor.ID, "physical") {
			serverList = append(serverList, refinedServers[i])
		}

	}
	return serverList, nil
}

func getStructServerField(v *Server, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

// Get requests details on a single server, by ID.
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{"X-OpenStack-Nova-API-Version": "2.26"},
	})
	return
}
