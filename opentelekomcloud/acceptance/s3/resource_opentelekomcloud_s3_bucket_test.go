package acceptance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccS3Bucket_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "opentelekomcloud_s3_bucket.bucket"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "region", env.OS_REGION_NAME),
					resource.TestCheckNoResourceAttr(resourceName, "website_endpoint"),
					resource.TestCheckResourceAttr(resourceName, "bucket", testAccBucketName(rInt)),
					resource.TestCheckResourceAttr(resourceName, "bucket_domain_name", testAccBucketDomainName(rInt)),
				),
			},
		},
	})
}

func TestAccAWSS3MultiBucket_withTags(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3MultiBucketConfigWithTags(rInt),
			},
		},
	})
}

func TestAccS3Bucket_namePrefix(t *testing.T) {
	resourceName := "opentelekomcloud_s3_bucket.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketConfig_namePrefix,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					resource.TestMatchResourceAttr(resourceName, "bucket", regexp.MustCompile("^tf-test-")),
				),
			},
		},
	})
}

func TestAccS3Bucket_generatedName(t *testing.T) {
	resourceName := "opentelekomcloud_s3_bucket.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketConfig_generatedName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
				),
			},
		},
	})
}

func TestAccS3Bucket_region(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "opentelekomcloud_s3_bucket.bucket"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketConfigWithRegion(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "region", env.OS_REGION_NAME),
				),
			},
		},
	})
}

func TestAccS3Bucket_Policy(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "opentelekomcloud_s3_bucket.bucket"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketConfigWithPolicy(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					testAccCheckS3BucketPolicy(resourceName, testAccS3BucketPolicy(rInt)),
				),
			},
			{
				Config: testAccS3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					testAccCheckS3BucketPolicy(resourceName, ""),
				),
			},
			{
				Config: testAccS3BucketConfigWithEmptyPolicy(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					testAccCheckS3BucketPolicy(resourceName, ""),
				),
			},
		},
	})
}

func TestAccS3Bucket_UpdateAcl(t *testing.T) {
	ri := acctest.RandInt()
	preConfig := fmt.Sprintf(testAccS3BucketConfigWithAcl, ri)
	postConfig := fmt.Sprintf(testAccS3BucketConfigWithAclUpdate, ri)
	resourceName := "opentelekomcloud_s3_bucket.bucket"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "acl", "public-read"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "acl", "private"),
				),
			},
		},
	})
}

func TestAccS3Bucket_Website_Simple(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "opentelekomcloud_s3_bucket.bucket"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketWebsiteConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					testAccCheckS3BucketWebsite(resourceName, "index.html", "", "", ""),
					resource.TestCheckResourceAttr(resourceName, "website_endpoint", testAccWebsiteEndpoint(rInt)),
				),
			},
			{
				Config: testAccS3BucketWebsiteConfigWithError(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					testAccCheckS3BucketWebsite(resourceName, "index.html", "error.html", "", ""),
					resource.TestCheckResourceAttr(resourceName, "website_endpoint", testAccWebsiteEndpoint(rInt)),
				),
			},
			{
				Config: testAccS3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					testAccCheckS3BucketWebsite(resourceName, "", "", "", ""),
					resource.TestCheckResourceAttr(resourceName, "website_endpoint", ""),
				),
			},
		},
	})
}

func TestAccS3Bucket_WebsiteRedirect(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "opentelekomcloud_s3_bucket.bucket"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketWebsiteConfigWithRedirect(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					testAccCheckS3BucketWebsite(resourceName, "", "", "", "hashicorp.com"),
					resource.TestCheckResourceAttr(resourceName, "website_endpoint", testAccWebsiteEndpoint(rInt)),
				),
			},
			{
				Config: testAccS3BucketWebsiteConfigWithHttpsRedirect(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					testAccCheckS3BucketWebsite(resourceName, "", "", "https", "hashicorp.com"),
					resource.TestCheckResourceAttr(resourceName, "website_endpoint", testAccWebsiteEndpoint(rInt)),
				),
			},
			{
				Config: testAccS3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					testAccCheckS3BucketWebsite(resourceName, "", "", "", ""),
					resource.TestCheckResourceAttr(resourceName, "website_endpoint", ""),
				),
			},
		},
	})
}

func TestAccS3Bucket_WebsiteRoutingRules(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "opentelekomcloud_s3_bucket.bucket"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketWebsiteConfigWithRoutingRules(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					testAccCheckS3BucketWebsite(resourceName, "index.html", "error.html", "", ""),
					testAccCheckS3BucketWebsiteRoutingRules(resourceName,
						[]*s3.RoutingRule{
							{
								Condition: &s3.Condition{
									KeyPrefixEquals: aws.String("docs/"),
								},
								Redirect: &s3.Redirect{
									ReplaceKeyPrefixWith: aws.String("documents/"),
								},
							},
						},
					),
					resource.TestCheckResourceAttr(resourceName, "website_endpoint", testAccWebsiteEndpoint(rInt)),
				),
			},
			{
				Config: testAccS3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					testAccCheckS3BucketWebsite(resourceName, "", "", "", ""),
					testAccCheckS3BucketWebsiteRoutingRules(resourceName, nil),
					resource.TestCheckResourceAttr(resourceName, "website_endpoint", ""),
				),
			},
		},
	})
}

// Test TestAccAWSS3Bucket_shouldFailNotFound is designed to fail with a "plan
// not empty" error in Terraform, to check against regressions.
// See https://github.com/hashicorp/terraform/pull/2925
func TestAccS3Bucket_shouldFailNotFound(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "opentelekomcloud_s3_bucket.bucket"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketDestroyedConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					testAccCheckS3DestroyBucket(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccS3Bucket_Versioning(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "opentelekomcloud_s3_bucket.bucket"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					testAccCheckS3BucketVersioning(resourceName, ""),
				),
			},
			{
				Config: testAccS3BucketConfigWithVersioning(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					testAccCheckS3BucketVersioning(resourceName, s3.BucketVersioningStatusEnabled),
				),
			},
			{
				Config: testAccS3BucketConfigWithDisableVersioning(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					testAccCheckS3BucketVersioning(resourceName, s3.BucketVersioningStatusSuspended),
				),
			},
		},
	})
}

func TestAccS3Bucket_VersioningSecond(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "opentelekomcloud_s3_bucket.bucket"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketConfigWithVersioning(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					testAccCheckS3BucketVersioning(resourceName, s3.BucketVersioningStatusEnabled),
				),
			},
			{
				Config: testAccS3BucketConfigWithDisableVersioning(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					testAccCheckS3BucketVersioning(resourceName, s3.BucketVersioningStatusSuspended),
				),
			},
		},
	})
}

func TestAccS3Bucket_Cors(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "opentelekomcloud_s3_bucket.bucket"

	updateBucketCors := func(n string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			rs, ok := s.RootModule().Resources[n]
			if !ok {
				return fmt.Errorf("not found: %s", n)
			}

			config := common.TestAccProvider.Meta().(*cfg.Config)
			client, err := config.S3Client(env.OS_REGION_NAME)
			if err != nil {
				return fmt.Errorf("error creating OpenTelekomCloud S3 client: %s", err)
			}
			_, err = client.PutBucketCors(&s3.PutBucketCorsInput{
				Bucket: aws.String(rs.Primary.ID),
				CORSConfiguration: &s3.CORSConfiguration{
					CORSRules: []*s3.CORSRule{
						{
							AllowedHeaders: []*string{aws.String("*")},
							AllowedMethods: []*string{aws.String("GET")},
							AllowedOrigins: []*string{aws.String("https://www.example.com")},
						},
					},
				},
			})
			if err != nil {
				if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() != "NoSuchCORSConfiguration" {
					return err
				}
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketConfigWithCORS(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					testAccCheckS3BucketCors(resourceName,
						[]*s3.CORSRule{
							{
								AllowedHeaders: []*string{aws.String("*")},
								AllowedMethods: []*string{aws.String("PUT"), aws.String("POST")},
								AllowedOrigins: []*string{aws.String("https://www.example.com")},
								ExposeHeaders:  []*string{aws.String("x-amz-server-side-encryption"), aws.String("ETag")},
								MaxAgeSeconds:  aws.Int64(3000),
							},
						},
					),
					updateBucketCors(resourceName),
				),
				ExpectNonEmptyPlan: true, // TODO: No diff in real life, so maybe a timing problem?
			},
			{
				Config: testAccS3BucketConfigWithCORS(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					testAccCheckS3BucketCors(resourceName,
						[]*s3.CORSRule{
							{
								AllowedHeaders: []*string{aws.String("*")},
								AllowedMethods: []*string{aws.String("PUT"), aws.String("POST")},
								AllowedOrigins: []*string{aws.String("https://www.example.com")},
								ExposeHeaders:  []*string{aws.String("x-amz-server-side-encryption"), aws.String("ETag")},
								MaxAgeSeconds:  aws.Int64(3000),
							},
						},
					),
				),
			},
		},
	})
}

func TestAccS3Bucket_Logging(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "opentelekomcloud_s3_bucket.bucket"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketConfigWithLogging(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					testAccCheckS3BucketLogging(resourceName, "opentelekomcloud_s3_bucket.log_bucket", "log/"),
				),
			},
		},
	})
}

func TestAccS3Bucket_Lifecycle(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "opentelekomcloud_s3_bucket.bucket"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketConfigWithLifecycle(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.0.id", "id1"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.0.expiration.2613713285.days", "365"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.0.expiration.2613713285.date", ""),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.0.expiration.2613713285.expired_object_delete_marker", "false"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.1.id", "id2"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.1.prefix", "path2/"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.1.expiration.2855832418.date", "2016-01-12"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.1.expiration.2855832418.days", "0"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.1.expiration.2855832418.expired_object_delete_marker", "false"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.2.id", "id3"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.2.prefix", "path3/"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.3.id", "id4"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.3.prefix", "path4/"),
				),
			},
			{
				Config: testAccS3BucketConfigWithVersioningLifecycle(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.0.id", "id1"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.0.noncurrent_version_expiration.80908210.days", "365"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.1.id", "id2"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.1.prefix", "path2/"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.1.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.1.noncurrent_version_expiration.80908210.days", "365"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.2.id", "id3"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.2.prefix", "path3/"),
				),
			},
			{
				Config: testAccS3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists(resourceName),
				),
			},
		},
	})
}

func testAccCheckS3BucketDestroy(_ *terraform.State) error {
	// UNDONE: Why instance check?
	// return testAccCheckInstanceDestroyWithProvider(s, TestAccProvider)
	return nil
}

func testAccCheckS3BucketExists(n string) resource.TestCheckFunc {
	providers := []*schema.Provider{common.TestAccProvider}
	return testAccCheckS3BucketExistsWithProviders(n, &providers)
}

func testAccCheckS3BucketExistsWithProviders(n string, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers
			if provider.Meta() == nil {
				continue
			}

			config := common.TestAccProvider.Meta().(*cfg.Config)
			conn, err := config.S3Client(env.OS_REGION_NAME)
			if err != nil {
				return fmt.Errorf("error creating OpenTelekomCloud S3 client: %s", err)
			}
			_, err = conn.HeadBucket(&s3.HeadBucketInput{
				Bucket: aws.String(rs.Primary.ID),
			})

			if err != nil {
				return fmt.Errorf("s3 Bucket error: %v", err)
			}
			return nil
		}

		return fmt.Errorf("instance not found")
	}
}

func testAccCheckS3DestroyBucket(n string) resource.TestCheckFunc {
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
		_, err = conn.DeleteBucket(&s3.DeleteBucketInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return fmt.Errorf("error destroying Bucket (%s) in testAccCheckS3DestroyBucket: %s", rs.Primary.ID, err)
		}
		return nil
	}
}

func testAccCheckS3BucketPolicy(n string, policy string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		config := common.TestAccProvider.Meta().(*cfg.Config)
		conn, err := config.S3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud s3 client: %s", err)
		}

		out, err := conn.GetBucketPolicy(&s3.GetBucketPolicyInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if policy == "" {
			if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "NoSuchBucketPolicy" {
				// expected
				return nil
			}
			if err == nil {
				return fmt.Errorf("expected no policy, got: %#v", *out.Policy)
			} else {
				return fmt.Errorf("getBucketPolicy error: %v, expected %s", err, policy)
			}
		}
		if err != nil {
			return fmt.Errorf("getBucketPolicy error: %v, expected %s", err, policy)
		}

		if v := out.Policy; v == nil {
			if policy != "" {
				return fmt.Errorf("bad policy, found nil, expected: %s", policy)
			}
		} else {
			expected := make(map[string]interface{})
			if err := json.Unmarshal([]byte(policy), &expected); err != nil {
				return err
			}
			actual := make(map[string]interface{})
			if err := json.Unmarshal([]byte(*v), &actual); err != nil {
				return err
			}

			if !reflect.DeepEqual(expected, actual) {
				return fmt.Errorf("bad policy, expected: %#v, got %#v", expected, actual)
			}
		}

		return nil
	}
}

func testAccCheckS3BucketWebsite(n string, indexDoc string, errorDoc string, redirectProtocol string, redirectTo string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.S3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud s3 client: %s", err)
		}

		out, err := client.GetBucketWebsite(&s3.GetBucketWebsiteInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			if indexDoc == "" {
				// If we want to assert that the website is not there, than
				// this error is expected
				return nil
			} else {
				return fmt.Errorf("S3BucketWebsite error: %v", err)
			}
		}

		if v := out.IndexDocument; v == nil {
			if indexDoc != "" {
				return fmt.Errorf("bad index doc, found nil, expected: %s", indexDoc)
			}
		} else {
			if *v.Suffix != indexDoc {
				return fmt.Errorf("bad index doc, expected: %s, got %#v", indexDoc, out.IndexDocument)
			}
		}

		if v := out.ErrorDocument; v == nil {
			if errorDoc != "" {
				return fmt.Errorf("bad error doc, found nil, expected: %s", errorDoc)
			}
		} else {
			if *v.Key != errorDoc {
				return fmt.Errorf("bad error doc, expected: %s, got %#v", errorDoc, out.ErrorDocument)
			}
		}

		if v := out.RedirectAllRequestsTo; v == nil {
			if redirectTo != "" {
				return fmt.Errorf("bad redirect to, found nil, expected: %s", redirectTo)
			}
		} else {
			if *v.HostName != redirectTo {
				return fmt.Errorf("bad redirect to, expected: %s, got %#v", redirectTo, out.RedirectAllRequestsTo)
			}
			if redirectProtocol != "" && v.Protocol != nil && *v.Protocol != redirectProtocol {
				return fmt.Errorf("bad redirect protocol to, expected: %s, got %#v", redirectProtocol, out.RedirectAllRequestsTo)
			}
		}

		return nil
	}
}

func testAccCheckS3BucketWebsiteRoutingRules(n string, routingRules []*s3.RoutingRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.S3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud s3 client: %s", err)
		}

		out, err := client.GetBucketWebsite(&s3.GetBucketWebsiteInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			if routingRules == nil {
				return nil
			}
			return fmt.Errorf("getBucketWebsite error: %v", err)
		}

		if !reflect.DeepEqual(out.RoutingRules, routingRules) {
			return fmt.Errorf("bad routing rule, expected: %v, got %v", routingRules, out.RoutingRules)
		}

		return nil
	}
}

func testAccCheckS3BucketVersioning(n string, versioningStatus string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.S3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud s3 client: %s", err)
		}

		out, err := client.GetBucketVersioning(&s3.GetBucketVersioningInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return fmt.Errorf("getBucketVersioning error: %v", err)
		}

		if v := out.Status; v == nil {
			if versioningStatus != "" {
				return fmt.Errorf("bad error versioning status, found nil, expected: %s", versioningStatus)
			}
		} else {
			if *v != versioningStatus {
				return fmt.Errorf("bad error versioning status, expected: %s, got %s", versioningStatus, *v)
			}
		}

		return nil
	}
}

func testAccCheckS3BucketCors(n string, corsRules []*s3.CORSRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.S3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud s3 client: %s", err)
		}

		out, err := client.GetBucketCors(&s3.GetBucketCorsInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return fmt.Errorf("getBucketCors error: %v", err)
		}

		if !reflect.DeepEqual(out.CORSRules, corsRules) {
			return fmt.Errorf("bad error cors rule, expected: %v, got %v", corsRules, out.CORSRules)
		}

		return nil
	}
}

func testAccCheckS3BucketLogging(n, b, p string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.S3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud s3 client: %s", err)
		}

		out, err := client.GetBucketLogging(&s3.GetBucketLoggingInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return fmt.Errorf("getBucketLogging error: %v", err)
		}

		tb := s.RootModule().Resources[b]

		if v := out.LoggingEnabled.TargetBucket; v == nil {
			if tb.Primary.ID != "" {
				return fmt.Errorf("bad target bucket, found nil, expected: %s", tb.Primary.ID)
			}
		} else {
			if *v != tb.Primary.ID {
				return fmt.Errorf("bad target bucket, expected: %s, got %s", tb.Primary.ID, *v)
			}
		}

		if v := out.LoggingEnabled.TargetPrefix; v == nil {
			if p != "" {
				return fmt.Errorf("bad target prefix, found nil, expected: %s", p)
			}
		} else {
			if *v != p {
				return fmt.Errorf("bad target prefix, expected: %s, got %s", p, *v)
			}
		}

		return nil
	}
}

// These need a bit of randomness as the name can only be used once globally
// within AWS
func testAccBucketName(randInt int) string {
	return fmt.Sprintf("tf-test-bucket-%d", randInt)
}

func testAccBucketDomainName(randInt int) string {
	return fmt.Sprintf("tf-test-bucket-%d.obs.eu-de.otc.t-systems.com", randInt)
}

func testAccWebsiteEndpoint(randInt int) string {
	return fmt.Sprintf("tf-test-bucket-%d.s3-website.%s.amazonaws.com", randInt, env.OS_REGION_NAME)
}

func testAccS3BucketPolicy(randInt int) string {
	return fmt.Sprintf(`{ "Version": "2008-10-17", "Statement": [ { "Effect": "Allow", "Principal": { "AWS": ["*"] }, "Action": ["s3:GetObject"], "Resource": ["arn:aws:s3:::tf-test-bucket-%d/*"] } ] }`, randInt)
}

func testAccS3BucketConfig(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "public-read"
}
`, randInt)
}

func testAccAWSS3MultiBucketConfigWithTags(randInt int) string {
	t := template.Must(template.New("t1").
		Parse(`
resource "opentelekomcloud_s3_bucket" "bucket1" {
	bucket = "tf-test-bucket-1-{{.GUID}}"
	acl = "private"
	force_destroy = true
	tags = {
		Name = "tf-test-bucket-1-{{.GUID}}"
		Environment = "{{.GUID}}"
	}
}

resource "opentelekomcloud_s3_bucket" "bucket2" {
	bucket = "tf-test-bucket-2-{{.GUID}}"
	acl = "private"
	force_destroy = true
	tags = {
		Name = "tf-test-bucket-2-{{.GUID}}"
		Environment = "{{.GUID}}"
	}
}

resource "opentelekomcloud_s3_bucket" "bucket3" {
	bucket = "tf-test-bucket-3-{{.GUID}}"
	acl = "private"
	force_destroy = true
	tags = {
		Name = "tf-test-bucket-3-{{.GUID}}"
		Environment = "{{.GUID}}"
	}
}

resource "opentelekomcloud_s3_bucket" "bucket4" {
	bucket = "tf-test-bucket-4-{{.GUID}}"
	acl = "private"
	force_destroy = true
	tags = {
		Name = "tf-test-bucket-4-{{.GUID}}"
		Environment = "{{.GUID}}"
	}
}

resource "opentelekomcloud_s3_bucket" "bucket5" {
	bucket = "tf-test-bucket-5-{{.GUID}}"
	acl = "private"
	force_destroy = true
	tags = {
		Name = "tf-test-bucket-5-{{.GUID}}"
		Environment = "{{.GUID}}"
	}
}

resource "opentelekomcloud_s3_bucket" "bucket6" {
	bucket = "tf-test-bucket-6-{{.GUID}}"
	acl = "private"
	force_destroy = true
	tags = {
		Name = "tf-test-bucket-6-{{.GUID}}"
		Environment = "{{.GUID}}"
	}
}
`))
	var doc bytes.Buffer
	_ = t.Execute(&doc, struct{ GUID int }{GUID: randInt})
	return doc.String()
}

func testAccS3BucketConfigWithRegion(randInt int) string {
	return fmt.Sprintf(`
provider "opentelekomcloud" {
  alias  = "reg1"
  region = "%s"
}

resource "opentelekomcloud_s3_bucket" "bucket" {
  provider = "reg1"
  bucket   = "tf-test-bucket-%d"
  region   = "%s"
}
`, env.OS_REGION_NAME, randInt, env.OS_REGION_NAME)
}

func testAccS3BucketWebsiteConfig(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "public-read"

  website {
    index_document = "index.html"
  }
}
`, randInt)
}

func testAccS3BucketWebsiteConfigWithError(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "public-read"

  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}
`, randInt)
}

func testAccS3BucketWebsiteConfigWithRedirect(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "public-read"

  website {
    redirect_all_requests_to = "hashicorp.com"
  }
}
`, randInt)
}

func testAccS3BucketWebsiteConfigWithHttpsRedirect(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "public-read"

  website {
    redirect_all_requests_to = "https://hashicorp.com"
  }
}
`, randInt)
}

func testAccS3BucketWebsiteConfigWithRoutingRules(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "public-read"

  website {
    index_document = "index.html"
    error_document = "error.html"
    routing_rules  = <<EOF
[{
	"Condition": {
		"KeyPrefixEquals": "docs/"
	},
	"Redirect": {
		"ReplaceKeyPrefixWith": "documents/"
	}
}]
EOF
  }
}
`, randInt)
}

func testAccS3BucketConfigWithPolicy(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "public-read"
  policy = %s
}
`, randInt, strconv.Quote(testAccS3BucketPolicy(randInt)))
}

func testAccS3BucketDestroyedConfig(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "public-read"
}
`, randInt)
}

func testAccS3BucketConfigWithEmptyPolicy(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "public-read"
  policy = ""
}
`, randInt)
}

func testAccS3BucketConfigWithVersioning(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "public-read"
  versioning {
    enabled = true
  }
}
`, randInt)
}

func testAccS3BucketConfigWithDisableVersioning(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "public-read"
  versioning {
    enabled = false
  }
}
`, randInt)
}

func testAccS3BucketConfigWithCORS(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "public-read"
  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = ["https://www.example.com"]
    expose_headers  = ["x-amz-server-side-encryption", "ETag"]
    max_age_seconds = 3000
  }
}
`, randInt)
}

func testAccS3BucketConfigWithLogging(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "log_bucket" {
  bucket        = "tf-test-log-bucket-%d"
  acl           = "log-delivery-write"
  force_destroy = "true"
}
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "private"
  logging {
    target_bucket = opentelekomcloud_s3_bucket.log_bucket.id
    target_prefix = "log/"
  }
}
`, randInt, randInt)
}

func testAccS3BucketConfigWithLifecycle(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "private"
  lifecycle_rule {
    id      = "id1"
    prefix  = "path1/"
    enabled = true

    expiration {
      days = 365
    }
  }
  lifecycle_rule {
    id      = "id2"
    prefix  = "path2/"
    enabled = true

    expiration {
      date = "2016-01-12"
    }
  }
  lifecycle_rule {
    id      = "id3"
    prefix  = "path3/"
    enabled = true

    expiration {
      days = "30"
    }
  }
  lifecycle_rule {
    id      = "id4"
    prefix  = "path4/"
    enabled = true

    expiration {
      date = "2016-01-12"
    }
  }
}
`, randInt)
}

func testAccS3BucketConfigWithVersioningLifecycle(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "private"
  versioning {
    enabled = false
  }
  lifecycle_rule {
    id      = "id1"
    prefix  = "path1/"
    enabled = true

    noncurrent_version_expiration {
      days = 365
    }
  }
  lifecycle_rule {
    id      = "id2"
    prefix  = "path2/"
    enabled = false

    noncurrent_version_expiration {
      days = 365
    }
  }
  lifecycle_rule {
    id      = "id3"
    prefix  = "path3/"
    enabled = true

    noncurrent_version_expiration {
      days = 30
    }
  }
}
`, randInt)
}

const (
	testAccS3BucketConfig_namePrefix = `
resource "opentelekomcloud_s3_bucket" "test" {
  bucket_prefix = "tf-test-"
}
`
	testAccS3BucketConfig_generatedName = `
resource "opentelekomcloud_s3_bucket" "test" {
  bucket_prefix = "tf-test-"
}
`

	testAccS3BucketConfigWithAcl = `
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl = "public-read"
}
`

	testAccS3BucketConfigWithAclUpdate = `
resource "opentelekomcloud_s3_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl = "private"
}
`
)
