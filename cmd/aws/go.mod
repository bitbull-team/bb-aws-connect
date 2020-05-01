module aws

go 1.14

require (
	awslib v0.0.0
	configlib v0.0.0
	github.com/AlecAivazis/survey/v2 v2.0.7
	github.com/aws/aws-sdk-go v1.30.9
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b
	github.com/urfave/cli/v2 v2.2.0
	shelllib v0.0.0
)

replace (
	awslib => ./../../lib/aws
	configlib => ./../../lib/config
	shelllib => ./../../lib/shell
)
