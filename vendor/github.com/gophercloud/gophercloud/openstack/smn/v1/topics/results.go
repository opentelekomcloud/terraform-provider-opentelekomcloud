package topics

import (
	"github.com/gophercloud/gophercloud"
)

type Topic struct {
	RequestId        string        `json:"request_id"`
	TopicUrn         string        `json:"topic_urn"`
}


type TopicList struct {
	TopicUrn         string         `json:"topic_urn"`
	DisplayName      string         `json:"display_name"`
	Name             string         `json:"name"`
	PushPolicy       string         `json:"push_policy"`
}


// Extract will get the Volume object out of the commonResult object.
func (r commonResult) Extract() (*Topic, error) {
	var s Topic
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

type commonResult struct {
	gophercloud.Result
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}

type DeleteResult struct {
	gophercloud.ErrResult
}

type GetResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type ListResult struct {
	gophercloud.Result
}

func (lr ListResult) Extract() ([]Topic, error) {
	var a struct {
		Topics []Topic `json:"topics"`
	}
	err := lr.Result.ExtractInto(&a)
	return a.Topics, err
}
