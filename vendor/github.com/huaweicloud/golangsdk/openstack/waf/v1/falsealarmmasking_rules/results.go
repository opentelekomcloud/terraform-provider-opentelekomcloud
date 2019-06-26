package falsealarmmasking_rules

import (
	"github.com/huaweicloud/golangsdk"
)

type ListResponse struct {
	// Total number of the Rules
	Total int `json:"total"`

	// List of AlarmMasking
	Items []AlarmMasking `json:"items"`
}

type AlarmMasking struct {
	//False Alarm Masking Rule ID
	Id string `json:"id"`
	//False Alarm Maksing Rule URL
	Url string `json:"url"`
	//Rule ID
	Rule string `json:"rule"`
	//Policy ID
	PolicyID string `json:"policyid"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a falsealarmmasking rule.
func (r commonResult) Extract() (*AlarmMasking, error) {
	var response AlarmMasking
	err := r.ExtractInto(&response)
	return &response, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a False Alarm Masking rule.
type CreateResult struct {
	commonResult
}

type ListResult struct {
	commonResult
}

func (r ListResult) Extract() ([]AlarmMasking, error) {
	var s ListResponse
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}
	return s.Items, nil
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
