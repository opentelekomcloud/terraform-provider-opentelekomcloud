#create SFS
resource "opentelekomcloud_sfs_file_system_v2" "sfs_1" {
  size              = "${var.size}"
  name              = "${var.share_name}"
  access_to         = "${var.vpc_id}"
  access_level      = "rw"
  description       = "${var.share_description}"
  availability_zone = "${var.availability_zone}"
  metadata = {
    "type" = "nfs"
  }
  share_proto = "NFS"
  is_public   = false
  #access_type = "ip"

}
resource "opentelekomcloud_sfs_file_system_v2" "sfs_2" {
  size         = "${var.size}"
  access_to    = "${var.vpc_id}"
  access_level = "rw"
}
