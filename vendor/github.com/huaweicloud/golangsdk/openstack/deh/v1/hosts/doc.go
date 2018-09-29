package hosts

/*
Package hosts enables management and retrieval of Dedicated Hosts

Example to Allocate Hosts
	opts := hosts.AllocateOpts{Name:"c2c-test",HostType:"h1",AvailabilityZone:"eu-de-02",AutoPlacement:"off",Quantity:1}
	allocatedHosts ,err := hosts.Allocate(client,opts).Extract()
	if err != nil {
		panic(err)
	}
	fmt.Println(allocatedHosts)


Example to Update Hosts
	updateopts := hosts.UpdateOpts{Name:"NewName3",AutoPlacement:"on"}
	update := hosts.Update(client,"8ea7381e-8d84-4f9f-a7ad-d32f1e1bb5b7",updateopts)
		if err != nil {
			panic(update.Err)
		}
	fmt.Println(update)

Example to delete Hosts
	delete := hosts.Delete(client,"94d94259-3734-4ad5-bc3b-5f9f3e96d5e8")
	if err != nil {
		panic(delete.Err)
	}
	fmt.Println(delete)

Example to List Hosts
	listdeh := hosts.ListOpts{}
	alldehs, err := hosts.List(client,listdeh).AllPages()
	list,err:=hosts.ExtractHosts(alldehs)
	if err != nil {
		panic(err)
	}
	fmt.Println(list)

Example to Get Host
	result := hosts.Get(client, "66156a61-27c2-4169-936b-910dd9c73da3")
	out, err := result.Extract()
	fmt.Println(out)

Example to List Servers
	listOpts := hosts.ListServerOpts{}
	allServers, err := hosts.ListServer(client, "671611d2-b45c-4648-9e78-06eb24522291",listOpts)
	if err != nil {
	panic(err)
	}

	for _, server := range allServers {
	fmt.Printf("%+v\n", server)
	}
*/
