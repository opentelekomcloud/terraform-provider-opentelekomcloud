resource "opentelekomcloud_mrs_cluster_v1" "cluster1" {
  cluster_name          = "${var.mrs_name}"
  region                = "eu-de"
  billing_type          = 12
  master_node_num       = 2
  core_node_num         = 3
  master_node_size      = "s1.xlarge.linux.mrs"
  core_node_size        = "s1.xlarge.linux.mrs"
  available_zone_id     = "eu-de-01"
  vpc_id                = "${var.vpc_id}"
  subnet_id             = "${var.subnet_id}"
  cluster_version       = "MRS 1.7.2"
  volume_type           = "SATA"
  volume_size           = 100
  safe_mode             = 1
  cluster_type          = 0
  node_public_cert_name = "${var.keypair_name}"
  cluster_admin_secret  = "t-systems@32158"
  component_list {
    component_name = "Hadoop"
  }
  component_list {
    component_name = "Spark"
  }
  component_list {
    component_name = "Hive"
  }
}

resource "opentelekomcloud_mrs_job_v1" "job1" {
  job_type   = 1
  job_name   = "${var.job_name}"
  cluster_id = "${opentelekomcloud_mrs_cluster_v1.cluster1.id}"
  jar_path   = "s3a://wordcount/program/hadoop-mapreduce-examples-2.7.5.jar"
  input      = "s3a://wordcount/input/"
  output     = "s3a://wordcount/output/"
  job_log    = "s3a://wordcount/log/"
  arguments  = "wordcount"
}
