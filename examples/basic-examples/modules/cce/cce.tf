# test required params
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name                   = "${var.cce_name}-1"
  cluster_type           = var.cce_type
  flavor_id              = "cce.s2.small"
  vpc_id                 = var.vpc_id
  subnet_id              = var.subnet_id
  container_network_type = "overlay_l2"
}
resource "opentelekomcloud_cce_cluster_v3" "cluster_2" {
  name                   = var.cce_name
  cluster_type           = var.cce_type
  flavor_id              = "cce.s2.small"
  vpc_id                 = var.vpc_id
  subnet_id              = var.subnet_id
  container_network_type = "overlay_l2"
  billing_mode           = 0
  cluster_version        = "v1.9.2"
  description            = "Create cluster by terraform"
}
resource "opentelekomcloud_cce_node_v3" "node_1" {
  cluster_id        = opentelekomcloud_cce_cluster_v3.cluster_2.id
  name              = var.cce_node_name
  flavor_id         = "m1.large"
  iptype            = "5_bgp"
  availability_zone = "eu-de-01"
  key_pair          = var.key_name
  root_volume {
    size       = 40,
    volumetype = "SATA"
  }
  sharetype      = "PER"
  bandwidth_size = 100

  data_volumes {
    size       = 100,
    volumetype = "SATA"
  }

  max_pods = 2

}
