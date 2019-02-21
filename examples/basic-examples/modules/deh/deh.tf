resource "opentelekomcloud_deh_host_v1" "deh_host" 
    {
        name = "${var.deh_name}"
        auto_placement = "on"
        availability_zone = "${var.availability_zone}"
        host_type = "${var.host_type}"
    }
resource "opentelekomcloud_deh_host_v1" "deh_host2" 
    {
        name = "${var.deh_name}_required"
        availability_zone = "${var.availability_zone}"
        host_type = "general"
    }