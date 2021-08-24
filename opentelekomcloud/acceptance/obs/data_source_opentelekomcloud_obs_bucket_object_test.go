package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/obs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccDataSourceObsBucketObject_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceOnlyConf, conf := testAccDataSourceObsObjectConfigBasic(rInt)

	var dsObj obs.GetObjectOutput

	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { common.TestAccPreCheck(t) },
		ProviderFactories:         common.TestAccProviderFactories,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: resourceOnlyConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketObjectExists("opentelekomcloud_obs_bucket_object.object"),
				),
			},
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsObjectDataSourceExists("data.opentelekomcloud_obs_bucket_object.obj", &dsObj),
					resource.TestCheckResourceAttr("data.opentelekomcloud_obs_bucket_object.obj", "content_length", "11"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_obs_bucket_object.obj", "content_type", "binary/octet-stream"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_obs_bucket_object.obj", "etag", "b10a8db164e0754105b7a99be72e3fe5"),
					resource.TestMatchResourceAttr("data.opentelekomcloud_obs_bucket_object.obj", "last_modified",
						regexp.MustCompile("^[a-zA-Z]{3}, [0-9]+ [a-zA-Z]+ [0-9]{4} [0-9:]+ [A-Z]+$")),
					resource.TestCheckNoResourceAttr("data.opentelekomcloud_obs_bucket_object.obj", "body"),
				),
			},
		},
	})
}

func TestAccDataSourceObsBucketObject_readableBody(t *testing.T) {
	rInt := acctest.RandInt()
	resourceOnlyConf, conf := testAccDataSourceObsObjectConfigReadableBody(rInt)

	var dsObj obs.GetObjectOutput

	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { common.TestAccPreCheck(t) },
		ProviderFactories:         common.TestAccProviderFactories,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: resourceOnlyConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketObjectExists("opentelekomcloud_obs_bucket_object.object"),
				),
			},
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsObjectDataSourceExists("data.opentelekomcloud_obs_bucket_object.obj", &dsObj),
					resource.TestCheckResourceAttr("data.opentelekomcloud_obs_bucket_object.obj", "content_length", "3"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_obs_bucket_object.obj", "content_type", "text/plain"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_obs_bucket_object.obj", "etag", "a6105c0a611b41b08f1209506350279e"),
					resource.TestMatchResourceAttr("data.opentelekomcloud_obs_bucket_object.obj", "last_modified",
						regexp.MustCompile("^[a-zA-Z]{3}, [0-9]+ [a-zA-Z]+ [0-9]{4} [0-9:]+ [A-Z]+$")),
					resource.TestCheckResourceAttr("data.opentelekomcloud_obs_bucket_object.obj", "body", "yes"),
				),
			},
		},
	})
}

func TestAccDataSourceObsBucketObject_allParams(t *testing.T) {
	t.Skip("Removing versioned bucket is broken, see GH-779")

	rInt := acctest.RandInt()
	resourceOnlyConf, conf := testAccDataSourceObsObjectConfigAllParams(rInt)

	var dsObj obs.GetObjectOutput

	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { common.TestAccPreCheck(t) },
		ProviderFactories:         common.TestAccProviderFactories,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: resourceOnlyConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketObjectExists("opentelekomcloud_obs_bucket_object.object"),
				),
			},
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsObjectDataSourceExists("data.opentelekomcloud_obs_bucket_object.obj", &dsObj),
					resource.TestCheckResourceAttr("data.opentelekomcloud_obs_bucket_object.obj", "content_length", "21"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_obs_bucket_object.obj", "content_type", "application/unknown"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_obs_bucket_object.obj", "etag", "723f7a6ac0c57b445790914668f98640"),
					resource.TestMatchResourceAttr("data.opentelekomcloud_obs_bucket_object.obj", "last_modified",
						regexp.MustCompile("^[a-zA-Z]{3}, [0-9]+ [a-zA-Z]+ [0-9]{4} [0-9:]+ [A-Z]+$")),
					resource.TestCheckNoResourceAttr("data.opentelekomcloud_obs_bucket_object.obj", "body"),
					// Encryption is off
					resource.TestCheckResourceAttr("data.opentelekomcloud_obs_bucket_object.obj", "expiration", ""),
					// Currently unsupported in opentelekomcloud_obs_bucket_object resource
					resource.TestCheckResourceAttr("data.opentelekomcloud_obs_bucket_object.obj", "expires", ""),
					resource.TestCheckResourceAttr("data.opentelekomcloud_obs_bucket_object.obj", "website_redirect_location", ""),
					resource.TestCheckResourceAttr("data.opentelekomcloud_obs_bucket_object.obj", "metadata.%", "0"),
				),
			},
		},
	})
}

func testAccCheckObsObjectDataSourceExists(n string, obj *obs.GetObjectOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find OBS object data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("OBS object data source ID not set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.NewObjectStorageClient(env.OsRegionName)
		if err != nil {
			return fmt.Errorf("error creating OBS client: %s", err)
		}
		out, err := client.GetObject(
			&obs.GetObjectInput{
				GetObjectMetadataInput: obs.GetObjectMetadataInput{
					Bucket: rs.Primary.Attributes["bucket"],
					Key:    rs.Primary.Attributes["key"],
				},
			})
		if err != nil {
			return fmt.Errorf("failed getting S3 Object from %s: %s",
				rs.Primary.Attributes["bucket"]+"/"+rs.Primary.Attributes["key"], err)
		}

		*obj = *out

		return nil
	}
}

func testAccDataSourceObsObjectConfigBasic(randInt int) (string, string) {
	resources := fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "object_bucket" {
	bucket = "tf-object-test-bucket-%d"
}
resource "opentelekomcloud_obs_bucket_object" "object" {
	bucket = opentelekomcloud_obs_bucket.object_bucket.bucket
	key = "tf-testing-obj-%d"
	content = "Hello World"
}
`, randInt, randInt)

	both := fmt.Sprintf(`%s
data "opentelekomcloud_obs_bucket_object" "obj" {
	bucket = "tf-object-test-bucket-%d"
	key = "tf-testing-obj-%d"
}`, resources, randInt, randInt)

	return resources, both
}

func testAccDataSourceObsObjectConfigReadableBody(randInt int) (string, string) {
	resources := fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "object_bucket" {
	bucket = "tf-object-test-bucket-%d"
}
resource "opentelekomcloud_obs_bucket_object" "object" {
	bucket = opentelekomcloud_obs_bucket.object_bucket.bucket
	key = "tf-testing-obj-%d-readable"
	content = "yes"
	content_type = "text/plain"
}
`, randInt, randInt)

	both := fmt.Sprintf(`%s
data "opentelekomcloud_obs_bucket_object" "obj" {
	bucket = "tf-object-test-bucket-%d"
	key = "tf-testing-obj-%d-readable"
}`, resources, randInt, randInt)

	return resources, both
}

func testAccDataSourceObsObjectConfigAllParams(randInt int) (string, string) { // nolint:unused
	resources := fmt.Sprintf(`
resource "opentelekomcloud_obs_bucket" "object_bucket" {
	bucket = "tf-object-test-bucket-%[1]d"
	versioning = true
}

resource "opentelekomcloud_obs_bucket_object" "object" {
	bucket = opentelekomcloud_obs_bucket.object_bucket.bucket
	key = "tf-testing-obj-%[1]d-all-params"
	content = <<CONTENT
{"msg": "Hi there!"}
CONTENT
	content_type = "application/unknown"
}
`, randInt)

	both := fmt.Sprintf(`%s
data "opentelekomcloud_obs_bucket_object" "obj" {
	bucket = "tf-object-test-bucket-%[2]d"
	key = "tf-testing-obj-%[2]d-all-params"
}`, resources, randInt)

	return resources, both
}
