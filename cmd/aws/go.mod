module aws

go 1.14

require (
	github.com/AlecAivazis/survey/v2 v2.0.7
	github.com/aws/aws-sdk-go v1.30.9
	github.com/urfave/cli/v2 v2.2.0
	awslib v0.0.0
	shelllib v0.0.0
	configlib v0.0.0
)

replace (
	awslib => ./../../lib/aws
	shelllib => ./../../lib/shell
	configlib => ./../../lib/config
)
