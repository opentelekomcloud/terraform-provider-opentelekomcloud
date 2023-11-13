#create SFS
resource "opentelekomcloud_sfs_file_system_v2" "sfs_1" {
  size              = var.size
  name              = var.share_name
  description       = var.share_description
  availability_zone = var.availability_zone
  metadata = {
    "type" = "nfs"
  }
  share_proto = "NFS"
  is_public   = false
  #access_type = "ip"

}

resource "opentelekomcloud_sfs_file_system_v2" "share-file" {
  name        = var.share_name
  size        = 50
  description = var.share_description
  share_proto = "NFS"

  tags = {
    muh = "kuh"
  }
}
