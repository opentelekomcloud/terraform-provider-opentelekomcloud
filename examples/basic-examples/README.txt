# The Basic Examples Of Resource Definitions

This folder contains the resource examples, currently there are more than 30 examples.
e.g. ECS, EVS, OBS ... please refer to the tf files.  

### Folders And Files

mainfile.tf: Module difinitions
provider.tf: Two provider difinitions.
var.tf: All variable difinitions.
outputs.tf: Module outputs.
public-2048.txt: the key using when creating keypair.
privite-xshell-2048.pem: The login private key.
\modules\ecs: ECS resource
\modules\ecskey:KeyPair imports
\modules\eip: EIP instance, EIP binding to instance.
\modules\lb|elb: ELB resources, adding ECS into ELB backend.
\modules\rds: RDS instance reource.
\modules\sg: Security groups and rules
\modules\vpc_subnets: Basic networking resource with v1
\modules\kms: KMS resources. 
\modules\sfs: File share resources.
\modules\ims: Image resource
\modules\evs: block volume resource.
\modules\obs: Object storage resource.
......

modules: Contains the basic examples, Currently has more than 30 resources, they are:
	main.tf: resource difinitions.
	outputs.tf: the output information difinitions.
	var.tf: the variables, please change its value with your case.
  	
