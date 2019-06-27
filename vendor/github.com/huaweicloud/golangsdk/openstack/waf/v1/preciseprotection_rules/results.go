package preciseprotection_rules

import (
	"github.com/huaweicloud/golangsdk"
)

type Precise struct {
	Id         string      `json:"id"`
	PolicyID   string      `json:"policyid"`
	Name       string      `json:"name"`
	Time       bool        `json:"time"`
	Start      int64       `json:"start"`
	End        int64       `json:"end"`
	Conditions []Condition `json:"conditions"`
	Action     Action      `json:"action"`
	Priority   int         `json:"priority"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a precise protection rule.
func (r commonResult) Extract() (*Precise, error) {
	var response Precise
	err := r.ExtractInto(&response)
	return &response, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a precise protection rule.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a precise protection rule.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
