package migrations

const (
	keyClientV1         = "dataarts-migrations-v1-client"
	errCreationV1Client = "error creating OpenTelekomCloud DataArts Migrations V1 client: %w"
)

type (
	ClusterType string
	FlavorType  string
)

const ClusterTypeCDM ClusterType = "cdm"

const (
	FlavorTypeSmall  FlavorType = "a79fd5ae-1833-448a-88e8-3ea2b913e1f6"
	FlavorTypeMedium FlavorType = "fb8fe666-6734-4b11-bc6c-43d11db3c745"
	FlavorTypeLarge  FlavorType = "5ddb1071-c5d7-40e0-a874-8a032e81a697"
	FlavorTypeXLarge FlavorType = "6ddb1072-c5d7-40e0-a874-8a032e81a698"
)
