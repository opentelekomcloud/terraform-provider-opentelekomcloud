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
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/helper/pathorcontents"
	"github.com/hashicorp/terraform-plugin-sdk/httpclient"
	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/credentials"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/objectstorage/v1/swauth"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/obs"
)

const (
	serviceProjectLevel string = "project"
	serviceDomainLevel  string = "domain"
)

type Config struct {
	AccessKey        string
	SecretKey        string
	CACertFile       string
	ClientCertFile   string
	ClientKeyFile    string
	Cloud            string
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
	SecurityToken    string
	Username         string
	UserID           string
	AgencyName       string
	AgencyDomainName string
	DelegatedProject string
	MaxRetries       int
	terraformVersion string

	HwClient *golangsdk.ProviderClient
	s3sess   *session.Session

	DomainClient *golangsdk.ProviderClient
}

func (c *Config) LoadAndValidate() error {
	if c.MaxRetries < 0 {
		return fmt.Errorf("max_retries should be a positive value")
	}

	if c.IdentityEndpoint == "" && c.Cloud == "" {
		return fmt.Errorf("one of 'auth_url' or 'cloud' must be specified")
	}

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
		return fmt.Errorf("invalid endpoint type provided")
	}

	if c.Cloud != "" {
		err := readCloudsYaml(c)
		if err != nil {
			return err
		}
	}

	err := fmt.Errorf("must config token or aksk or username password to be authorized")

	if c.Token != "" {
		err = buildClientByToken(c)

	} else if c.AccessKey != "" && c.SecretKey != "" {
		err = buildClientByAKSK(c)

	} else if c.Password != "" && (c.Username != "" || c.UserID != "") {
		err = buildClientByPassword(c)
	}
	if err != nil {
		return err
	}

	var osDebug bool
	if os.Getenv("OS_DEBUG") != "" {
		osDebug = true
	}
	return c.newS3Session(osDebug)
}

func readCloudsYaml(c *Config) error {
	env := openstack.NewEnv("OS")
	cloud, err := env.Cloud()
	if err != nil {
		return err
	}

	// Auth data
	c.TenantName = cloud.AuthInfo.ProjectName
	c.TenantID = cloud.AuthInfo.ProjectID
	c.DomainName = cloud.AuthInfo.DomainName
	if c.DomainName == "" {
		c.DomainName = cloud.AuthInfo.ProjectDomainName
	}
	c.DomainID = cloud.AuthInfo.DomainID
	if c.DomainID == "" {
		c.DomainID = cloud.AuthInfo.ProjectDomainID
	}
	c.IdentityEndpoint = cloud.AuthInfo.AuthURL
	c.Token = cloud.AuthInfo.Token
	c.Username = cloud.AuthInfo.Username
	c.UserID = cloud.AuthInfo.UserID
	c.Password = cloud.AuthInfo.Password

	// General cloud info
	if c.Region == "" && cloud.RegionName != "" {
		c.Region = cloud.RegionName
	}
	if c.CACertFile == "" && cloud.CACertFile != "" {
		c.CACertFile = cloud.CACertFile
	}
	if c.ClientCertFile == "" && cloud.ClientCertFile != "" {
		c.ClientCertFile = cloud.ClientCertFile
	}
	if c.ClientKeyFile == "" && cloud.ClientKeyFile != "" {
		c.ClientKeyFile = cloud.ClientKeyFile
	}
	if cloud.Verify != nil {
		c.Insecure = !*cloud.Verify
	}
	return nil
}

func generateTLSConfig(c *Config) (*tls.Config, error) {
	config := &tls.Config{}
	if c.CACertFile != "" {
		caCert, _, err := pathorcontents.Read(c.CACertFile)
		if err != nil {
			return nil, fmt.Errorf("error reading CA Cert: %s", err)
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
			return nil, fmt.Errorf("error reading Client Cert: %s", err)
		}
		clientKey, _, err := pathorcontents.Read(c.ClientKeyFile)
		if err != nil {
			return nil, fmt.Errorf("error reading Client Key: %s", err)
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

			return fmt.Errorf("error loading credentials for Swift S3 Provider: %s", err)
		}

		log.Printf("[INFO] Swift S3 Auth provider used: %q", cp.ProviderName)

		awsConfig := &aws.Config{
			Credentials: creds,
			Region:      aws.String(GetRegion(nil, c)),
			// MaxRetries:       aws.Int(c.MaxRetries),
			HTTPClient: cleanhttp.DefaultClient(),
			// S3ForcePathStyle: aws.Bool(c.S3ForcePathStyle),
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

func buildClientByToken(c *Config) error {
	var pao, dao golangsdk.AuthOptions

	if c.AgencyDomainName != "" && c.AgencyName != "" {
		pao = golangsdk.AuthOptions{
			AgencyName:       c.AgencyName,
			AgencyDomainName: c.AgencyDomainName,
			DelegatedProject: c.DelegatedProject,
		}

		dao = golangsdk.AuthOptions{
			AgencyName:       c.AgencyName,
			AgencyDomainName: c.AgencyDomainName,
		}
	} else {
		pao = golangsdk.AuthOptions{
			DomainID:   c.DomainID,
			DomainName: c.DomainName,
			TenantID:   c.TenantID,
			TenantName: c.TenantName,
		}

		dao = golangsdk.AuthOptions{
			DomainID:   c.DomainID,
			DomainName: c.DomainName,
		}
	}

	for _, ao := range []*golangsdk.AuthOptions{&pao, &dao} {
		ao.IdentityEndpoint = c.IdentityEndpoint
		ao.TokenID = c.Token

	}
	return genClients(c, pao, dao)
}

func buildClientByAKSK(c *Config) error {
	var pao, dao golangsdk.AKSKAuthOptions

	if c.AgencyDomainName != "" && c.AgencyName != "" {
		pao = golangsdk.AKSKAuthOptions{
			DomainID:         c.DomainID,
			Domain:           c.DomainName,
			AgencyName:       c.AgencyName,
			AgencyDomainName: c.AgencyDomainName,
			DelegatedProject: c.DelegatedProject,
		}

		dao = golangsdk.AKSKAuthOptions{
			DomainID:         c.DomainID,
			Domain:           c.DomainName,
			AgencyName:       c.AgencyName,
			AgencyDomainName: c.AgencyDomainName,
		}
	} else {
		pao = golangsdk.AKSKAuthOptions{
			ProjectName: c.TenantName,
			ProjectId:   c.TenantID,
		}

		dao = golangsdk.AKSKAuthOptions{
			DomainID: c.DomainID,
			Domain:   c.DomainName,
		}
	}

	for _, ao := range []*golangsdk.AKSKAuthOptions{&pao, &dao} {
		ao.IdentityEndpoint = c.IdentityEndpoint
		ao.AccessKey = c.AccessKey
		ao.SecretKey = c.SecretKey
	}
	return genClients(c, pao, dao)
}

func buildClientByPassword(c *Config) error {
	var pao, dao golangsdk.AuthOptions

	if c.AgencyDomainName != "" && c.AgencyName != "" {
		pao = golangsdk.AuthOptions{
			DomainID:         c.DomainID,
			DomainName:       c.DomainName,
			AgencyName:       c.AgencyName,
			AgencyDomainName: c.AgencyDomainName,
			DelegatedProject: c.DelegatedProject,
		}

		dao = golangsdk.AuthOptions{
			DomainID:         c.DomainID,
			DomainName:       c.DomainName,
			AgencyName:       c.AgencyName,
			AgencyDomainName: c.AgencyDomainName,
		}
	} else {
		pao = golangsdk.AuthOptions{
			DomainID:   c.DomainID,
			DomainName: c.DomainName,
			TenantID:   c.TenantID,
			TenantName: c.TenantName,
		}

		dao = golangsdk.AuthOptions{
			DomainID:   c.DomainID,
			DomainName: c.DomainName,
		}
	}

	for _, ao := range []*golangsdk.AuthOptions{&pao, &dao} {
		ao.IdentityEndpoint = c.IdentityEndpoint
		ao.Password = c.Password
		ao.Username = c.Username
		ao.UserID = c.UserID
	}
	return genClients(c, pao, dao)
}

func genClients(c *Config, pao, dao golangsdk.AuthOptionsProvider) error {
	client, err := genClient(c, pao)
	if err != nil {
		return err
	}
	c.HwClient = client

	client, err = genClient(c, dao)
	if err == nil {
		c.DomainClient = client
	}
	return err
}

func genClient(c *Config, ao golangsdk.AuthOptionsProvider) (*golangsdk.ProviderClient, error) {
	client, err := openstack.NewClient(ao.GetIdentityEndpoint())
	if err != nil {
		return nil, err
	}

	// Set UserAgent
	client.UserAgent.Prepend(httpclient.TerraformUserAgent(c.terraformVersion))

	config, err := generateTLSConfig(c)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config}

	// if OS_DEBUG is set, log the requests and responses
	var osDebug bool
	if os.Getenv("OS_DEBUG") != "" {
		osDebug = true
	}

	client.HTTPClient = http.Client{
		Transport: &LogRoundTripper{
			Rt:         transport,
			OsDebug:    osDebug,
			MaxRetries: c.MaxRetries,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if client.AKSKAuthOptions.AccessKey != "" {
				golangsdk.ReSign(req, golangsdk.SignOptions{
					AccessKey: client.AKSKAuthOptions.AccessKey,
					SecretKey: client.AKSKAuthOptions.SecretKey,
				})
			}
			return nil
		},
	}

	// If using Swift Authentication, there's no need to validate authentication normally.
	if !c.Swauth {
		err = openstack.Authenticate(client, ao)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
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
		return nil, fmt.Errorf("missing credentials for Swift S3 Provider, need access_key and secret_key values for provider")
	}

	client, err := openstack.NewOBSService(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
	if err != nil {
		return nil, err
	}

	awsS3Sess := c.s3sess.Copy(&aws.Config{Endpoint: aws.String(client.Endpoint)})
	s3conn := s3.New(awsS3Sess)

	return s3conn, err
}

func setUpOBSLogging() {
	// init log
	if os.Getenv("OS_DEBUG") != "" {
		var logfile = "./.obs-sdk.log"
		// maxLogSize:10M, backups:10
		if err := obs.InitLog(logfile, 1024*1024*10, 10, obs.LEVEL_DEBUG, false); err != nil {
			log.Printf("[WARN] initial obs sdk log failed: %s", err)
		}
	}
}

// setupTemporaryCredentials creates temporary AK/SK, which can be used to auth in OBS when AK/SK is not provided
func (c *Config) setupTemporaryCredentials() error {
	if c.SecurityToken != "" || (c.AccessKey != "" && c.SecretKey != "") {
		return nil
	}
	client, err := c.identityV3Client()
	if err != nil {
		return fmt.Errorf("error creating identity v3 domain client: %s", err)
	}
	credential, err := credentials.CreateTemporary(client, credentials.CreateTemporaryOpts{
		Methods: []string{"token"},
		Token:   client.Token(),
	}).Extract()
	if err != nil {
		return fmt.Errorf("error creating temporary AK/SK: %s", err)
	}
	c.AccessKey = credential.AccessKey
	c.SecretKey = credential.SecretKey
	c.SecurityToken = credential.SecurityToken
	return nil
}

func (c *Config) newObjectStorageClient(region string) (*obs.ObsClient, error) {
	if err := c.setupTemporaryCredentials(); err != nil {
		return nil, fmt.Errorf("failed to construct OBS client without AK/SK: %s", err)
	}

	client, err := openstack.NewOBSService(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getEndpointType(),
	})
	if err != nil {
		return nil, err
	}

	setUpOBSLogging()

	return obs.New(c.AccessKey, c.SecretKey, client.Endpoint, obs.WithSecurityToken(c.SecurityToken))
}

func (c *Config) blockStorageV1Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewBlockStorageV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) blockStorageV2Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewBlockStorageV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) blockStorageV3Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewBlockStorageV3(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) cbrV3Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewCBRService(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) computeV1Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewComputeV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getEndpointType(),
	})
}

func (c *Config) computeV2Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewComputeV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) dnsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewDNSV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) identityV3Client(_ ...string) (*golangsdk.ServiceClient, error) {
	return openstack.NewIdentityV3(c.DomainClient, golangsdk.EndpointOpts{
		Availability: c.getEndpointType(),
	})
}

// identityV30Client - provides client is used for use with endpoints with invalid "v3.0" URLs
func (c *Config) identityV30Client() (*golangsdk.ServiceClient, error) {
	service, err := openstack.NewIdentityV3(c.DomainClient, golangsdk.EndpointOpts{
		Availability: c.getEndpointType(),
	})
	if err != nil {
		return nil, err
	}
	service.Endpoint = strings.Replace(service.IdentityEndpoint, "v3/", "v3.0/", 1)
	return service, nil
}

func (c *Config) imageV1Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewImageServiceV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) imageV2Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewImageServiceV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) networkingV1Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewNetworkV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) networkingV2Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewNetworkV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) objectStorageV1Client(region string) (*golangsdk.ServiceClient, error) {
	// If Swift Authentication is being used, return a swauth client.
	if c.Swauth {
		return swauth.NewObjectStorageV1(c.HwClient, swauth.AuthOpts{
			User: c.Username,
			Key:  c.Password,
		})
	}

	return openstack.NewObjectStorageV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) smnV2Client(projectName ProjectName) (*golangsdk.ServiceClient, error) {
	newConfig, err := reconfigProjectName(c, projectName)
	if err != nil {
		return nil, err
	}
	return openstack.NewSMNV2(newConfig.HwClient, golangsdk.EndpointOpts{
		Region:       GetRegion(nil, c),
		Availability: c.getEndpointType(),
	})
}

func (c *Config) cesV1Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewCESClient(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) getEndpointType() golangsdk.Availability {
	if c.EndpointType == "internal" || c.EndpointType == "internalURL" {
		return golangsdk.AvailabilityInternal
	}
	if c.EndpointType == "admin" || c.EndpointType == "adminURL" {
		return golangsdk.AvailabilityAdmin
	}
	return golangsdk.AvailabilityPublic
}

func (c *Config) kmsKeyV1Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewKMSV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) natV2Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewNatV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) orchestrationV1Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewOrchestrationV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) sfsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewSharedFileSystemV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) vbsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewVBS(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) autoscalingV1Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewAutoScalingService(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) csbsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewCSBSService(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) dehV1Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewDeHServiceV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) dmsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewDMSServiceV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) MrsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewMapReduceV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) elbV1Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewELBV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) rdsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewRDSV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) antiddosV1Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewAntiDDoSV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) ctsV1Client(projectName ProjectName) (*golangsdk.ServiceClient, error) {
	newConfig, err := reconfigProjectName(c, projectName)
	if err != nil {
		return nil, err
	}
	return openstack.NewCTSService(newConfig.HwClient, golangsdk.EndpointOpts{
		Region:       GetRegion(nil, c),
		Availability: c.getEndpointType(),
	})
}

func (c *Config) cssV1Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewCSSService(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) cceV3Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewCCE(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) cceV3AddonClient(region string) (*golangsdk.ServiceClient, error) {
	client, err := c.cceV3Client(region)
	if err != nil {
		return nil, err
	}
	client.ResourceBase = fmt.Sprintf("%sapi/v3/", client.Endpoint)
	return client, nil
}

func (c *Config) dcsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewDCSServiceV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) rdsTagV1Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewRdsTagV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) wafV1Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewWAFV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) rdsV3Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewRDSV3(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) sdrsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.SDRSV1(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) ltsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewLTSV2(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}

func (c *Config) ddsV3Client(region string) (*golangsdk.ServiceClient, error) {
	return openstack.NewDDSServiceV3(c.HwClient, golangsdk.EndpointOpts{
		Region:       region,
		Availability: c.getEndpointType(),
	})
}
