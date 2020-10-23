module bb-aws-connect

go 1.14

require (
	awslib v0.0.0
	configlib v0.0.0
	ecs v0.0.0
	github.com/aws/aws-sdk-go v1.30.9
	github.com/urfave/cli/v2 v2.2.0
	shelllib v0.0.0
	ssm v0.0.0
)

replace (

	awslib => ./pkg/aws
	configlib => ./pkg/config
	ecs => ./cmd/ecs
	shelllib => ./pkg/shell
	ssm => ./cmd/ssm
)
