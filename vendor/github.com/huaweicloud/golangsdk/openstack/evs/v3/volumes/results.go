package volumes

import (
	"github.com/huaweicloud/golangsdk"
)

// Volume contains all the information associated with a Volume.
type Volume struct {
	// Unique identifier for the volume.
	ID string `json:"id"`
	// wwn of the volume.
	WWN string `json:"wwn"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract will get the Volume object out of the commonResult object.
func (r commonResult) Extract() (*Volume, error) {
	var s Volume
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractInto converts our response data into a volume struct
func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "volume")
}

// GetResult contains the response body and error from a Get request.
type GetResult struct {
	commonResult
}
