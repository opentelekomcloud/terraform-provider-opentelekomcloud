package keypairs

import (
	"github.com/opentelekomcloud/gophertelekomcloud"
)

const resourcePath = "os-keypairs"

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}
