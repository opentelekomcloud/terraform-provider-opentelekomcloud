package ecs

const (
	errCreateClient   = "error creating OpenTelekomCloud ComputeV1 client: %w"
	errCreateV2Client = "error creating OpenTelekomCloud ComputeV2 client: %w"

	keyClientV2  = "ecs-v2-client"
	keyClientV1  = "ecs-v1-client"
	microversion = "2.55"
)
