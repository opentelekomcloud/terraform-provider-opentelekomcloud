     resource "opentelekomcloud_networking_secgroup_v2" "secgroup_dcs" {
         name = "secgroup_dcs"
         description = "secgroup_dcs"
       }
      data "opentelekomcloud_dcs_az_v1" "az_1" {
         port = "8002"
        }
       data "opentelekomcloud_dcs_product_v1" "product_1" {
          spec_code = "dcs.master_standby"
        }

resource "opentelekomcloud_dcs_instance_v1" "instance_1" {
          name  = "${var.dcs_name}_required"
          capacity = "${var.capacity}"
          vpc_id = "${var.vpc_id}"
          security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_dcs.id}"
          subnet_id = "${var.subnet_id}"
          available_zones = ["${data.opentelekomcloud_dcs_az_v1.az_1.id}"]
          #product_id = "${data.opentelekomcloud_dcs_product_v1.product_1.id}"
		  product_id = "OTC_DCS_MS"
		  save_days = 1
          backup_type = "manual"
          begin_at = "00:00-01:00"
          period_type = "weekly"
		  backup_at = [1]
		  engine = "Redis"
		  password = "Huawei@123"
		  engine_version = "3.0.7"    
        }
resource "opentelekomcloud_dcs_instance_v1" "instance_2" {
          name  = "${var.dcs_name}_redis"
		  description  = "${var.dcs_desc}"		  
          engine_version = "3.0.7"         
          engine = "Redis"
          capacity = 4
		  access_user= ""
		  password = "Huawei_test"
          vpc_id = "${var.vpc_id}"
          security_group_id = "${opentelekomcloud_networking_secgroup_v2.secgroup_dcs.id}"
          subnet_id = "${var.subnet_id}"
          available_zones = ["${data.opentelekomcloud_dcs_az_v1.az_1.id}"]
          # product_id = "${data.opentelekomcloud_dcs_product_v1.product_1.id}"
		  product_id = "OTC_DCS_MS"
          save_days = 1
          backup_type = "manual"
          begin_at = "00:00-01:00"
          period_type = "weekly"
		  maintain_begin  = "02:00"
		  maintain_end  = "06:00"
          backup_at = [1]
          depends_on = ["data.opentelekomcloud_dcs_product_v1.product_1", "opentelekomcloud_networking_secgroup_v2.secgroup_dcs"]
}
