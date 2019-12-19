package logtopics

import "github.com/huaweicloud/golangsdk"

// Log topic Create response
type CreateResponse struct {
	ID string `json:"log_topic_id"`
}

// CreateResult is a struct that contains all the return parameters of creation
type CreateResult struct {
	golangsdk.Result
}

// Extract from CreateResult
func (r CreateResult) Extract() (*CreateResponse, error) {
	s := new(CreateResponse)
	err := r.Result.ExtractInto(s)
	return s, err
}

// DeleteResult is a struct which contains the result of deletion
type DeleteResult struct {
	golangsdk.ErrResult
}

// Log topic response
type LogTopic struct {
	ID           string `json:"log_topic_id,omitempty"`
	Name         string `json:"log_topic_name"`
	CreationTime int64  `json:"creation_time"`
	IndexEnabled bool   `json:"index_enabled"`
}

// GetResult contains the body of getting detailed
type GetResult struct {
	golangsdk.Result
}

// Extract from GetResult
func (r GetResult) Extract() (*LogTopic, error) {
	s := new(LogTopic)
	err := r.Result.ExtractInto(s)
	return s, err
}
