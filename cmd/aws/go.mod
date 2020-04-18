module aws

go 1.14

require (
	awslib v0.0.0
	github.com/aws/aws-sdk-go v1.30.9
	github.com/urfave/cli/v2 v2.2.0
)

replace awslib => ./../../lib/aws
