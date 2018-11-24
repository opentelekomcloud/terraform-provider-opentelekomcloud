package suspendresume

import "github.com/huaweicloud/golangsdk"

func actionURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL("servers", id, "action")
}

// Suspend is the operation responsible for suspending a Compute server.
func Suspend(client *golangsdk.ServiceClient, id string) (r SuspendResult) {
	_, r.Err = client.Post(actionURL(client, id), map[string]interface{}{"suspend": nil}, nil, nil)
	return
}

// Resume is the operation responsible for resuming a Compute server.
func Resume(client *golangsdk.ServiceClient, id string) (r UnsuspendResult) {
	_, r.Err = client.Post(actionURL(client, id), map[string]interface{}{"resume": nil}, nil, nil)
	return
}
