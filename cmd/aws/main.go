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
	}
}

// CreateAWSSession return a new AWS client session
func CreateAWSSession(c *cli.Context) *session.Session {
	var config Config
	configlib.LoadConfig("", &config)

	// Check for AWS profile
	profile := c.String("profile")
	if len(profile) == 0 {
		profile = config.AWS.Profile
	}
	c.Set("profile", profile)

	// Check for region
	region := c.String("region")
	awsConfig := aws.Config{}
	if len(region) != 0 {
		awsConfig.Region = aws.String(c.String("region"))
	} else if len(region) == 0 && len(config.AWS.Region) != 0 {
		awsConfig.Region = aws.String(region)
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
			EnvVars: []string{"AWS_PROFILE"},
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
			Subcommands: []*cli.Command{
				{
					Name:   "ssm:connect",
					Usage:  "Connect to an EC2 instance using SSM session",
					Action: SSMSelectInstance,
					Flags: append(globalFlags, []cli.Flag{
						&cli.StringFlag{
							Name:    "service",
							Aliases: []string{"s"},
							Usage:   "Service Type (example: bastion, frontend, varnish)",
						},
						&cli.StringFlag{
							Name:    "env",
							Aliases: []string{"e"},
							Usage:   "Environment (example: test, stage, prod)",
						},
						&cli.StringFlag{
							Name:    "instance",
							Aliases: []string{"i"},
							Usage:   "Instace ID (example: i-xxxxxxxxxxxxxxxxx)",
						},
						&cli.StringFlag{
							Name:  "cwd",
							Usage: "Current working directory (example: /var/www/)",
							Value: "/",
						},
						&cli.StringFlag{
							Name:  "user",
							Usage: "User to use in the session",
							Value: "root",
						},
						&cli.StringFlag{
							Name:  "shell",
							Usage: "Shell used in session",
							Value: "/bin/bash",
						},
						&cli.StringFlag{
							Name:  "command",
							Usage: "Use a custom command as entrypoint",
						},
					}...),
				},
				{
					Name:      "ssm:run",
					Usage:     "Run command to EC2 instances using a SSM command",
					ArgsUsage: "[command to execute]",
					Action:    SSMSelectInstances,
					Flags: append(globalFlags, []cli.Flag{
						&cli.StringFlag{
							Name:    "service",
							Aliases: []string{"s"},
							Usage:   "Service Type (example: bastion, frontend, varnish)",
						},
						&cli.StringFlag{
							Name:    "env",
							Aliases: []string{"e"},
							Usage:   "Environment (example: test, stage, prod)",
						},
						&cli.StringSliceFlag{
							Name:    "instance",
							Aliases: []string{"i"},
							Usage:   "Instace ID (example: i-xxxxxxxxxxxxxxxxx)",
						},
						&cli.StringFlag{
							Name:  "file",
							Usage: "Script file path to execute (example: ./my-script.sh)",
						},
						&cli.BoolFlag{
							Name:    "auto-select",
							Aliases: []string{"a"},
							Usage:   "Automatically select all instance listed without asking",
						},
					}...),
				},
				{
					Name:   "ecs:connect",
					Usage:  "Connect to an ECS Task container",
					Action: ECSListServices,
					Flags: append(globalFlags, []cli.Flag{
						&cli.StringFlag{
							Name:     "cluster",
							Aliases:  []string{"c"},
							Usage:    "Cluster Name",
							Value:    "default",
							Required: true,
						},
						&cli.StringFlag{
							Name:    "service",
							Aliases: []string{"s"},
							Usage:   "Service name (example: my-service)",
						},
						&cli.StringFlag{
							Name:    "task",
							Aliases: []string{"t"},
							Usage:   "Task ID (example: xxxxxxxxxxxxxxxxxxxx)",
						},
						&cli.StringFlag{
							Name:   "container",
							Hidden: true,
						},
						&cli.StringFlag{
							Name:   "instance",
							Hidden: true,
						},
						&cli.StringFlag{
							Name:    "workdir",
							Aliases: []string{"w"},
							Usage:   "Docker exec 'workdir' parameters (example: /app)",
						},
						&cli.StringFlag{
							Name:    "user",
							Aliases: []string{"u"},
							Usage:   "Docker exec 'user' parameters (example: www-data)",
						},
						&cli.StringFlag{
							Name:  "command",
							Usage: "Use a custom command as entrypoint",
							Value: "/bin/bash",
						},
					}...),
				},
			},
		},
	}
}
