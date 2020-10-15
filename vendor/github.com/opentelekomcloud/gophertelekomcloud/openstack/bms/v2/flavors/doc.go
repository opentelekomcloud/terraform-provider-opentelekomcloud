/*
Package flavors enables management and retrieval of Flavors
BMS service.

Example to List Flavors

	listOpts := flavors.ListOpts{}
	allFlavors, err := flavors.List(bmsClient, listOpts)
	if err != nil {
		panic(err)
	}

	for _, flavor := range allFlavors {
		fmt.Printf("%+v\n", flavor)
	}
*/
package flavors
