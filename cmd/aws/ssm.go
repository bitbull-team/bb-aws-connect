package aws

import (
	"awslib"
	"fmt"
	"shelllib"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/urfave/cli/v2"
)

// SSMListInstances list instances to connect to
func SSMListInstances(c *cli.Context) error {
	// Check if instance is provided
	instanceID := c.String("instance")
	if len(instanceID) != 0 {
		// Start SSM session
		return SSMStartSession(c)
	}

	// Create AWS session
	currentSession := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           c.String("profile"),
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Build filters
	var tagFilters []awslib.TagFilter
	env := c.String("env")
	if env != "" {
		tagFilters = append(tagFilters, awslib.TagFilter{
			Name:  "Environment",
			Value: env,
		})
	}
	serviceType := c.String("service")
	if serviceType != "" {
		tagFilters = append(tagFilters, awslib.TagFilter{
			Name:  "ServiceType",
			Value: serviceType,
		})
	}

	// List available instance
	instances, err := awslib.EC2ListInstances(currentSession, tagFilters)
	if err != nil {
		return cli.Exit("Error during EC2 instance list: "+err.Error(), -1)
	}
	if len(instances) == 0 {
		return cli.Exit("No instances found", -1)
	}

	// Build table
	header := fmt.Sprintf("%s\t%s\t%s\t%s", "Instace ID    ", "IP address", "Environment", "ServiceType")
	var options []string
	for _, instance := range instances {
		options = append(options, fmt.Sprintf("%s\t%s\t%s    \t%s", instance.ID, instance.IP, instance.Environment, instance.ServiceType))
	}

	// Ask selection
	instanceSelected := ""
	prompt := &survey.Select{
		Message:  "Select an instance: \n\n  " + header + "\n",
		Options:  options,
		PageSize: 15,
	}
	survey.AskOne(prompt, &instanceSelected)

	// Check response
	instanceID = strings.Split(instanceSelected, "\t")[0]
	if len(instanceID) == 0 {
		fmt.Println("No instances selected")
		return nil
	}

	// Start SSM session
	c.Set("instance", instanceID)
	return SSMStartSession(c)
}

// SSMStartSession connect to a instance
func SSMStartSession(c *cli.Context) error {
	// Get parameters
	profile := c.String("profile")
	region := c.String("region")
	instanceID := c.String("instance")

	// Build arguments
	args := []string{
		"ssm", "start-session",
		"--profile", profile,
		"--region", region,
		"--target", instanceID,
		"--document-name", "AWS-StartInteractiveCommand",
	}

	// Check command
	command := c.String("command")
	if len(command) == 0 {
		// Additional arguments
		cwd := c.String("cwd")
		user := c.String("user")
		shell := c.String("shell")

		// Build extra arguments
		if cwd != "" || user != "" || shell != "" {
			// Change CWD
			if cwd != "" {
				command = fmt.Sprintf("cd %s", cwd)
			}

			// Concatenate CWD and user/shell
			if (user != "" || shell != "") && len(command) > 0 {
				command += " && "
			}

			// Change user and shell
			if user == "" && shell != "" {
				command += shell
			} else if user != "" && shell == "" {
				command += fmt.Sprintf("sudo su %s", user)
			} else if user != "" && shell != "" {
				command += fmt.Sprintf("sudo su %s -s %s", user, shell)
			}
		}
	}

	// Check if command is still 0 after checking cwd, user and shell
	if len(command) != 0 {
		args = append(args, "--parameters", fmt.Sprintf("command=\"%s\"", command))
	}

	// Start SSM session
	shelllib.ExecuteCommandForeground("aws", args...)
	return nil
}
