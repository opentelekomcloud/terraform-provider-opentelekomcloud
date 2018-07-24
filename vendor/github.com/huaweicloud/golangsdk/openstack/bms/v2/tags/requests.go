package tags

import (
	"github.com/huaweicloud/golangsdk"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// create request.
type CreateOptsBuilder interface {
	ToTagsCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new tag.
type CreateOpts struct {
	Tag []string `json:"tags" required:"true"`
}

// ToTagsCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToTagsCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create will create a new Tag based on the values in CreateOpts. To extract
// the Tag object from the response, call the Extract method on the
// CreateResult.
func Create(c *golangsdk.ServiceClient, serverId string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToTagsCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, serverId), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get retrieves a particular tag based on its unique ID.
func Get(c *golangsdk.ServiceClient, serverId string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, serverId), &r.Body, nil)
	return
}

// Delete will permanently delete a particular tag based on its unique ID.
func Delete(c *golangsdk.ServiceClient, serverId string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, serverId), nil)
	return
}
