/*
Package keypairs provides the ability to manage key pairs as well as create
bms with a specified key pair.

Example to List Key Pairs

	listkeypair := keypairs.ListOpts{Name: "c2c-keypair1"}
	allkeypair, err := keypairs.List(client,listkeypair)
	if err != nil {
		panic(err)
	}

	fmt.Println(allkeypair)


*/
package keypairs
