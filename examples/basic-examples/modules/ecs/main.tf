# create random number
resource "random_id" "ecs" {
  byte_length = 4
}
# create ECS with Required params
resource "opentelekomcloud_compute_instance_v2" "ecs_1" {
  region            = var.region
  availability_zone = var.availability_zone
  name              = "${var.ecs_name}-notags"
  image_id          = var.image_id
  flavor_id         = var.flavor_id

  network {
    uuid = var.subnet_id
  }
}

# create ECS with required and  optional params
resource "opentelekomcloud_compute_instance_v2" "ecs_2" {
  region            = var.region
  availability_zone = var.availability_zone
  name              = "${var.ecs_name}-${random_id.ecs.id}"
  image_id          = var.image_id
  flavor_id         = var.flavor_id
  key_pair          = var.key_name
  security_groups   = [var.security_groups]
  user_data         = "#cloud-config\nhostname: instance_1.example.com\nfqdn: instance_1.example.com"
  config_drive      = "true"
  network {
    uuid = var.subnet_id
  }
  network {
    uuid = var.subnet_id2
  }
  admin_pass = var.admin_pass

  tags = {
    foo  = "bar1"
    key  = "value"
    key2 = "value2"
  }
  auto_recovery = var.auto_recovery
}


resource "opentelekomcloud_networking_floatingip_v2" "floatip_1" {
  pool = "admin_external_net"
}

resource "opentelekomcloud_compute_floatingip_associate_v2" "fip_2" {
  floating_ip = opentelekomcloud_networking_floatingip_v2.floatip_1.address
  instance_id = opentelekomcloud_compute_instance_v2.ecs_2.id
}

# boot from volume
resource "opentelekomcloud_compute_instance_v2" "boot-from-volume" {
  region            = var.region
  availability_zone = var.availability_zone
  name              = "${var.ecs_name}-${random_id.ecs.id}"
  #image_id         = var.image_id
  flavor_id       = var.flavor_id
  key_pair        = var.key_name
  security_groups = [var.security_groups]

  block_device {
    uuid                  = var.image_id
    source_type           = "image"
    volume_size           = 1000
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
  }

  network {
    uuid = var.subnet_id
  }
  tags = {
    tag = "test"
  }
}

#Boot Instance and Attach Existing Volume as a Block Device

resource "opentelekomcloud_blockstorage_volume_v2" "myvol" {
  name = "myvol"
  size = 10
}
resource "opentelekomcloud_compute_instance_v2" "instance" {
  region            = var.region
  availability_zone = var.availability_zone
  name              = "${var.ecs_name}-${random_id.ecs.id}"
  #image_id         = var.image_id

  flavor_id       = var.flavor_id
  key_pair        = var.key_name
  security_groups = [var.security_groups]

  block_device {
    uuid                  = var.image_id
    source_type           = "image"
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
    volume_size           = 10
  }

  #block_device {
  #  uuid                  = opentelekomcloud_blockstorage_volume_v2.myvol.id
  #  source_type           = "volume"
  #  boot_index            = 1
  #  destination_type      = "volume"
  #  delete_on_termination = true
  #}

  network {
    uuid = var.subnet_id
  }
}

##Boot Instance, Create Volume, and Attach Volume as a Block Device
#resource "opentelekomcloud_compute_instance_v2" "instance_1" {
#  region = var.region
#  availability_zone = var.availability_zone
#  name            = var.ecs_name}-${random_id.ecs.id
#  image_id       = var.image_id
#  flavor_id       = var.flavor_id
#  key_pair        = var.key_name
#  security_groups = [var.security_groups]
#
#  block_device {
#    uuid                  = var.image_id
#    source_type           = "image"
#    boot_index            = 0
#    destination_type      = "local"
#    delete_on_termination = true
#  }
#  block_device {
#    source_type           = "blank"
#    boot_index            = 1
#    destination_type      = "volume"
#    volume_size           = 10
#    delete_on_termination = true
#  }
#
#  network {
#    uuid = var.subnet_id
#  }
#}
