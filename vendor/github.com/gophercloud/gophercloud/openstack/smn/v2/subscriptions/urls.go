package subscriptions

import "github.com/gophercloud/gophercloud"

func createURL(c *gophercloud.ServiceClient, topicUrn string) string {
	return c.ServiceURL("topics", topicUrn, "subscriptions" )
}

func deleteURL(c *gophercloud.ServiceClient, subscriptionUrn string) string {
	return c.ServiceURL("subscriptions", subscriptionUrn)
}

func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("subscriptions?offset=0&limit=100")
}

func listFromTopicURL(c *gophercloud.ServiceClient, topicUrn string) string {
	return c.ServiceURL("topics", topicUrn, "subscriptions?offset=0&limit=100")
}
