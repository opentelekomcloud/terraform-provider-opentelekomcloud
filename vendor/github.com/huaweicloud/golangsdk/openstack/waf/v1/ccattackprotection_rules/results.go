package ccattackprotection_rules

import (
	"github.com/huaweicloud/golangsdk"
)

type CcAttack struct {
	Id           string       `json:"id"`
	PolicyID     string       `json:"policyid"`
	Url          string       `json:"url"`
	LimitNum     int          `json:"limit_num"`
	LimitPeriod  int          `json:"limit_period"`
	LockTime     int          `json:"lock_time"`
	TagType      string       `json:"tag_type"`
	TagIndex     string       `json:"tag_index"`
	TagCondition TagCondition `json:"tag_condition"`
	Action       Action       `json:"action"`
	Default      bool         `json:"default"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a cc attack protection rule.
func (r commonResult) Extract() (*CcAttack, error) {
	var response CcAttack
	err := r.ExtractInto(&response)
	return &response, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a cc attack protection rule.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a cc attack protection rule.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
