package members

import "github.com/huaweicloud/golangsdk"

func imageMembersURL(c *golangsdk.ServiceClient, imageID string) string {
	return c.ServiceURL("images", imageID, "members")
}

func listMembersURL(c *golangsdk.ServiceClient, imageID string) string {
	return imageMembersURL(c, imageID)
}

func createMemberURL(c *golangsdk.ServiceClient, imageID string) string {
	return imageMembersURL(c, imageID)
}

func imageMemberURL(c *golangsdk.ServiceClient, imageID string, memberID string) string {
	return c.ServiceURL("images", imageID, "members", memberID)
}

func getMemberURL(c *golangsdk.ServiceClient, imageID string, memberID string) string {
	return imageMemberURL(c, imageID, memberID)
}

func updateMemberURL(c *golangsdk.ServiceClient, imageID string, memberID string) string {
	return imageMemberURL(c, imageID, memberID)
}

func deleteMemberURL(c *golangsdk.ServiceClient, imageID string, memberID string) string {
	return imageMemberURL(c, imageID, memberID)
}
