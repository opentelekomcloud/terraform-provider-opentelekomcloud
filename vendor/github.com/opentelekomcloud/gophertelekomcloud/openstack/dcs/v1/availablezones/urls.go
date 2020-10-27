package availablezones

import (
	"strings"

	"github.com/opentelekomcloud/gophertelekomcloud"
)

// endpoint/availablezones
const resourcePath = "availableZones"

// getURL will build the get url of get function
func getURL(client *golangsdk.ServiceClient) string {
	// remove projectid from endpoint
	return strings.Replace(client.ServiceURL(resourcePath), "/"+client.ProjectID, "", -1)
}
