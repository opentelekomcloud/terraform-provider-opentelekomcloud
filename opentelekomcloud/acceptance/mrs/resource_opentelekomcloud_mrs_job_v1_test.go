package acceptance

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/mrs/v1/job"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccMRSV1Job_basic(t *testing.T) {
	var jobGet job.JobExecution

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheckRequiredEnvVars(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV1JobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMRSV1JobConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV1JobExists("opentelekomcloud_mrs_job_v1.job1", &jobGet),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_mrs_job_v1.job1", "job_state", "Completed"),
				),
			},
		},
	})
}

func testAccCheckMRSV1JobDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	mrsClient, err := config.MrsV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating opentelekomcloud mrs: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_mrs_job_v1" {
			continue
		}

		_, err := job.Get(mrsClient, rs.Primary.ID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault400); ok {
				return nil
			}
			if _, ok := err.(golangsdk.ErrDefault500); ok {
				continue
			}
			return fmt.Errorf("job still exists. err : %s", err)
		}
	}

	return nil
}

func testAccCheckMRSV1JobExists(n string, jobGet *job.JobExecution) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s. ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set. ")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		mrsClient, err := config.MrsV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating opentelekomcloud mrs client: %s ", err)
		}

		found, err := job.Get(mrsClient, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("job not found. ")
		}

		*jobGet = *found
		time.Sleep(15 * time.Second)

		return nil
	}
}

var testAccMRSV1JobConfigBasic = fmt.Sprintf(`
%s

resource "opentelekomcloud_mrs_cluster_v1" "cluster1" {
  cluster_name             = "mrs-cluster-acc"
  billing_type             = 12
  master_node_num          = 1
  core_node_num            = 1
  master_node_size         = "c3.2xlarge.4.linux.mrs"
  core_node_size           = "c3.2xlarge.4.linux.mrs"
  available_zone_id        = "%s"
  vpc_id                   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id                = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  cluster_version          = "MRS 1.7.2"
  master_data_volume_type  = "SAS"
  master_data_volume_size  = 100
  master_data_volume_count = 1
  core_data_volume_type    = "SATA"
  core_data_volume_size    = 100
  core_data_volume_count   = 2
  safe_mode                = 0
  cluster_type             = 0
  node_public_cert_name    = "%s"
  cluster_admin_secret     = "SuperQwert!123"
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
  job_name   = "test_mapreduce_job1"
  cluster_id = "9e9dfe0f-2b8a-4b40-bd49-2323b8268a47"
  jar_path   = "s3a://tf-mrs-test/program/hadoop-mapreduce-examples-2.7.5.jar"
  input      = "s3a://tf-mrs-test/input/"
  output     = "s3a://tf-mrs-test/output/"
  job_log    = "s3a://tf-mrs-test/joblog/"
  arguments  = "wordcount"
}
`, common.DataSourceSubnet, env.OS_AVAILABILITY_ZONE, env.OS_KEYPAIR_NAME)
