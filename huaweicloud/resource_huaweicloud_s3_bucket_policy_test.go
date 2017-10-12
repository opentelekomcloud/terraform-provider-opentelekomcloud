package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jen20/awspolicyequivalence"
)

// PASS
func TestAccAWSS3BucketPolicy_basic(t *testing.T) {
	name := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())

	expectedPolicyText := fmt.Sprintf(
		`{"Version":"2008-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:*"],"Resource":["arn:aws:s3:::%s/*","arn:aws:s3:::%s"]}]}`,
		name, name)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3BucketPolicyConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("huaweicloud_s3_bucket.bucket"),
					testAccCheckAWSS3BucketHasPolicy("huaweicloud_s3_bucket.bucket", expectedPolicyText),
				),
			},
		},
	})
}

// PASS
func TestAccAWSS3BucketPolicy_policyUpdate(t *testing.T) {
	name := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())

	expectedPolicyText1 := fmt.Sprintf(
		`{"Version":"2008-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:*"],"Resource":["arn:aws:s3:::%s/*","arn:aws:s3:::%s"]}]}`,
		name, name)

	expectedPolicyText2 := fmt.Sprintf(
		`{"Version":"2008-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:DeleteBucket", "s3:ListBucket", "s3:ListBucketVersions"], "Resource":["arn:aws:s3:::%s/*","arn:aws:s3:::%s"]}]}`,
		name, name)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3BucketPolicyConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("huaweicloud_s3_bucket.bucket"),
					testAccCheckAWSS3BucketHasPolicy("huaweicloud_s3_bucket.bucket", expectedPolicyText1),
				),
			},

			{
				Config: testAccAWSS3BucketPolicyConfig_updated(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketExists("huaweicloud_s3_bucket.bucket"),
					testAccCheckAWSS3BucketHasPolicy("huaweicloud_s3_bucket.bucket", expectedPolicyText2),
				),
			},
		},
	})
}

func testAccCheckAWSS3BucketHasPolicy(n string, expectedPolicyText string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No S3 Bucket ID is set")
		}

		conn := testAccProvider.Meta().(*Config).s3conn

		policy, err := conn.GetBucketPolicy(&s3.GetBucketPolicyInput{
			Bucket: aws.String(rs.Primary.ID),
		})
		if err != nil {
			return fmt.Errorf("GetBucketPolicy error: %v", err)
		}

		actualPolicyText := *policy.Policy

		equivalent, err := awspolicy.PoliciesAreEquivalent(actualPolicyText, expectedPolicyText)
		if err != nil {
			return fmt.Errorf("Error testing policy equivalence: %s", err)
		}
		if !equivalent {
			return fmt.Errorf("Non-equivalent policy error:\n\nexpected: %s\n\n     got: %s\n",
				expectedPolicyText, actualPolicyText)
		}

		return nil
	}
}

func testAccAWSS3BucketPolicyConfig(bucketName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_s3_bucket" "bucket" {
	bucket = "%s"
	#tags {
	#	TestName = "TestAccAWSS3BucketPolicy_basic"
	#}
}

resource "huaweicloud_s3_bucket_policy" "bucket" {
	bucket = "${huaweicloud_s3_bucket.bucket.bucket}"
	policy =<<POLICY
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

func testAccAWSS3BucketPolicyConfig_updated(bucketName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_s3_bucket" "bucket" {
	bucket = "%s"
	#tags {
	#	TestName = "TestAccAWSS3BucketPolicy_basic"
	#}
}

resource "huaweicloud_s3_bucket_policy" "bucket" {
	bucket = "${huaweicloud_s3_bucket.bucket.bucket}"
	policy =<<POLICY
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
