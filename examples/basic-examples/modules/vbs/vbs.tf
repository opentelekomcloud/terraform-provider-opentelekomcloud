resource "opentelekomcloud_vbs_backup_v2" "mybackup2" {
  volume_id   = var.volume_id
  name        = "${var.backup_name}_test"
  description = var.backup_desc
  tags {
      key   = "k1"
      value = "v1"
  }
}

resource "opentelekomcloud_vbs_backup_policy_v2" "policy1" {
  name                = "policy_001"
  start_time          = "12:00"
  status              = "ON"
  retain_first_backup = "N"
  rentention_num      = 2
  frequency           = 1
  tags {
      key   = "k1"
      value = "v1"
  }
}

resource "opentelekomcloud_vbs_backup_share_v2" "backupshare" {
  backup_id      = opentelekomcloud_vbs_backup_v2.mybackup2.id
  to_project_ids = [var.to_project_ids]
}
