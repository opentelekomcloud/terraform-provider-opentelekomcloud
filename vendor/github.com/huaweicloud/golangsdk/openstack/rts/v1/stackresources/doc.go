/*
Package resources enables management and retrieval of Resources
RTS service.

Example to List Resources

listOpts := stackresources.ListOpts{}
allResources, err := stackresources.List(orchestrationClient, listOpts)
if err != nil {
panic(err)
}

for _, resource := range allResources {
fmt.Printf("%+v\n", resource)
}
*/

package stackresources
