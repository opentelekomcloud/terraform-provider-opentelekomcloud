package tasks

import (
	"net/url"
	"strings"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/utils"
)

const resourcePath = "tasks"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, taskID string) string {
	return c.ServiceURL(resourcePath, taskID)
}

func listURL(c *golangsdk.ServiceClient) string {
	return rootURL(c)
}

func getURL(c *golangsdk.ServiceClient, taskID string) string {
	return resourceURL(c, taskID)
}

func createURL(c *golangsdk.ServiceClient) string {
	return rootURL(c)
}

func nextPageURL(serviceURL, requestedNext string) (string, error) {
	base, err := utils.BaseEndpoint(serviceURL)
	if err != nil {
		return "", err
	}

	requestedNextURL, err := url.Parse(requestedNext)
	if err != nil {
		return "", err
	}

	base = golangsdk.NormalizeURL(base)
	nextPath := base + strings.TrimPrefix(requestedNextURL.Path, "/")

	nextURL, err := url.Parse(nextPath)
	if err != nil {
		return "", err
	}

	nextURL.RawQuery = requestedNextURL.RawQuery

	return nextURL.String(), nil
}
