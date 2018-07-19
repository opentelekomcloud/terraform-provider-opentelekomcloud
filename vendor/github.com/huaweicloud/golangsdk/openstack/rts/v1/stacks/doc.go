// Package stacks provides operation for working with Heat stacks. A stack is a
// group of resources (servers, load balancers, databases, and so forth)
// combined to fulfill a useful purpose. Based on a template, Heat orchestration
// engine creates an instantiated set of resources (a stack) to run the
// application framework or component specified (in the template). A stack is a
// running instance of a template. The result of creating a stack is a deployment
// of the application framework or component.

/*
Package resources enables management and retrieval of
RTS service.

Example to List Stacks

lis:=stacks.ListOpts{SortDir:stacks.SortAsc,SortKey:stacks.SortStatus}
	getstack,err:=stacks.List(client,lis).AllPages()
	stacks,err:=stacks.ExtractStacks(getstack)
	fmt.Println(stacks)

Example to Create a Stacks
template := new(stacks.Template)
	template.Bin = []byte(`
			  {
			 "heat_template_version": "2013-05-23",
			 "description": "Simple template to deploy",
			 "parameters": {
			  "image_id": {
			   "type": "string",
			   "description": "Image to be used for compute instance",
			   "label": "Image ID",
			   "default": "ea67839e-fd7a-4b99-9f81-13c4c8dc317c"
			  },
			  "net_id": {
			   "type": "string",
			   "description": "The network to be used",
			   "label": "Network UUID",
			   "default": "7eb54ab6-5cdb-446a-abbe-0dda1885c76e"
			  },
			  "instance_type": {
			   "type": "string",
			   "description": "Type of instance (flavor) to be used",
			   "label": "Instance Type",
			   "default": "s1.medium"
			  }
			 },
			 "resources": {
			  "my_instance": {
			   "type": "OS::Nova::Server",
			   "properties": {
			    "image": {
			     "get_param": "image_id"
			    },
			    "flavor": {
			     "get_param": "instance_type"
			    },
			    "networks": [
			     {
			      "network": {
			       "get_param": "net_id"
			      }
			     }
			    ]
			   }
			  }
			 }
			}`)

	fmt.Println(template)
	stack:=stacks.CreateOpts{Name:"terraform-providerr_disssssss",Timeout: 60,
		TemplateOpts:    template, }
	outstack,err:=stacks.Create(client,stack).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a Stacks

template := new(stacks.Template)
	template.Bin = []byte(`
			  {
			 "heat_template_version": "2013-05-23",
			 "description": "Simple template disha",
			 "parameters": {
			  "image_id": {
			   "type": "string",
			   "description": "Image to be used for compute instance",
			   "label": "Image ID",
			   "default": "ea67839e-fd7a-4b99-9f81-13c4c8dc317c"
			  },
			  "net_id": {
			   "type": "string",
			   "description": "The network to be used",
			   "label": "Network UUID",
			   "default": "7eb54ab6-5cdb-446a-abbe-0dda1885c76e"
			  },
			  "instance_type": {
			   "type": "string",
			   "description": "Type of instance (flavor) to be used",
			   "label": "Instance Type",
			   "default": "s1.medium"
			  }
			 },
			 "resources": {
			  "my_instance": {
			   "type": "OS::Nova::Server",
			   "properties": {
			    "image": {
			     "get_param": "image_id"
			    },
			    "flavor": {
			     "get_param": "instance_type"
			    },
			    "networks": [
			     {
			      "network": {
			       "get_param": "net_id"
			      }
			     }
			    ]
			   }
			  }
			 }
			}`)
	myopt:=stacks.UpdateOpts{TemplateOpts:template}
	errup:=stacks.Update(client,"terraform-providerr_stack_bigdata","b631d6ce-9010-4c34-b994-cd5a10b13c86",myopt).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Delete a Stacks

	err:=stacks.Delete(client,"terraform-providerr_stack_bigdata","b631d6ce-9010-4c34-b994-cd5a10b13c86")

	if err != nil {
		panic(err)
		}
*/
package stacks
