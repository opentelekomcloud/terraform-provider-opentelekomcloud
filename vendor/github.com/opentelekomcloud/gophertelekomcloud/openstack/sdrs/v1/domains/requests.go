package domains

import (
	"github.com/opentelekomcloud/gophertelekomcloud"
)

// Get active-active domains
func Get(client *golangsdk.ServiceClient) (r GetResult) {
	_, r.Err = client.Get(getURL(client), &r.Body, nil)
	return
}
