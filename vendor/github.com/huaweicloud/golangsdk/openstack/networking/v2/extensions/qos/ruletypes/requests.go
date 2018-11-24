package ruletypes

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// ListRuleTypes returns the list of rule types from the server
func ListRuleTypes(c *golangsdk.ServiceClient) (result pagination.Pager) {
	return pagination.NewPager(c, listRuleTypesURL(c), func(r pagination.PageResult) pagination.Page {
		return ListRuleTypesPage{pagination.SinglePageBase(r)}
	})
}
