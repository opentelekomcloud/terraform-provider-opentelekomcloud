package loggroups

import "github.com/huaweicloud/golangsdk"

// CreateOptsBuilder is used for creating log group parameters.
type CreateOptsBuilder interface {
	ToLogGroupsCreateMap() (map[string]interface{}, error)
}

// CreateOpts is a struct that contains all the parameters.
type CreateOpts struct {
	// Specifies the log group name.
	LogGroupName string `json:"log_group_name" required:"true"`

	// Specifies the log expiration time. The value is fixed to 7 days.
	TTL int `json:"ttl_in_daysion,omitempty"`
}

// ToLogGroupsCreateMap is used for type convert
func (ops CreateOpts) ToLogGroupsCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(ops, "")
}

// Create a log group with given parameters.
func Create(client *golangsdk.ServiceClient, ops CreateOptsBuilder) (r CreateResult) {
	b, err := ops.ToLogGroupsCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{201},
	})

	return
}

// Delete a log group by id
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}

// Get a log group with detailed information by id
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}
