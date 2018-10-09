package shares

import (
	"reflect"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// SortDir is a type for specifying in which direction to sort a list of Shares.
type SortDir string

var (
	// SortAsc is used to sort a list of Shares in ascending order.
	SortAsc SortDir = "asc"
	// SortDesc is used to sort a list of Shares in descending order.
	SortDesc SortDir = "desc"
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the share attributes you want to see returned.
type ListOpts struct {
	ID               string
	SnapshotID       string
	ShareToMe        bool    `q:"share_to_me"`
	Name             string  `q:"name"`
	Status           string  `q:"status"`
	BackupID         string  `q:"backup_id"`
	FromProjectID    string  `q:"from_project_id"`
	ToProjectID      string  `q:"to_project_id"`
	AvailabilityZone string  `q:"availability_zone"`
	SortDir          SortDir `q:"sort_dir"`
	Limit            int     `q:"limit"`
	Offset           int     `q:"offset"`
	VolumeID         string  `q:"volume_id"`
}

// List returns collection of
// share. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
// Default policy settings return only those share that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]Share, error) {
	q, err := golangsdk.BuildQueryString(&opts)
	if err != nil {
		return nil, err
	}
	u := listURL(c) + q.String()
	pages, err := pagination.NewPager(c, u, func(r pagination.PageResult) pagination.Page {
		return SharePage{pagination.LinkedPageBase{PageResult: r}}
	}).AllPages()

	allShares, err := ExtractShareList(pages)
	if err != nil {
		return nil, err
	}

	return FilterShares(allShares, opts)
}

func FilterShares(shares []Share, opts ListOpts) ([]Share, error) {

	var refinedShares []Share
	var matched bool
	m := map[string]FilterStruct{}

	if opts.ID != "" {
		m["ID"] = FilterStruct{Value: opts.ID}
	}
	if opts.SnapshotID != "" {
		m["SnapshotID"] = FilterStruct{Value: opts.SnapshotID, Driller: []string{"Backup"}}
	}

	if len(m) > 0 && len(shares) > 0 {
		for _, share := range shares {
			matched = true

			for key, value := range m {
				if sVal := GetStructNestedField(&share, key, value.Driller); !(sVal == value.Value) {
					matched = false
				}
			}
			if matched {
				refinedShares = append(refinedShares, share)
			}
		}

	} else {
		refinedShares = shares
	}

	return refinedShares, nil
}

type FilterStruct struct {
	Value   string
	Driller []string
}

func GetStructNestedField(v *Share, field string, structDriller []string) string {
	r := reflect.ValueOf(v)
	for _, drillField := range structDriller {
		f := reflect.Indirect(r).FieldByName(drillField).Interface()
		r = reflect.ValueOf(f)
	}
	f1 := reflect.Indirect(r).FieldByName(field)
	return string(f1.String())
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToShareCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new share.
type CreateOpts struct {
	//ID of the backup to be shared
	BackupID string `json:"backup_id" required:"true"`
	//IDs of projects with which the backup is shared
	ToProjectIDs []string `json:"to_project_ids" required:"true"`
}

// ToShareCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToShareCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "shared")
}

// Create will create a new Share based on the values in CreateOpts. To extract
// the Share object from the response, call the Extract method on the
// CreateResult.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToShareCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular share based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

//DeleteOptsBuilder is an interface which can be able to build the query string
//of share deletion.
type DeleteOptsBuilder interface {
	ToShareDeleteQuery() (string, error)
}

type DeleteOpts struct {
	//Whether the ID in the URL is a backup share ID or a backup ID
	IsBackupID bool `q:"is_backup_id"`
}

func (opts DeleteOpts) ToShareDeleteQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

//Delete is a method by which can be able to delete one or all shares of a backup.
func Delete(client *golangsdk.ServiceClient, id string, opts DeleteOptsBuilder) (r DeleteResult) {
	url := resourceURL(client, id)
	if opts != nil {
		q, err := opts.ToShareDeleteQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += q
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = client.Delete(url, reqOpt)
	return
}
