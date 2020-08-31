#create KMS
resource "opentelekomcloud_kms_key_v1" "key_1" {
  key_alias = "${var.kms_name}"
}
resource "opentelekomcloud_kms_key_v1" "key_2" {
  key_alias       = "${var.kms_name}_2"
  pending_days    = "${var.kms_pending_days}"
  key_description = "key"
  realm           = "${var.kms_region}"
  is_enabled      = true
}
