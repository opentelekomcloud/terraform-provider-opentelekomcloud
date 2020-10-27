package tags

import (
	"github.com/opentelekomcloud/gophertelekomcloud"
)

func resourceURL(client *golangsdk.ServiceClient, serverId string) string {
	return client.ServiceURL("servers", serverId, "tags")
}
