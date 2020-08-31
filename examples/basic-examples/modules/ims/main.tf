
resource "opentelekomcloud_images_image_v2" "rancheros" {
  name             = "RancherOS"
  image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
  container_format = "bare"
  disk_format      = "qcow2"
}

#resource "opentelekomcloud_images_image_v2" "BIGIP" {
#  name   = "CentOS7_created_by_terraform"
#  local_file_path = "/opt/terraform/opentelekomcloud/BIGIP-14.0.0.3-0.0.4.qcow2"
#  container_format = "bare"
#  disk_format = "qcow2"
#  tags = ["foo.bar", "tag.value"]
#  min_disk_gb = 40
#  min_ram_mb = 10240
#  visibility = "private"
#}
