package availablezones

import (
	"github.com/opentelekomcloud/gophertelekomcloud"
)

// Get available zones
func Get(client *golangsdk.ServiceClient) (r GetResult) {
	_, r.Err = client.Get(getURL(client), &r.Body, nil)
	return
}
