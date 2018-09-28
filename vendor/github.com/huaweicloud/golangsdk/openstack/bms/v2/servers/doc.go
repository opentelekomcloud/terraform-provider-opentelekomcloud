/*
Package servers enables management and retrieval of Servers
BMS service.

Example to List Servers

	listOpts := servers.ListOpts{}
	allServers, err := servers.List(bmsClient, listOpts)
	if err != nil {
		panic(err)
	}

	for _, server := range allServers {
		fmt.Printf("%+v\n", server)
	}
*/
package servers
