package shares

import (
	"encoding/json"
	"time"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type Share struct {
	//Details about the source backup
	Backup Backup `json:"backup"`
	//Backup ID
	BackupID string `json:"backup_id"`
	//Backup share ID
	ID string `json:"id"`
	//ID of the project with which the backup is shared
	ToProjectID string `json:"to_project_id"`
	//ID of the project that shares the backup
	FromProjectID string `json:"from_project_id"`
	//Creation time of the backup share
	CreatedAt time.Time `json:"-"`
	//Update time of the backup share
	UpdatedAt time.Time `json:"-"`
	//Whether the backup has been deleted
	Deleted string `json:"deleted"`
	//Deletion time
	DeletedAt time.Time `json:"-"`
}

type Backup struct {
	//Backup ID
	ID string `json:"id"`
	//Backup name
	Name string `json:"name"`
	//Backup status
	Status string `json:"status"`
	//Backup description
	Description string `json:"description"`
	//AZ where the backup resides
	AvailabilityZone string `json:"availability_zone"`
	//Source volume ID of the backup
	VolumeID string `json:"volume_id"`
	//Cause of the backup failure
	FailReason string `json:"fail_reason"`
	//Backup size
	Size int `json:"size"`
	//Number of objects on OBS for the disk data
	ObjectCount int `json:"object_count"`
	//Container of the backup
	Container string `json:"container"`
	//Backup creation time
	CreatedAt time.Time `json:"-"`
	//Backup metadata
	ServiceMetadata string `json:"service_metadata"`
	//Time when the backup was updated
	UpdatedAt time.Time `json:"-"`
	//Current time
	DataTimeStamp time.Time `json:"-"`
	//Whether a dependent backup exists
	DependentBackups bool `json:"has_dependent_backups"`
	//ID of the snapshot associated with the backup
	SnapshotID string `json:"snapshot_id"`
	//Whether the backup is an incremental backup
	IsIncremental bool `json:"is_incremental"`
}

type commonResult struct {
	golangsdk.Result
}

// SharePage is the page returned by a pager when traversing over a
// collection of Shares.
type SharePage struct {
	pagination.LinkedPageBase
}

// Extract is a function that accepts a result and extracts shares.
func (r commonResult) Extract() ([]Share, error) {
	var s struct {
		Share []Share `json:"shared"`
	}
	err := r.ExtractInto(&s)
	return s.Share, err
}

// ExtractShare is a function that accepts a result and extracts a share.
func (r commonResult) ExtractShare() (*Share, error) {
	var s struct {
		Share *Share `json:"shared"`
	}
	err := r.ExtractInto(&s)
	return s.Share, err
}

// NextPageURL is invoked when a paginated collection of Shares has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r SharePage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"shared_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a SharePage struct is empty.
func (r SharePage) IsEmpty() (bool, error) {
	is, err := ExtractShareList(r)
	return len(is) == 0, err
}

// ExtractShareList accepts a Page struct, specifically a SharePage struct,
// and extracts the elements into a slice of Shares struct. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractShareList(r pagination.Page) ([]Share, error) {
	var s struct {
		Shares []Share `json:"shared"`
	}
	err := (r.(SharePage)).ExtractInto(&s)
	return s.Shares, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Share.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its ExtractShare
// method to interpret it as a Share.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}

// UnmarshalJSON overrides the default, to convert the JSON API response into our Backup struct
func (r *Backup) UnmarshalJSON(b []byte) error {
	type tmp Backup
	var s struct {
		tmp
		CreatedAt     golangsdk.JSONRFC3339MilliNoZ `json:"created_at"`
		UpdatedAt     golangsdk.JSONRFC3339MilliNoZ `json:"updated_at"`
		DataTimeStamp golangsdk.JSONRFC3339MilliNoZ `json:"data_timestamp"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*r = Backup(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)
	r.DataTimeStamp = time.Time(s.DataTimeStamp)

	return nil
}

// UnmarshalJSON overrides the default, to convert the JSON API response into our Share struct
func (r *Share) UnmarshalJSON(b []byte) error {
	type tmp Share
	var s struct {
		tmp
		CreatedAt golangsdk.JSONRFC3339MilliNoZ `json:"created_at"`
		UpdatedAt golangsdk.JSONRFC3339MilliNoZ `json:"updated_at"`
		DeletedAt golangsdk.JSONRFC3339MilliNoZ `json:"deleted_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*r = Share(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)
	r.DeletedAt = time.Time(s.DeletedAt)

	return nil
}
