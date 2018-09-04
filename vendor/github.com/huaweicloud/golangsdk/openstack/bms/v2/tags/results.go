package tags

import (
	"github.com/huaweicloud/golangsdk"
)

type Tags struct {
	//Specifies the tags of a BMS
	Tags []string `json:"tags"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract interprets any commonResult as a Tags.
func (r commonResult) Extract() (*Tags, error) {
	var s *Tags
	err := r.ExtractInto(&s)
	return s, err
}

// CreateResult represents the result of an create operation. Call its Extract
// method to interpret it as a Tag.
type CreateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Tag.
type GetResult struct {
	commonResult
}
