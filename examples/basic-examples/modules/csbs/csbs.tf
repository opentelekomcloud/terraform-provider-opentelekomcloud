
resource "opentelekomcloud_csbs_backup_v1" "backup_v1" {
  backup_name   = "${var.backup_name}_csbs"
  resource_id   = var.resource_id
  description   = var.backup_desc
  resource_type = "OS::Nova::Server"
  tags {
      key   = "k1"
      value = "v1"
    }
}

resource "opentelekomcloud_csbs_backup_v1" "backup_2" {
  resource_id = var.resource_id
}
