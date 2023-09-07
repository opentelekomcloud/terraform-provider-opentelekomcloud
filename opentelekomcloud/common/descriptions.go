package common

var Descriptions = map[string]string{
	"access_key": "The access key for API operations. You can retrieve this\n" +
		"from the 'My Credential' section of the console.",

	"secret_key": "The secret key for API operations. You can retrieve this\n" +
		"from the 'My Credential' section of the console.",

	"auth_url": "The Identity authentication URL.",

	"region": "The OpenTelekomCloud region to connect to.",

	"user_name": "Username to login with.",

	"user_id": "User ID to login with.",

	"tenant_id": "The ID of the Tenant (Identity v2) or Project (Identity v3)\n" +
		"to login with.",

	"tenant_name": "The name of the Tenant (Identity v2) or Project (Identity v3)\n" +
		"to login with.",

	"password": "Password to login with.",

	"token": "Authentication token to use as an alternative to username/password.",

	"security_token": "Security token to use for OBS federated authentication.",

	"domain_id": "The ID of the Domain to scope to (Identity v3).",

	"domain_name": "The name of the Domain to scope to (Identity v3).",

	"insecure": "Trust self-signed certificates.",

	"cacert_file": "A Custom CA certificate.",

	"endpoint_type": "The catalog endpoint type to use.",

	"cert": "A client certificate to authenticate with.",

	"key": "A client private key to authenticate with.",

	"swauth": "Use Swift's authentication system instead of Keystone. Only used for\n" +
		"interaction with Swift.",

	"agency_name": "The name of agency",

	"agency_domain_name": "The name of domain who created the agency (Identity v3).",

	"delegated_project": "The name of delegated project (Identity v3).",

	"cloud": "An entry in a `clouds.yaml` file to use.",

	"max_retries": "How many times HTTP connection should be retried until giving up.",

	"max_backoff_retries": "How many times HTTP request should be retried when rate limit reached",

	"backoff_retry_timeout": "Timeout in seconds for backoff retry",

	"passcode": "One-time MFA passcode",
}
