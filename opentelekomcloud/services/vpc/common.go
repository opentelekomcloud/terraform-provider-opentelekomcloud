package vpc

const (
	errCreationV3Client = "error creating OpenTelekomCloud NetworkingV3 client: %w"
	errCreationV2Client = "error creating OpenTelekomCloud NetworkingV2 client: %w"
	errCreationV1Client = "error creating OpenTelekomCloud NetworkingV1 client: %w"
	keyClientV2         = "vpc-v2-client"
	keyClientV1         = "vpc-v1-client"
	// MaxCreateRoutes is the limitation of creating API
	MaxCreateRoutes int = 5
)
