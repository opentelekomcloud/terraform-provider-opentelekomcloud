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
	var jobGet job.Job

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckMrs(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV1JobDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccMRSV1JobConfig_basic,
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
	mrsClient, err := config.MrsV1Client(env.OsRegionName)
	if err != nil {
		return fmt.Errorf("error creating opentelekomcloud mrs: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_mrs_job_v1" {
			continue
		}

		_, err := job.Get(mrsClient, rs.Primary.ID).Extract()
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

func testAccCheckMRSV1JobExists(n string, jobGet *job.Job) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s. ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set. ")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		mrsClient, err := config.MrsV1Client(env.OsRegionName)
		if err != nil {
			return fmt.Errorf("error creating opentelekomcloud mrs client: %s ", err)
		}

		found, err := job.Get(mrsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("job not found. ")
		}

		*jobGet = *found
		time.Sleep(15 * time.Second)

		return nil
	}
}

var TestAccMRSV1JobConfig_basic = fmt.Sprintf(`
resource "opentelekomcloud_mrs_cluster_v1" "cluster1" {
  cluster_name = "mrs-cluster-acc"
  billing_type = 12
  master_node_num = 2
  core_node_num = 3
  master_node_size = "h1.2xlarge.4.linux.mrs"
  core_node_size = "s1.xlarge.linux.mrs"
  available_zone_id = "%s"
  vpc_id = "%s"
  subnet_id = "%s"
  cluster_version = "MRS 1.5.0"
  volume_type = "SAS"
  volume_size = 100
  safe_mode = 0
  cluster_type = 0
  node_public_cert_name = "KeyPair-ci"
  cluster_admin_secret = ""
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
  job_type = 1
  job_name = "test_mapreduce_job1"
  cluster_id = opentelekomcloud_mrs_cluster_v1.cluster1.id
  jar_path = "s3a://tf-mrs-test/program/hadoop-mapreduce-examples-2.7.5.jar"
  input = "s3a://tf-mrs-test/input/"
  output = "s3a://tf-mrs-test/output/"
  job_log = "s3a://tf-mrs-test/joblog/"
  arguments = "wordcount"
}`, env.OsAvailabilityZone, env.OsRouterID, env.OsNetworkID)
