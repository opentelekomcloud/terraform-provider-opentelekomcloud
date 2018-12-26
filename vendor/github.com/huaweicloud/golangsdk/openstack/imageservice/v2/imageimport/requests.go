package imageimport

import "github.com/huaweicloud/golangsdk"

// ImportMethod represents valid Import API method.
type ImportMethod string

const (
	// GlanceDirectMethod represents glance-direct Import API method.
	GlanceDirectMethod ImportMethod = "glance-direct"

	// WebDownloadMethod represents web-download Import API method.
	WebDownloadMethod ImportMethod = "web-download"
)

// Get retrieves Import API information data.
func Get(c *golangsdk.ServiceClient) (r GetResult) {
	_, r.Err = c.Get(infoURL(c), &r.Body, nil)
	return
}
