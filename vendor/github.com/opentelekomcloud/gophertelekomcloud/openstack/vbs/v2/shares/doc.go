/*
Package shares enables management and retrieval of Shares
VBS service.

Example to List Shares

	listOpts := shares.ListOpts{}
	allShares, err := shares.List(vbsClient, listOpts)
	if err != nil {
		panic(err)
	}

	for _, share := range allShares {
		fmt.Printf("%+v\n", share)
	}

Example to Get a Share

   getshare,err:=shares.Get(vbsClient, "6149e448-dcac-4691-96d9-041e09ef617f").ExtractShare()
   if err != nil {
         panic(err)
		}

   fmt.Println(getshare)

Example to Create a Share

	createOpts := shares.CreateOpts{BackupID:"87566ed6-72cb-4053-aa6e-6f6216b3d507",
									ToProjectIDs:[]string{"91d687759aed45d28b5f6084bc2fa8ad"}}

	share, err := shares.Create(vbsClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Share

	shareID := "4e8e5957-649f-477b-9e5b-f1f75b21c03c"

	deleteopts := shares.DeleteOpts{IsBackupID:false}
	err := shares.Delete(vbsclient,shareID,deleteopts).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package shares
