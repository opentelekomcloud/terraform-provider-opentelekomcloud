package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	awspolicy "github.com/jen20/awspolicyequivalence"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceName = "opentelekomcloud_obs_bucket.bucket"

func TestAccObsBucketPolicyBasic(t *testing.T) {
	name := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())

	expectedPolicyText := fmt.Sprintf(
		`{"Version":"2008-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:*"],"Resource":["arn:aws:s3:::%s/*","arn:aws:s3:::%s"]}]}`,
		name, name)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketPolicyConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					testAccCheckObsBucketHasPolicy(resourceName, expectedPolicyText),
				),
			},
		},
	})
}

func TestAccObsBucketPolicyUpdate(t *testing.T) {
	name := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())

	expectedPolicyText1 := fmt.Sprintf(
		`{"Version":"2008-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:*"],"Resource":["arn:aws:s3:::%s/*","arn:aws:s3:::%s"]}]}`,
		name, name)

	expectedPolicyText2 := fmt.Sprintf(
		`{"Version":"2008-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:DeleteBucket", "s3:ListBucket", "s3:ListBucketVersions"], "Resource":["arn:aws:s3:::%s/*","arn:aws:s3:::%s"]}]}`,
		name, name)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketPolicyConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					testAccCheckObsBucketHasPolicy(resourceName, expectedPolicyText1),
				),
			},

			{
				Config: testAccObsBucketPolicyConfigUpdated(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					testAccCheckObsBucketHasPolicy(resourceName, expectedPolicyText2),
				),
			},
		},
	})
}

func TestAccObsBucketPolicyMalformed(t *testing.T) {
	name := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { common.TestAccPreCheck(t) },
		Providers:    common.TestAccProviders,
		CheckDestroy: testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccObsBucketPolicyConfigWrongPolicy(name),
				ExpectError: regexp.MustCompile(`error putting OBS policy.+`),
			},
		},
	})
}

func testAccCheckObsBucketHasPolicy(n string, expectedPolicyText string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no OBS Bucket ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.NewObjectStorageClient(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OBS client: %s", err)
		}

		policy, err := client.GetBucketPolicy(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("GetBucketPolicy error: %v", err)
		}

		equivalent, err := awspolicy.PoliciesAreEquivalent(policy.Policy, expectedPolicyText)
		if err != nil {
			return fmt.Errorf("error testing policy equivalence: %s", err)
		}
		if !equivalent {
			return fmt.Errorf("non-equivalent policy error:\n\nexpected: %s\n\n     got: %s\n",
				expectedPolicyText, policy.Policy)
		}

		return nil
	}
}

func testAccObsBucketPolicyConfig(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket = "%s"
}

resource "opentelekomcloud_obs_bucket_policy" "bucket" {
  bucket = opentelekomcloud_obs_bucket.bucket.bucket
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

func testAccObsBucketPolicyConfigUpdated(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket = "%s"
}

resource "opentelekomcloud_obs_bucket_policy" "bucket" {
  bucket = opentelekomcloud_obs_bucket.bucket.bucket
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

func testAccObsBucketPolicyConfigWrongPolicy(bucketName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "bucket" {
  bucket = "%s"
}

resource "opentelekomcloud_obs_bucket_policy" "bucket" {
  bucket = opentelekomcloud_obs_bucket.bucket.bucket
  policy =<<POLICY
{
    "Id": "BUCKET_POLICY",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "ID": [
                    "domain/12345:user/67890"
                ]
            },
            "Action": [
                "s3:GetObject",
                "s3:PutObject",
                "s3:ListBucket",
                "s3:ListBucketVersions",
                "s3:ListBucketMultipartUploads",
                "s3:GetBucketLocation",
                "s3:GetBucketStorage"
            ],
            "Resource": [
                "arn:aws:s3:::bucket/*",
                "arn:aws:s3:::bucket"
            ]
        }
    ]
}
POLICY
}
`, bucketName)
}
