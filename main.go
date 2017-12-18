package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/Click2Cloud/terraform-provider-opentelekomcloud/opentelekomcloud" // TODO: Revert path when merge
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: opentelekomcloud.Provider})
}

/*
package main

import (

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"fmt"
	"github.com/gophercloud/gophercloud/openstack/networking/v1/extensions/vpcs"

)
func main() {

	opts1 := gophercloud.AuthOptions{

		IdentityEndpoint: "https://iam.eu-de.otc.t-systems.com/v3",
		Username:         "lizhonghua",
		Password:         "slob@123",
		DomainName:	      "OTC00000000001000010501",
		//TenantID:         "87a56a48977e42068f70ad3280c50f0e",
		TenantName:       "eu-de_Nordea",
		//TokenID:		  "MIIDmQYJKoZIhvcNAQcCoIIDijCCA4YCAQExDTALBglghkgBZQMEAgEwggFnBgkqhkiG9w0BBwGgggFYBIIBVHsidG9rZW4iOnsiaXNzdWVkX2F0IjoiMjAxNy0xMS0wNlQwNDoyMzo0Ni45ODQwMDBaIiwiZXhwaXJlc19hdCI6IjIwMTctMTEtMDdUMDQ6MjM6NDYuOTg0MDAwWiIsIm1ldGhvZHMiOlsicGFzc3dvcmQiXSwidXNlciI6eyJkb21haW4iOnsibmFtZSI6Ik9UQzAwMDAwMDAwMDAxMDAwMDEwNTA3IiwiaWQiOiI3MTYyOTA3MzgxOGI0NDVkOGQzMDU5MWM1NDUwNTZkMCIsInhkb21haW5fdHlwZSI6IlRTSSIsInhkb21haW5faWQiOiIwMDAwMDAwMDAwMTAwMDAxMDUwNyJ9LCJpZCI6IjczY2YxOTVkMjFjNzQzNjM5NzVkMTllMmIwMGNjMDk1IiwibmFtZSI6ImxpeW9uZ2xlIn0sImNhdGFsb2ciOltdfX0xggIFMIICAQIBATBcMFcxCzAJBgNVBAYTAlVTMQ4wDAYDVQQIDAVVbnNldDEOMAwGA1UEBwwFVW5zZXQxDjAMBgNVBAoMBVVuc2V0MRgwFgYDVQQDDA93d3cuZXhhbXBsZS5jb20CAQEwCwYJYIZIAWUDBAIBMA0GCSqGSIb3DQEBAQUABIIBgJeSqcO-SQuw1zyoxkFMFwG71JYrQk439VMKSgC+CbGR3uYN6AauNimQKSgE42zWA+t0uGBXea4K5wS183iEumCElBDnr4qtsTv0zNBPA81lfzKxzRzvIwxp9qrOJO4xKY4B4avo4UXHuyCTxqaIr8qc1gfnBudWnCQkRjCJccQazs9iD4NaNUkwS2GiVLm2vQs-xVdjcSoNx1WRMdL2pHk-VpsM-NLC8WPioVOZrTJlbt53tUb5EOFwysxbsAYiaIjthqWuAlziRyPE5by-hibchucmoggwec1T99IZjAaBfMu0XkKobDnfGiHv1np+auTABKwQK+J8+B2PsqtqbQW1MQxUTeooF4YuK-6Ekv6Q14ybogJEhq6yQOwV99Gbv+aGlxDwRt1zFYwzFpLC824qeAdp68fvHmL9Zlu0Ja0bsViuUvQzoRmquETmaAYu+9O9aWCAJpAR8169BxKh5GU-7O5yz0l4t9eFPcy4imCDqyHDGV7pmuNaIOC6MxRqww==",

	}


	*/
/*opts1 := gophercloud.AuthOptions{

		IdentityEndpoint: "https://iam.eu-de.otc.t-systems.com/v1",
		Username:         "liyongle",
		Password:         "OpenStackSDK",
		DomainName:	      "OTC00000000001000010507",
		TenantID:         "18899b93e7be46c2b7f0d2568efabc33",
		//TokenID:		  "MIIDmQYJKoZIhvcNAQcCoIIDijCCA4YCAQExDTALBglghkgBZQMEAgEwggFnBgkqhkiG9w0BBwGgggFYBIIBVHsidG9rZW4iOnsiaXNzdWVkX2F0IjoiMjAxNy0xMS0wNlQwNDoyMzo0Ni45ODQwMDBaIiwiZXhwaXJlc19hdCI6IjIwMTctMTEtMDdUMDQ6MjM6NDYuOTg0MDAwWiIsIm1ldGhvZHMiOlsicGFzc3dvcmQiXSwidXNlciI6eyJkb21haW4iOnsibmFtZSI6Ik9UQzAwMDAwMDAwMDAxMDAwMDEwNTA3IiwiaWQiOiI3MTYyOTA3MzgxOGI0NDVkOGQzMDU5MWM1NDUwNTZkMCIsInhkb21haW5fdHlwZSI6IlRTSSIsInhkb21haW5faWQiOiIwMDAwMDAwMDAwMTAwMDAxMDUwNyJ9LCJpZCI6IjczY2YxOTVkMjFjNzQzNjM5NzVkMTllMmIwMGNjMDk1IiwibmFtZSI6ImxpeW9uZ2xlIn0sImNhdGFsb2ciOltdfX0xggIFMIICAQIBATBcMFcxCzAJBgNVBAYTAlVTMQ4wDAYDVQQIDAVVbnNldDEOMAwGA1UEBwwFVW5zZXQxDjAMBgNVBAoMBVVuc2V0MRgwFgYDVQQDDA93d3cuZXhhbXBsZS5jb20CAQEwCwYJYIZIAWUDBAIBMA0GCSqGSIb3DQEBAQUABIIBgJeSqcO-SQuw1zyoxkFMFwG71JYrQk439VMKSgC+CbGR3uYN6AauNimQKSgE42zWA+t0uGBXea4K5wS183iEumCElBDnr4qtsTv0zNBPA81lfzKxzRzvIwxp9qrOJO4xKY4B4avo4UXHuyCTxqaIr8qc1gfnBudWnCQkRjCJccQazs9iD4NaNUkwS2GiVLm2vQs-xVdjcSoNx1WRMdL2pHk-VpsM-NLC8WPioVOZrTJlbt53tUb5EOFwysxbsAYiaIjthqWuAlziRyPE5by-hibchucmoggwec1T99IZjAaBfMu0XkKobDnfGiHv1np+auTABKwQK+J8+B2PsqtqbQW1MQxUTeooF4YuK-6Ekv6Q14ybogJEhq6yQOwV99Gbv+aGlxDwRt1zFYwzFpLC824qeAdp68fvHmL9Zlu0Ja0bsViuUvQzoRmquETmaAYu+9O9aWCAJpAR8169BxKh5GU-7O5yz0l4t9eFPcy4imCDqyHDGV7pmuNaIOC6MxRqww==",

	}*//*

	provider, err := openstack.AuthenticatedClient(opts1)
	if err != nil {
		fmt.Println(err)
	}

	endpoint:=gophercloud.EndpointOpts{
		///Name:   "neutron",
		Region: "eu-de",
	}

	client,err:=openstack.NewVpcV1(provider,endpoint)

	//Creating vpc
	vpc:=vpcs.CreateOpts{Name:"terraform-provider-mytestVPC12"}
	outvpc,err:=vpcs.Create(client,vpc).Extract()
	fmt.Println(outvpc)

	//Querying VPC Details
	*/
/*result:=vpcs.Get(client,"abda1f6e-ae7c-4ff5-8d06-53425dc11f34")
	out,err:=result.Extract()
	fmt.Println(out)*//*


	//Querying VPCs
	*/
/*	listvpc:=vpcs.ListOpts{}
		out,err:=vpcs.List(client,listvpc).AllPages()
		result,err:=vpcs.ExtractVpcs(out)
		fmt.Println(result[0])*//*



	//Updating VPC Information
	*/
/*updatevpc:=vpcs.UpdateOpts{Name:"myvpc", CIDR:"192.168.0.0/16"}
	out,err:=vpcs.Update(client,"4f0d165f-5c89-4b1e-8c1c-530496546833",updatevpc).Extract()
	fmt.Println(out)
	fmt.Println(err)
*//*

	//Deleting a VPC
	*/
/*out:=vpcs.Delete(client,"b808fa6d-fd36-416e-bc68-91ef209462aa")
	fmt.Println(out)*//*


	*/
/*ro:=routers.CreateOpts{Name:"testdish"}
	list,err:= routers.Create(client,ro).Extract()
	fmt.Println(list)*//*













}
*/
