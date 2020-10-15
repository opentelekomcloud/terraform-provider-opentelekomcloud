/*
Package tags enables management and retrieval of Tags
BMS service.

Example to Create a Tag

	createOpts := tags.CreateOpts{
		Tag: []string{"__type_baremetal"},

	}

	tag, err := tags.Create(bmsClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Tag

	serverID := "4e8e5957-649f-477b-9e5b-f1f75b21c03c"
	err := tags.Delete(bmsClient, serverID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package tags
