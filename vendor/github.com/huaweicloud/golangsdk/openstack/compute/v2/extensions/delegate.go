package extensions

import (
	"github.com/huaweicloud/golangsdk"
	common "github.com/huaweicloud/golangsdk/openstack/common/extensions"
	"github.com/huaweicloud/golangsdk/pagination"
)

// ExtractExtensions interprets a Page as a slice of Extensions.
func ExtractExtensions(page pagination.Page) ([]common.Extension, error) {
	return common.ExtractExtensions(page)
}

// Get retrieves information for a specific extension using its alias.
func Get(c *golangsdk.ServiceClient, alias string) common.GetResult {
	return common.Get(c, alias)
}

// List returns a Pager which allows you to iterate over the full collection of extensions.
// It does not accept query parameters.
func List(c *golangsdk.ServiceClient) pagination.Pager {
	return common.List(c)
}
