/*
Package softwareconfig enables management and retrieval of Software Configs

Example to List Software Configs

	listOpts := softwareconfig.ListOpts{}
	allConfigs, err := softwareconfig.List(client,listOpts)
	if err != nil {
		panic(err)
	}

	for _, config := range allConfigs {
		fmt.Printf("%+v\n", config)
	}

Example to Create a Software Configs

	createOpts := vsoftwareconfig.CreateOpts{
		Name:         "config_test",
	}

	config, err := softwareconfig.Create(client, createOpts).Extract()
	if err != nil {
		panic(err)
	}


Example to Delete a Software Configs

	configID := "8de48948-b6d6-4417-82a5-071f7811af91"
	del:=softwareconfig.Delete(client,configID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package softwareconfig
