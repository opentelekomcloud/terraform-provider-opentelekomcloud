resource "opentelekomcloud_csbs_backup_v1" "backup_v1" {
  backup_name   = "${var.project}-backup"
  description   = "mybackup"
  resource_id   = opentelekomcloud_compute_instance_v2.webserver.id
  resource_type = "OS::Nova::Server"
}
