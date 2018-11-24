package imageimport

import "github.com/huaweicloud/golangsdk"

const (
	infoPath     = "info"
	resourcePath = "import"
)

func infoURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(infoPath, resourcePath)
}
