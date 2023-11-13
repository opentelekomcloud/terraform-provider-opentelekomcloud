data "opentelekomcloud_rds_flavors_v1" "flavor" {
  region            = var.region
  datastore_name    = var.datastore_name
  datastore_version = var.datastore_version
  speccode          = var.flavor_name
}

resource "opentelekomcloud_rds_instance_v1" "instance" {
  name = var.rds_name
  datastore {
    type    = var.datastore_name
    version = var.datastore_version
  }
  #flavorref = data.opentelekomcloud_rds_flavors_v1.flavor.id
  flavorref = "6ba2c53b-386b-41bf-bc70-72cc78a867a6"
  volume {
    type = var.rds_volume_type
    size = var.rds_volume_size
  }
  region           = var.region
  availabilityzone = var.availabilityzone
  vpc              = var.vpc_id
  nics {
    subnetid = var.subnetid
  }
  securitygroup {
    id = var.securitygroupid
  }
  dbport = "8635"
  dbrtpd = var.passwd
  backupstrategy {
    starttime = "04:00:00"
    keepdays  = 4
  }
  ha {
    enable          = true
    replicationmode = "async"
  }
}
