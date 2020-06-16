package aws

import (
	"configlib"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/urfave/cli/v2"
)

// Config is struct for AWS
type Config struct {
	AWS struct {
		Profile string
		Region  string
		ECS     struct {
			Cluster string
		}
	}
}

// Global config
var config Config

// CreateAWSSession return a new AWS client session
func CreateAWSSession(c *cli.Context) *session.Session {
	// Check for AWS profile
	profile := c.String("profile")
	if len(profile) == 0 && len(config.AWS.Profile) != 0 {
		profile = config.AWS.Profile
	}

	// Check for region
	region := c.String("region")
	awsConfig := aws.Config{}
	if len(region) != 0 {
		awsConfig.Region = aws.String(c.String("region"))
	} else if len(region) == 0 && len(config.AWS.Region) != 0 {
		awsConfig.Region = aws.String(config.AWS.Region)
	} else {
		awsConfig.Region = aws.String("eu-west-1")
	}
	c.Set("region", *awsConfig.Region)

	return session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
		Config:            awsConfig,
	}))
}

// Commands - Return all commands
func Commands() []*cli.Command {
	globalFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "profile",
			Aliases: []string{"p"},
			Usage:   "AWS profile name",
			EnvVars: []string{"AWS_PROFILE", "AWS_DEFAULT_PROFILE"},
		},
		&cli.StringFlag{
			Name:    "region",
			Aliases: []string{"r"},
			Usage:   "AWS region",
			EnvVars: []string{"AWS_DEFAULT_REGION"},
		},
	}

	return []*cli.Command{
		{
			Name: "aws",
			Before: func(c *cli.Context) error {
				configlib.LoadConfig(c.String("config"), &config)
				// Store args into config, child commands cannot access it
				config.AWS.Region = c.String("region")
				config.AWS.Profile = c.String("profile")
				return nil
			},
			Flags: globalFlags,
			Subcommands: []*cli.Command{
				NewECSConnectCommand(globalFlags),
				NewSSMConnectCommand(globalFlags),
				NewSSMRunCommand(globalFlags),
				NewSSMRunDocumentCommand(globalFlags),
				NewSSMTunnelCommand(globalFlags),
			},
		},
	}
}
