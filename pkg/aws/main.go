package awslib

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/urfave/cli/v2"
)

// Config is struct for AWS
type Config struct {
	Profile string
	Region  string
}

// CreateAWSSession return a new AWS client session
func CreateAWSSession(c *cli.Context, config Config) *session.Session {
	// Check for AWS profile
	profile := c.String("profile")
	if len(profile) == 0 && len(config.Profile) != 0 {
		profile = config.Profile
	}

	// Check for region
	region := c.String("region")
	awsConfig := aws.Config{}
	if len(region) != 0 {
		awsConfig.Region = aws.String(c.String("region"))
	} else if len(region) == 0 && len(config.Region) != 0 {
		awsConfig.Region = aws.String(config.Region)
	} else {
		awsConfig.Region = aws.String("eu-west-1")
	}
	c.Set("region", *awsConfig.Region)

	// Check for debug mode
	debugMode := os.Getenv("BB_AWS_CONNECT_AWS_DEBUG")
	if debugMode != "" {
		awsConfig.LogLevel = aws.LogLevel(aws.LogDebugWithHTTPBody)
	}

	return session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
		Config:            awsConfig,
	}))
}
