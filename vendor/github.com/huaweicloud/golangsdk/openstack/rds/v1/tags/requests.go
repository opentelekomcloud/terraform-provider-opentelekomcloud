package tags

import (
	"github.com/huaweicloud/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Tag is a structure of key value pair.
type CreateOpts struct {
	//tag key
	Key string `json:"key" required:"true"`
	//tag value
	Value string `json:"value" required:"true"`
}

type DeleteOpts struct {
	//tag key
	Key string `json:"key" required:"true"`
}

// ToCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "tag")
}

// ToDeleteMap assembles a request body based on the contents of a
// DeleteOpts.
func (opts DeleteOpts) ToDeleteMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

type CreateOptsBuilder interface {
	ToCreateMap() (map[string]interface{}, error)
}

type DeleteOptsBuilder interface {
	ToDeleteMap() (map[string]interface{}, error)
}

// Create implements tag create request.
func Create(client *golangsdk.ServiceClient, id string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(resourceURL(client, id), b, nil, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}

// Delete implements tag delete request.
func Delete(client *golangsdk.ServiceClient, id string, opts DeleteOptsBuilder) (r DeleteResult) {
	b, err := opts.ToDeleteMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.DeleteWithBody(resourceURL(client, id), b, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}

// Get implements tag get request.
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}
