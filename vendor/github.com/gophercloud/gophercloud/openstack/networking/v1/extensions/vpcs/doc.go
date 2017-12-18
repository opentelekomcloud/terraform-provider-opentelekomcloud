/*
Package vpcs enables management and retrieval of Vpcs from the OpenStack
Networking service.

Example to List Vpcs

	listOpts := vpcs.ListOpts{}
	allPages, err := vpcs.List(networkClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allVpcs, err := vpcs.ExtractVpcs(allPages)
	if err != nil {
		panic(err)
	}

	for _, vpc := range allRoutes {
		fmt.Printf("%+v\n", vpc)
	}

Example to Create a Vpc

	iTrue := true
	gwi := vpcs.GatewayInfo{
		NetworkID: "8ca37218-28ff-41cb-9b10-039601ea7e6b",
	}

	createOpts := vpcs.CreateOpts{
		Name:         "vpc_1",
		AdminStateUp: &iTrue,
		GatewayInfo:  &gwi,
	}

	vpc, err := vpcs.Create(networkClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a Vpc

	vpcID := "4e8e5957-649f-477b-9e5b-f1f75b21c03c"

	routes := []vpcs.Route{{
		DestinationCIDR: "40.0.1.0/24",
		NextHop:         "10.1.0.10",
	}}

	updateOpts := vpcs.UpdateOpts{
		Name:   "new_name",
		Routes: routes,
	}

	vpc, err := vpcs.Update(networkClient, vpcID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Remove all Routes from a Vpc

	vpcID := "4e8e5957-649f-477b-9e5b-f1f75b21c03c"

	routes := []vpcs.Route{}

	updateOpts := vpcs.UpdateOpts{
		Routes: routes,
	}

	vpc, err := vpcs.Update(networkClient, vpcID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Vpc

	vpcID := "4e8e5957-649f-477b-9e5b-f1f75b21c03c"
	err := vpcs.Delete(networkClient, vpcID).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Add an Interface to a Vpc

	vpcID := "4e8e5957-649f-477b-9e5b-f1f75b21c03c"

	intOpts := vpcs.AddInterfaceOpts{
		SubnetID: "a2f1f29d-571b-4533-907f-5803ab96ead1",
	}

	interface, err := vpcs.AddInterface(networkClient, vpcID, intOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Remove an Interface from a Vpc

	vpcID := "4e8e5957-649f-477b-9e5b-f1f75b21c03c"

	intOpts := vpcs.RemoveInterfaceOpts{
		SubnetID: "a2f1f29d-571b-4533-907f-5803ab96ead1",
	}

	interface, err := vpcs.RemoveInterface(networkClient, vpcID, intOpts).Extract()
	if err != nil {
		panic(err)
	}
*/
package vpcs
