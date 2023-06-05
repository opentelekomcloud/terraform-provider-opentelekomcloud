package acceptance

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	s3s "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/s3"
)

func TestAccS3BucketObject_source(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "tf-acc-s3-obj-source")
	th.AssertNoErr(t, err)
	t.Cleanup(func() {
		th.AssertNoErr(t, os.Remove(tmpFile.Name()))
	})

	rInt := acctest.RandInt()
	// first write some data to the tempfile just so it's not 0 bytes.
	err = ioutil.WriteFile(tmpFile.Name(), []byte("{anything will do }"), 0644)
	th.AssertNoErr(t, err)
	var obj s3.GetObjectOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketObjectConfigSource(rInt, tmpFile.Name()),
				Check:  testAccCheckS3BucketObjectExists("opentelekomcloud_s3_bucket_object.object", &obj),
			},
		},
	})
}

func TestAccS3BucketObject_content(t *testing.T) {
	rInt := acctest.RandInt()
	var obj s3.GetObjectOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccS3BucketObjectConfigContent(rInt),
				Check:     testAccCheckS3BucketObjectExists("opentelekomcloud_s3_bucket_object.object", &obj),
			},
		},
	})
}

func TestAccS3BucketObject_withContentCharacteristics(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "tf-acc-s3-obj-content-characteristics")
	th.AssertNoErr(t, err)
	t.Cleanup(func() {
		th.AssertNoErr(t, os.Remove(tmpFile.Name()))
	})

	rInt := acctest.RandInt()
	// first write some data to the tempfile just so it's not 0 bytes.
	err = ioutil.WriteFile(tmpFile.Name(), []byte("{anything will do }"), 0644)
	th.AssertNoErr(t, err)

	var obj s3.GetObjectOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketObjectConfigWithContentCharacteristics(rInt, tmpFile.Name()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketObjectExists("opentelekomcloud_s3_bucket_object.object", &obj),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket_object.object", "content_type", "binary/octet-stream"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket_object.object", "website_redirect", "http://google.com"),
				),
			},
		},
	})
}

func TestAccS3BucketObject_updates(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "tf-acc-s3-obj-updates")
	th.AssertNoErr(t, err)
	t.Cleanup(func() {
		th.AssertNoErr(t, os.Remove(tmpFile.Name()))
	})

	rInt := acctest.RandInt()
	err = ioutil.WriteFile(tmpFile.Name(), []byte("initial object state"), 0644)
	th.AssertNoErr(t, err)
	var obj s3.GetObjectOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketObjectConfigUpdates(rInt, tmpFile.Name()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketObjectExists("opentelekomcloud_s3_bucket_object.object", &obj),
					resource.TestCheckResourceAttr("opentelekomcloud_s3_bucket_object.object", "etag", "647d1d58e1011c743ec67d5e8af87b53"),
				),
			},
			{
				PreConfig: func() {
					err = ioutil.WriteFile(tmpFile.Name(), []byte("modified object"), 0644)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testAccS3BucketObjectConfigUpdates(rInt, tmpFile.Name()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketObjectExists("opentelekomcloud_s3_bucket_object.object", &obj),
					resource.TestCheckResourceAttr("opentelekomcloud_s3_bucket_object.object", "etag", "1c7fd13df1515c2a13ad9eb068931f09"),
				),
			},
		},
	})
}

func TestAccS3BucketObject_updatesWithVersioning(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "tf-acc-s3-obj-updates-w-versions")
	th.AssertNoErr(t, err)
	t.Cleanup(func() {
		th.AssertNoErr(t, os.Remove(tmpFile.Name()))
	})

	rInt := acctest.RandInt()
	err = ioutil.WriteFile(tmpFile.Name(), []byte("initial versioned object state"), 0644)
	th.AssertNoErr(t, err)

	var originalObj, modifiedObj s3.GetObjectOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketObjectConfigUpdatesWithVersioning(rInt, tmpFile.Name()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketObjectExists("opentelekomcloud_s3_bucket_object.object", &originalObj),
					resource.TestCheckResourceAttr("opentelekomcloud_s3_bucket_object.object", "etag", "cee4407fa91906284e2a5e5e03e86b1b"),
				),
			},
			{
				PreConfig: func() {
					err = ioutil.WriteFile(tmpFile.Name(), []byte("modified versioned object"), 0644)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testAccS3BucketObjectConfigUpdatesWithVersioning(rInt, tmpFile.Name()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketObjectExists("opentelekomcloud_s3_bucket_object.object", &modifiedObj),
					resource.TestCheckResourceAttr("opentelekomcloud_s3_bucket_object.object", "etag", "00b8c73b1b50e7cc932362c7225b8e29"),
					testAccCheckS3BucketObjectVersionIdDiffers(&originalObj, &modifiedObj),
				),
			},
		},
	})
}

func testAccCheckS3BucketObjectVersionIdDiffers(first, second *s3.GetObjectOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if first.VersionId == nil {
			return fmt.Errorf("expected first object to have VersionId: %s", first)
		}
		if second.VersionId == nil {
			return fmt.Errorf("expected second object to have VersionId: %s", second)
		}

		if *first.VersionId == *second.VersionId {
			return fmt.Errorf("expected Version IDs to differ, but they are equal (%s)", *first.VersionId)
		}

		return nil
	}
}

func testAccCheckS3BucketObjectDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	s3conn, err := config.S3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud s3 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_s3_bucket_object" {
			continue
		}

		_, err := s3conn.HeadObject(
			&s3.HeadObjectInput{
				Bucket:  aws.String(rs.Primary.Attributes["bucket"]),
				Key:     aws.String(rs.Primary.Attributes["key"]),
				IfMatch: aws.String(rs.Primary.Attributes["etag"]),
			})
		if err == nil {
			return fmt.Errorf("swift S3 Object still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckS3BucketObjectExists(n string, obj *s3.GetObjectOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no S3 Bucket Object ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		s3conn, err := config.S3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud s3 client: %s", err)
		}
		out, err := s3conn.GetObject(
			&s3.GetObjectInput{
				Bucket:  aws.String(rs.Primary.Attributes["bucket"]),
				Key:     aws.String(rs.Primary.Attributes["key"]),
				IfMatch: aws.String(rs.Primary.Attributes["etag"]),
			})
		if err != nil {
			return fmt.Errorf("S3Bucket Object error: %s", err)
		}

		*obj = *out

		return nil
	}
}

func TestAccS3BucketObject_sse(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "tf-acc-s3-obj-source-sse")
	th.AssertNoErr(t, err)
	t.Cleanup(func() {
		th.AssertNoErr(t, os.Remove(tmpFile.Name()))
	})

	// first write some data to the tempfile just so it's not 0 bytes.
	err = ioutil.WriteFile(tmpFile.Name(), []byte("{anything will do}"), 0644)
	th.AssertNoErr(t, err)

	rInt := acctest.RandInt()
	var obj s3.GetObjectOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccS3BucketObjectConfigWithSSE(rInt, tmpFile.Name()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketObjectExists(
						"opentelekomcloud_s3_bucket_object.object",
						&obj),
					testAccCheckS3BucketObjectSSE(
						"opentelekomcloud_s3_bucket_object.object",
						"aws:kms"),
				),
			},
		},
	})
}

func TestAccS3BucketObject_acl(t *testing.T) {
	rInt := acctest.RandInt()
	var obj s3.GetObjectOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckS3(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckS3BucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketObjectConfigAcl(rInt, "private"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketObjectExists(
						"opentelekomcloud_s3_bucket_object.object", &obj),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket_object.object",
						"acl",
						"private"),
					testAccCheckS3BucketObjectAcl(
						"opentelekomcloud_s3_bucket_object.object",
						[]string{"FULL_CONTROL"}),
				),
			},
			{
				Config: testAccS3BucketObjectConfigAcl(rInt, "public-read"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketObjectExists(
						"opentelekomcloud_s3_bucket_object.object",
						&obj),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_s3_bucket_object.object",
						"acl",
						"public-read"),
					testAccCheckS3BucketObjectAcl(
						"opentelekomcloud_s3_bucket_object.object",
						[]string{"FULL_CONTROL", "READ"}),
				),
			},
		},
	})
}

func testAccCheckS3BucketObjectAcl(n string, expectedPerms []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		config := common.TestAccProvider.Meta().(*cfg.Config)
		s3conn, err := config.S3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud s3 client: %s", err)
		}

		out, err := s3conn.GetObjectAcl(&s3.GetObjectAclInput{
			Bucket: aws.String(rs.Primary.Attributes["bucket"]),
			Key:    aws.String(rs.Primary.Attributes["key"]),
		})

		if err != nil {
			return fmt.Errorf("getObjectAcl error: %v", err)
		}

		var perms []string
		for _, v := range out.Grants {
			perms = append(perms, *v.Permission)
		}
		sort.Strings(perms)

		if !reflect.DeepEqual(perms, expectedPerms) {
			return fmt.Errorf("expected ACL permissions to be %v, got %v", expectedPerms, perms)
		}

		return nil
	}
}

func TestResourceS3BucketObjectAcl_validation(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "incorrect",
			ErrCount: 1,
		},
		{
			Value:    "public-read",
			ErrCount: 0,
		},
		{
			Value:    "public-read-write",
			ErrCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Value, func(t *testing.T) {
			tx := tc
			t.Parallel()
			_, errors := s3s.ValidateS3BucketObjectAclType(tx.Value, "acl")
			if len(errors) != tx.ErrCount {
				t.Fatalf("Expected to trigger %d validation errors, but got %d", tx.ErrCount, len(errors))
			}
		})
	}
}

func TestResourceS3BucketObjectStorageClass_validation(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "incorrect",
			ErrCount: 1,
		},
		{
			Value:    "STANDARD",
			ErrCount: 0,
		},
		{
			Value:    "REDUCED_REDUNDANCY",
			ErrCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Value, func(t *testing.T) {
			tx := tc
			t.Parallel()
			_, errors := validateS3BucketObjectStorageClassType(tx.Value, "storage_class")
			if len(errors) != tx.ErrCount {
				t.Fatalf("Expected not to trigger a validation error")
			}
		})
	}
}

func testAccCheckS3BucketObjectSSE(n, expectedSSE string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		config := common.TestAccProvider.Meta().(*cfg.Config)
		s3conn, err := config.S3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud s3 client: %s", err)
		}

		out, err := s3conn.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(rs.Primary.Attributes["bucket"]),
			Key:    aws.String(rs.Primary.Attributes["key"]),
		})

		if err != nil {
			return fmt.Errorf("headObject error: %v", err)
		}

		if out.ServerSideEncryption == nil {
			return fmt.Errorf("expected a non %v Server Side Encryption.", out.ServerSideEncryption)
		}

		sse := *out.ServerSideEncryption
		if sse != expectedSSE {
			return fmt.Errorf("expected Server Side Encryption %v, got %v.",
				expectedSSE, sse)
		}

		return nil
	}
}

func testAccS3BucketObjectConfigSource(randInt int, source string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "object_bucket" {
  bucket = "tf-object-test-bucket-%d"
}
resource "opentelekomcloud_s3_bucket_object" "object" {
  bucket       = opentelekomcloud_s3_bucket.object_bucket.bucket
  key          = "test-key"
  source       = "%s"
  content_type = "binary/octet-stream"
}
`, randInt, source)
}

func testAccS3BucketObjectConfigWithContentCharacteristics(randInt int, source string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "object_bucket_2" {
  bucket = "tf-object-test-bucket-%d"
}

resource "opentelekomcloud_s3_bucket_object" "object" {
  bucket           = opentelekomcloud_s3_bucket.object_bucket_2.bucket
  key              = "test-key"
  source           = "%s"
  content_language = "en"
  content_type     = "binary/octet-stream"
  website_redirect = "http://google.com"
}
`, randInt, source)
}

func testAccS3BucketObjectConfigContent(randInt int) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "object_bucket" {
  bucket = "tf-object-test-bucket-%d"
}
resource "opentelekomcloud_s3_bucket_object" "object" {
  bucket  = opentelekomcloud_s3_bucket.object_bucket.bucket
  key     = "test-key"
  content = "some_bucket_content"
}
`, randInt)
}

func testAccS3BucketObjectConfigUpdates(randInt int, source string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "object_bucket_3" {
  bucket = "tf-object-test-bucket-%d"
}

resource "opentelekomcloud_s3_bucket_object" "object" {
  bucket = opentelekomcloud_s3_bucket.object_bucket_3.bucket
  key    = "updateable-key"
  source = "%s"
  etag   = md5(file("%s"))
}
`, randInt, source, source)
}

func testAccS3BucketObjectConfigUpdatesWithVersioning(randInt int, source string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "object_bucket_3" {
  bucket = "tf-object-test-bucket-%d"
  versioning {
    enabled = true
  }
}

resource "opentelekomcloud_s3_bucket_object" "object" {
  bucket = opentelekomcloud_s3_bucket.object_bucket_3.bucket
  key    = "updateable-key"
  source = "%s"
  etag   = md5(file("%s"))
}
`, randInt, source, source)
}

func testAccS3BucketObjectConfigWithSSE(randInt int, source string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "object_bucket" {
  bucket = "tf-object-test-bucket-%d"
}

resource "opentelekomcloud_s3_bucket_object" "object" {
  bucket                 = opentelekomcloud_s3_bucket.object_bucket.bucket
  key                    = "test-key"
  source                 = "%s"
  server_side_encryption = "aws:kms"
}
`, randInt, source)
}

func testAccS3BucketObjectConfigAcl(randInt int, acl string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_s3_bucket" "object_bucket" {
  bucket = "tf-object-test-bucket-%d"
}
resource "opentelekomcloud_s3_bucket_object" "object" {
  bucket  = opentelekomcloud_s3_bucket.object_bucket.bucket
  key     = "test-key"
  content = "some_bucket_content"
  acl     = "%s"
}
`, randInt, acl)
}

func validateS3BucketObjectStorageClassType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	storageClass := map[string]bool{
		s3.StorageClassStandard:          true,
		s3.StorageClassReducedRedundancy: true,
		s3.StorageClassStandardIa:        true,
	}

	if _, ok := storageClass[value]; !ok {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid Storage Class type %q. Valid types are either %q, %q, or %q",
			k, value, s3.StorageClassStandard, s3.StorageClassReducedRedundancy,
			s3.StorageClassStandardIa))
	}
	return
}
