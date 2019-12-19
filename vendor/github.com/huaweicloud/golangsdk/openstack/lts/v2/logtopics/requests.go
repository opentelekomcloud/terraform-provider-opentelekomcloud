package logtopics

import "github.com/huaweicloud/golangsdk"

// CreateOptsBuilder is used for creating log group parameters.
type CreateOptsBuilder interface {
	ToLogTopicsCreateMap() (map[string]interface{}, error)
}

// CreateOpts is a struct that contains all the parameters.
type CreateOpts struct {
	// Specifies the log group name.
	LogTopicName string `json:"log_topic_name" required:"true"`
}

// ToLogTopicsCreateMap is used for type convert
func (ops CreateOpts) ToLogTopicsCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(ops, "")
}

// Create a log topic with given parameters.
func Create(client *golangsdk.ServiceClient, groupId string, ops CreateOptsBuilder) (r CreateResult) {
	b, err := ops.ToLogTopicsCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client, groupId), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{201},
	})

	return
}

// Delete a log topic by id
func Delete(client *golangsdk.ServiceClient, groupId string, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, groupId, id), nil)
	return
}

// Get a log topic with detailed information by id
func Get(client *golangsdk.ServiceClient, groupId string, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, groupId, id), &r.Body, nil)
	return
}
