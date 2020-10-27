/*
Package softwaredeployment enables management and retrieval of Software Deployments

Example to List Software Deployments

	listOpts := softwaredeployment.ListOpts{}
	allDeployments, err := softwaredeployment.List(client,listOpts)
	if err != nil {
		panic(err)
	}

	for _, deployment := range allDeployments {
		fmt.Printf("%+v\n", allDeployments)
	}

Example to Get Software Deployment

	deploymentID:="bd7d48a5-6e33-4b95-aa28-d0d3af46c635"

 	deployments,err:=softwaredeployment.Get(client,deploymentID).Extract()

	if err != nil {
		panic(err)
	}

Example to Create a Software Deployments

	input:=map[string]interface{}{"name":"foo"}

	createOpts := softwaredeployment.CreateOpts{
		Status:"IN_PROGRESS",
		ServerId:"f274ac7d-334d-41ff-83bd-1de669f7781b",
		ConfigId:"a6ff3598-f2e0-4111-81b0-aa3e1cac2529",
		InputValues:input,
		TenantId:"17fbda95add24720a4038ba4b1c705ed",
		Action:"CREATE"
	}

	deployment, err := softwaredeployment.Create(client, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a Software Deployments
	deploymentID:="bd7d48a5-6e33-4b95-aa28-d0d3af46c635"

	ouput:=map[string]interface{}{"deploy_stdout":"Writing to /tmp/baaaaa\nWritten to /tmp/baaaaa\n","deploy_stderr":"+ echo Writing to /tmp/baaaaa\n+ echo fooooo\n+ cat /tmp/baaaaa\n+ echo -n The file /tmp/baaaaa contains fooooo for server ec14c864-096e-4e27-bb8a-2c2b4dc6f3f5 during CREATE\n+ echo Written to /tmp/baaaaa\n+ echo Output to stderr\nOutput to stderr\n",
		"deploy_status_code":0,"result":"The file /tmp/baaaaa contains fooooo for server ec14c864-096e-4e27-bb8a-2c2b4dc6f3f5 during CREATE"}

	updateOpts := softwaredeployment.UpdateOpts{
	Status:"COMPLETE",
	ConfigId:"a6ff3598-f2e0-4111-81b0-aa3e1cac2529",
	OutputValues:ouput,
	StatusReason:"Outputs received"}

	deployment, err := softwaredeployment.Update(client, deploymentID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}


Example to Delete a Software Deployments

	deploymentID:="bd7d48a5-6e33-4b95-aa28-d0d3af46c635"
	del:=softwaredeployment.Delete(client,deploymentID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package softwaredeployment
