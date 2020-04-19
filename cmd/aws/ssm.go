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

// SSMListAndStartSession will list instances to connect to
func SSMListAndStartSession(c *cli.Context) error {
	// Check if instance is provided
	instanceID := c.String("instance")
	if len(instanceID) != 0 {
		// Start SSM session
		c.Set("instance", instanceID)
		return SSMStartSession(c)
	}

	// Create AWS session
	profile := c.String("profile")
	currentSession := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           profile,
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

	// List availlable instance
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

// SSMStartSession will connect to a instance
func SSMStartSession(c *cli.Context) error {
	// Get parameters
	profile := c.String("profile")
	region := c.String("region")
	instanceID := c.String("instance")

	// Start SSM session
	shelllib.ExecuteCommandForeground("aws", "ssm", "start-session", "--profile", profile, "--region", region, "--target", instanceID)
	return nil
}
