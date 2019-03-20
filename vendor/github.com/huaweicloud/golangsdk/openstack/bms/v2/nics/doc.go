/*
Package nics enables management and retrieval of NICs
BMS service.

Example to List Flavors

	listOpts := nics.ListOpts{}
	allNics, err := nics.List(bmsClient, listOpts)
	if err != nil {
		panic(err)
	}

	for _, nic := range allNics {
		fmt.Printf("%+v\n", nic)
	}
*/
package nics
