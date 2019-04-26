package tags

import (
	"github.com/huaweicloud/golangsdk"
)

type Tag struct {
	//tag key
	Key string `json:"key"`
	//tag value
	Value string `json:"value"`
}

type RespTags struct {
	//contains list of tags, i.e.key value pair
	Tags []Tag `json:"tags"`
}

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	golangsdk.ErrResult
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type GetResult struct {
	commonResult
}

func (r commonResult) Extract() (*RespTags, error) {
	var response RespTags
	err := r.ExtractInto(&response)
	return &response, err
}
