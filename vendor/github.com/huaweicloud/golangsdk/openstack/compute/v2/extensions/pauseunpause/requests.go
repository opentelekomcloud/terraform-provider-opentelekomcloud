package pauseunpause

import "github.com/huaweicloud/golangsdk"

func actionURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL("servers", id, "action")
}

// Pause is the operation responsible for pausing a Compute server.
func Pause(client *golangsdk.ServiceClient, id string) (r PauseResult) {
	_, r.Err = client.Post(actionURL(client, id), map[string]interface{}{"pause": nil}, nil, nil)
	return
}

// Unpause is the operation responsible for unpausing a Compute server.
func Unpause(client *golangsdk.ServiceClient, id string) (r UnpauseResult) {
	_, r.Err = client.Post(actionURL(client, id), map[string]interface{}{"unpause": nil}, nil, nil)
	return
}
