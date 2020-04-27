package aws

import (
	"awslib"
	"fmt"
	"shelllib"

	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"
)

// SSMSelectInstance list instances to connect to
func SSMSelectInstance(c *cli.Context) error {
	// Check if instance is provided
	instanceID := c.String("instance")
	if len(instanceID) != 0 {
		// Start SSM session
		return SSMStartSession(c)
	}

	// Create AWS session
	currentSession := CreateAWSSession(c)

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

	// If only one instance is found connect to it
	if len(instances) == 1 {
		c.Set("instance", *instances[0].ID)
		return SSMStartSession(c)
	}

	// Build table
	header := fmt.Sprintf("%-20s\t%-15s\t%s\t%s", "Instace ID", "IP address", "Environment", "ServiceType")
	var options []string
	for _, instance := range instances {
		options = append(options, fmt.Sprintf("%-20s\t%-15s\t%-8s\t%s", *instance.ID, *instance.IP, *instance.Environment, *instance.ServiceType))
	}

	// Ask selection
	instanceSelectedIndex := -1
	prompt := &survey.Select{
		Message:  "Select an instance: \n\n  " + header + "\n",
		Options:  options,
		PageSize: 15,
	}
	survey.AskOne(prompt, &instanceSelectedIndex)
	fmt.Println("")

	// Check response
	if instanceSelectedIndex == -1 {
		fmt.Println("\nNo instance selected")
		return nil
	}

	// Start SSM session
	c.Set("instance", *instances[instanceSelectedIndex].ID)
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
