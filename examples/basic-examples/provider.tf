provider "opentelekomcloud" {
    #please follow the documentation to change options
    #user_name  = "${var.user_name}"
    tenant_name = "${var.tenant_name}"
    domain_name = "h00454348" 
    #password  = "${var.password}"	
    access_key  = "your ak"
    secret_key = "your sk"        
    insecure = "true"
    region = "eu-de"
    auth_url  = "${var.auth_url}"
    version = "1.6.0"
}

provider "opentelekomcloud" {
    # please follow the documentation to change options
    alias = "iam"
    user_name  = "your user_name"
    domain_name = "your domain_name" 
    password  = "your passwd"		
    insecure = "true"
    region = "eu-de"
    auth_url  = "https://iam.eu-de.otc.t-systems.com/v3"
    version = "1.6.0"
}

provider "opentelekomcloud" {
    cloud = "test"  # name of cloud in clouds.yaml file
}
