package acceptance

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	awspolicy "github.com/jen20/awspolicyequivalence"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccS3BucketPolicy_basic(t *testing.T) {
	name := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())

	expectedPolicyText := fmt.Sprintf(
		`{"Version":"2008-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:*"],"Resource":["arn:aws:s3:::%s/*","arn:aws:s3:::%s"]}]}`,
		name, name)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketPolicyConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					testAccCheckS3BucketHasPolicy("opentelekomcloud_s3_bucket.bucket", expectedPolicyText),
				),
			},
		},
	})
}

func TestAccS3BucketPolicy_policyUpdate(t *testing.T) {
	name := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())

	expectedPolicyText1 := fmt.Sprintf(
		`{"Version":"2008-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:*"],"Resource":["arn:aws:s3:::%s/*","arn:aws:s3:::%s"]}]}`,
		name, name)

	expectedPolicyText2 := fmt.Sprintf(
		`{"Version":"2008-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:DeleteBucket", "s3:ListBucket", "s3:ListBucketVersions"], "Resource":["arn:aws:s3:::%s/*","arn:aws:s3:::%s"]}]}`,
		name, name)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketPolicyConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					testAccCheckS3BucketHasPolicy("opentelekomcloud_s3_bucket.bucket", expectedPolicyText1),
				),
			},

			{
				Config: testAccS3BucketPolicyConfig_updated(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("opentelekomcloud_s3_bucket.bucket"),
					testAccCheckS3BucketHasPolicy("opentelekomcloud_s3_bucket.bucket", expectedPolicyText2),
				),
			},
		},
	})
}

func testAccCheckS3BucketHasPolicy(n string, expectedPolicyText string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no S3 Bucket ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		conn, err := config.S3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud s3 client: %s", err)
		}

		policy, err := conn.GetBucketPolicy(&s3.GetBucketPolicyInput{
			Bucket: aws.String(rs.Primary.ID),
		})
		if err != nil {
			return fmt.Errorf("getBucketPolicy error: %v", err)
		}

		actualPolicyText := *policy.Policy

		equivalent, err := awspolicy.PoliciesAreEquivalent(actualPolicyText, expectedPolicyText)
		if err != nil {
			return fmt.Errorf("error testing policy equivalence: %s", err)
		}
		if !equivalent {
			return fmt.Errorf("non-equivalent policy error:\n\nexpected: %s\n\n     got: %s\n",
				expectedPolicyText, actualPolicyText)
		}

		return nil
	}
}

func testAccS3BucketPolicyConfig(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "%s"
}

resource "opentelekomcloud_s3_bucket_policy" "bucket" {
  bucket = opentelekomcloud_s3_bucket.bucket.bucket
  policy = <<POLICY
{
	"Version": "2008-10-17",
	"Statement": [{
		"Effect": "Allow",
		"Principal": {
			"AWS": ["*"]
		},
		"Action": [
			"s3:*"
		],
		"Resource": [
			"arn:aws:s3:::%s",
			"arn:aws:s3:::%s/*"
		]
	}]
}
POLICY
}
`, bucketName, bucketName, bucketName)
}

func testAccS3BucketPolicyConfig_updated(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "%s"
}

resource "opentelekomcloud_s3_bucket_policy" "bucket" {
  bucket = opentelekomcloud_s3_bucket.bucket.bucket
  policy = <<POLICY
{
	"Version": "2008-10-17",
	"Statement": [{
		"Effect": "Allow",
		"Principal": {
			"AWS": ["*"]
		},
		"Action": [
			"s3:DeleteBucket",
			"s3:ListBucket",
			"s3:ListBucketVersions"
		],
		"Resource": [
			"arn:aws:s3:::%s",
			"arn:aws:s3:::%s/*"
		]
	}]
}
POLICY
}
`, bucketName, bucketName, bucketName)
}
