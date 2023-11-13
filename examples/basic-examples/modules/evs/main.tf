resource "opentelekomcloud_blockstorage_volume_v2" "volume_2" {
  #region      = var.region
  #name        = var.volume_name
  #availability_zone = var.availability_zone
  #description = var.desc
  size = var.volume_size
  #volume_type = var.volume_type
}

resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  region            = var.region
  name              = var.volume_name
  availability_zone = var.availability_zone
  description       = var.desc
  size              = var.volume_size
  volume_type       = var.volume_type
}

#attach single evs to ecs
resource "opentelekomcloud_compute_volume_attach_v2" "va_1" {
  instance_id = var.ecs_id
  volume_id   = opentelekomcloud_blockstorage_volume_v2.volume_1.id
}

#attach multiple volumes to ecs
resource "opentelekomcloud_blockstorage_volume_v2" "volumes" {
  count             = 2
  availability_zone = var.availability_zone
  name              = format("vol-%02d", count.index + 1)
  size              = var.volume_size
}
resource "opentelekomcloud_compute_volume_attach_v2" "attachments" {
  count       = 2
  instance_id = var.ecs_id
  volume_id   = element(opentelekomcloud_blockstorage_volume_v2.volumes.*.id, count.index)
}
