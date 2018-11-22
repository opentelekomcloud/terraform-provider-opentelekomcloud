package opentelekomcloud

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/swauth"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform/helper/pathorcontents"
	"github.com/hashicorp/terraform/terraform"
	"github.com/huaweicloud/golangsdk"
	huaweisdk "github.com/huaweicloud/golangsdk/openstack"
)

type Config struct {
	AccessKey        string
	SecretKey        string
	CACertFile       string
	ClientCertFile   string
	ClientKeyFile    string
	DomainID         string
	DomainName       string
	EndpointType     string
	IdentityEndpoint string
	Insecure         bool
	Password         string
	Region           string
	Swauth           bool
	TenantID         string
	TenantName       string
	Token            string
	Username         string
	UserID           string

	OsClient *gophercloud.ProviderClient
	HwClient *golangsdk.ProviderClient
	s3sess   *session.Session
}

func (c *Config) LoadAndValidate() error {
	validEndpoint := false
	validEndpoints := []string{
		"internal", "internalURL",
		"admin", "adminURL",
		"public", "publicURL",
		"",
	}

	for _, endpoint := range validEndpoints {
		if c.EndpointType == endpoint {
			validEndpoint = true
		}
	}

	if !validEndpoint {
		return fmt.Errorf("Invalid endpoint type provided")
	}

	config, err := generateTLSConfig(c)
	if err != nil {
		return err
	}
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config}

	// if OS_DEBUG is set, log the requests and responses
	var osDebug bool
	if os.Getenv("OS_DEBUG") != "" {
		osDebug = true
	}

	err = c.newopenstackClient(transport, osDebug)
	if err != nil {
		return err
	}

	err = c.newhwClient(transport, osDebug)
	if err != nil {
		return err
	}

	return c.newS3Session(osDebug)
}

func generateTLSConfig(c *Config) (*tls.Config, error) {
	config := &tls.Config{}
	if c.CACertFile != "" {
		caCert, _, err := pathorcontents.Read(c.CACertFile)
		if err != nil {
			return nil, fmt.Errorf("Error reading CA Cert: %s", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))
		config.RootCAs = caCertPool
	}

	if c.Insecure {
		config.InsecureSkipVerify = true
	}

	if c.ClientCertFile != "" && c.ClientKeyFile != "" {
		clientCert, _, err := pathorcontents.Read(c.ClientCertFile)
		if err != nil {
			return nil, fmt.Errorf("Error reading Client Cert: %s", err)
		}
		clientKey, _, err := pathorcontents.Read(c.ClientKeyFile)
		if err != nil {
			return nil, fmt.Errorf("Error reading Client Key: %s", err)
		}

		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil {
			return nil, err
		}

		config.Certificates = []tls.Certificate{cert}
		config.BuildNameToCertificate()
	}

	return config, nil
}

func (c *Config) newS3Session(osDebug bool) error {
	// Don't get AWS session unless we need it for Accesskey, SecretKey.
	if c.AccessKey != "" && c.SecretKey != "" {
		// Setup AWS/S3 client/config information for Swift S3 buckets
		log.Println("[INFO] Building Swift S3 auth structure")
		creds, err := GetCredentials(c)
		if err != nil {
			return err
		}
		// Call Get to check for credential provider. If nothing found, we'll get an
		// error, and we can present it nicely to the user
		cp, err := creds.Get()
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "NoCredentialProviders" {
				return fmt.Errorf(`No valid credential sources found for Swift S3 Provider.
  Please see https://terraform.io/docs/providers/aws/index.html for more information on
  providing credentials for the S3 Provider`)
			}

			return fmt.Errorf("Error loading credentials for Swift S3 Provider: %s", err)
		}

		log.Printf("[INFO] Swift S3 Auth provider used: %q", cp.ProviderName)

		awsConfig := &aws.Config{
			Credentials: creds,
			Region:      aws.String(c.Region),
			//MaxRetries:       aws.Int(c.MaxRetries),
			HTTPClient: cleanhttp.DefaultClient(),
			//S3ForcePathStyle: aws.Bool(c.S3ForcePathStyle),
		}

		if osDebug {
			awsConfig.LogLevel = aws.LogLevel(aws.LogDebugWithHTTPBody | aws.LogDebugWithRequestRetries | aws.LogDebugWithRequestErrors)
			awsConfig.Logger = awsLogger{}
		}

		if c.Insecure {
			transport := awsConfig.HTTPClient.Transport.(*http.Transport)
			transport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		// Set up base session for AWS/Swift S3
		c.s3sess, err = session.NewSession(awsConfig)
		if err != nil {
			return errwrap.Wrapf("Error creating Swift S3 session: {{err}}", err)
		}
	}
	return nil
}

func (c *Config) newopenstackClient(transport *http.Transport, osDebug bool) error {
	ao := gophercloud.AuthOptions{
		DomainID:         c.DomainID,
		DomainName:       c.DomainName,
		IdentityEndpoint: c.IdentityEndpoint,
		Password:         c.Password,
		TenantID:         c.TenantID,
		TenantName:       c.TenantName,
		TokenID:          c.Token,
		Username:         c.Username,
		UserID:           c.UserID,
	}

	client, err := openstack.NewClient(ao.IdentityEndpoint)
	if err != nil {
		return err
	}

	// Set UserAgent
	client.UserAgent.Prepend(terraform.UserAgentString())

	client.HTTPClient = http.Client{
		Transport: &LogRoundTripper{
			Rt:      transport,
			OsDebug: osDebug,
		},
	}

	// If using Swift Authentication, there's no need to validate authentication normally.
	if !c.Swauth {
		err = openstack.Authenticate(client, ao)
		if err != nil {
			return err
		}
	}

	c.OsClient = client
	return nil
}

func (c *Config) newhwClient(transport *http.Transport, osDebug bool) error {
	ao := golangsdk.AuthOptions{
		DomainID:         c.DomainID,
		DomainName:       c.DomainName,
		IdentityEndpoint: c.IdentityEndpoint,
		Password:         c.Password,
		TenantID:         c.TenantID,
		TenantName:       c.TenantName,
		TokenID:          c.Token,
		Username:         c.Username,
		UserID:           c.UserID,
	}

	client, err := huaweisdk.NewClient(ao.IdentityEndpoint)
	if err != nil {
		return err
	}

	// Set UserAgent
	client.UserAgent.Prepend(terraform.UserAgentString())

	client.HTTPClient = http.Client{
		Transport: &LogRoundTripper{
			Rt:      transport,
			OsDebug: osDebug,
		},
	}

	// If using Swift Authentication, there's no need to validate authentication normally.
	if !c.Swauth {
		err = huaweisdk.Authenticate(client, ao)
		if err != nil {
			return err
		}
	}

	c.HwClient = client
	return nil
}

type awsLogger struct{}

func (l awsLogger) Log(args ...interface{}) {
	tokens := make([]string, 0, len(args))
	for _, arg := range args {
		if token, ok := arg.(string); ok {
			tokens = append(tokens, token)
		}
	}
	log.Printf("[DEBUG] [aws-sdk-go] %s", strings.Join(tokens, " "))
}

func (c *Config) determineRegion(region string) string {
	// If a resource-level region was not specified, and a provider-level region was set,
	// use the provider-level region.
	if region == "" && c.Region != "" {
		region = c.Region
	}

	log.Printf("[DEBUG] OpenTelekomCloud Region is: %s", region)
	return region
}

func (c *Config) computeS3conn(region string) (*s3.S3, error) {
	if c.s3sess == nil {
		return nil, fmt.Errorf("Missing credentials for Swift S3 Provider, need access_key and secret_key values for provider.")
	}

	client, err := openstack.NewImageServiceV2(c.OsClient, gophercloud.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getEndpointType(),
	})
	// Bit of a hack, seems the only way to compute this.
	endpoint := strings.Replace(client.Endpoint, "//ims", "//obs", 1)

	awsS3Sess := c.s3sess.Copy(&aws.Config{Endpoint: aws.String(endpoint)})
	s3conn := s3.New(awsS3Sess)

	return s3conn, err
}

func (c *Config) blockStorageV1Client(region string) (*gophercloud.ServiceClient, error) {
	return openstack.NewBlockStorageV1(c.OsClient, gophercloud.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getEndpointType(),
	})
}

func (c *Config) blockStorageV2Client(region string) (*gophercloud.ServiceClient, error) {
	return openstack.NewBlockStorageV2(c.OsClient, gophercloud.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getEndpointType(),
	})
}

func (c *Config) computeV2Client(region string) (*gophercloud.ServiceClient, error) {
	return openstack.NewComputeV2(c.OsClient, gophercloud.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getEndpointType(),
	})
}

func (c *Config) dnsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewDNSV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) identityV3Client(region string) (*gophercloud.ServiceClient, error) {
	return openstack.NewIdentityV3(c.OsClient, gophercloud.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getEndpointType(),
	})
}

func (c *Config) imageV2Client(region string) (*gophercloud.ServiceClient, error) {
	return openstack.NewImageServiceV2(c.OsClient, gophercloud.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getEndpointType(),
	})
}

func (c *Config) networkingV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewNetworkV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) networkingV2Client(region string) (*gophercloud.ServiceClient, error) {
	return openstack.NewNetworkV2(c.OsClient, gophercloud.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getEndpointType(),
	})
}

func (c *Config) objectStorageV1Client(region string) (*gophercloud.ServiceClient, error) {
	// If Swift Authentication is being used, return a swauth client.
	if c.Swauth {
		return swauth.NewObjectStorageV1(c.OsClient, swauth.AuthOpts{
			User: c.Username,
			Key:  c.Password,
		})
	}

	return openstack.NewObjectStorageV1(c.OsClient, gophercloud.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getEndpointType(),
	})
}

func (c *Config) SmnV2Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewSmnServiceV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) getEndpointType() gophercloud.Availability {
	if c.EndpointType == "internal" || c.EndpointType == "internalURL" {
		return gophercloud.AvailabilityInternal
	}
	if c.EndpointType == "admin" || c.EndpointType == "adminURL" {
		return gophercloud.AvailabilityAdmin
	}
	return gophercloud.AvailabilityPublic
}

func (c *Config) loadCESClient(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewCESClient(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) getHwEndpointType() golangsdk.Availability {
	if c.EndpointType == "internal" || c.EndpointType == "internalURL" {
		return golangsdk.AvailabilityInternal
	}
	if c.EndpointType == "admin" || c.EndpointType == "adminURL" {
		return golangsdk.AvailabilityAdmin
	}
	return golangsdk.AvailabilityPublic
}

func (c *Config) loadECSV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewComputeV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) kmsKeyV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewKmsKeyV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) hwNetworkV2Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewNetworkV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) loadEVSV2Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewBlockStorageV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) natV2Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewNatV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) orchestrationV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewOrchestrationV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) sfsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewSharedFileSystemV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) vbsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewVBS(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

//computeV2HWClient used to access the v2 bms Services i.e. flavor, nic, keypair.
func (c *Config) computeV2HWClient(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewComputeV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

//bmsClient used to access the v2.1 bms Services i.e. servers, tags.
func (c *Config) bmsClient(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewBMSV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) autoscalingV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewAutoScalingService(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) csbsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewCSBSService(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) dehV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewDeHServiceV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) dmsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewDMSServiceV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) MrsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewMapReduceV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) elbV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewELBV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) rdsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewRDSV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}

func (c *Config) cceV3Client(region string) (*golangsdk.ServiceClient, error) {
	return huaweisdk.NewCCEV3(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getHwEndpointType(),
	})
}
