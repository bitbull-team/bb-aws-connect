package aws

import (
	"awslib"
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/urfave/cli/v2"
)

// SSMSelectInstances select multiple instances
func SSMSelectInstances(c *cli.Context) error {
	// Check if instance is provided
	instanceIDs := c.StringSlice("instance")
	if len(instanceIDs) != 0 {
		// Start SSM session
		return SSMRunCommand(c)
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
	header := fmt.Sprintf("%-20s\t%-15s\t%s\t%s", "Instace ID", "IP address", "Environment", "ServiceType")
	var options []string
	for _, instance := range instances {
		options = append(options, fmt.Sprintf("%-20s\t%-15s\t%-8s\t%s", *instance.ID, *instance.IP, *instance.Environment, *instance.ServiceType))
	}

	// Ask selection
	instancesSelectedIndex := []int{}
	prompt := &survey.MultiSelect{
		Message:  "Select an instance: \n\n  " + header + "\n",
		Options:  options,
		PageSize: 15,
	}
	survey.AskOne(prompt, &instancesSelectedIndex)
	fmt.Println("")

	// Check response
	if len(instancesSelectedIndex) == 0 {
		fmt.Println("\nNo instances selected")
		return nil
	}

	// Set instance ids
	for _, instanceSelectedIndex := range instancesSelectedIndex {
		c.Set("instance", *instances[instanceSelectedIndex].ID)
	}

	return SSMRunCommand(c)
}

// SSMRunCommand execute command to remote instance
func SSMRunCommand(c *cli.Context) error {
	// Check command
	scriptFile := c.String("file")
	var command string
	var commandRows []string

	// Check command from first argument
	if len(command) == 0 && len(scriptFile) == 0 && c.Args().Present() {
		command = c.Args().First()
	}

	// Ask command
	if len(command) == 0 && len(scriptFile) == 0 {
		prompt := &survey.Multiline{
			Message: "Type command to execute: ",
		}
		survey.AskOne(prompt, &command)
		fmt.Println("")

		if len(command) == 0 {
			return cli.Exit("No command or file arguments provided", -1)
		}
	}

	// Check command arguments
	if len(command) == 0 {
		// Read script file
		file, err := os.Open(scriptFile)
		if err != nil {
			return cli.Exit("Cannot open script file "+scriptFile+": "+err.Error(), -1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			commandRows = append(commandRows, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			return cli.Exit("Error reading script file "+scriptFile+": "+err.Error(), -1)
		}
	} else {
		// Split command by line
		commandRows = strings.Split(command, "\n")
	}

	// Create AWS session
	currentSession := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           c.String("profile"),
		SharedConfigState: session.SharedConfigEnable,
	}))

	// List available service
	commandID, err := awslib.SSMExecuteCommand(
		currentSession,
		c.StringSlice("instance"),
		commandRows,
		"AWS-RunShellScript",
		"Executed from bb-cli",
	)
	if err != nil {
		return cli.Exit("Error before SSM command execution: "+err.Error(), -1)
	}

	// Wait until all commands ends
	responses, allSuccess, errWait := awslib.SSMWaitCommand(currentSession, commandID)
	if errWait != nil {
		return cli.Exit("Error during SSM command execution: "+errWait.Error(), -1)
	}

	if allSuccess {
		fmt.Println("All commands ends successfully!")
	} else {
		fmt.Println("Some commands ends with errors")
	}

	// Show output
	fmt.Println("")
	for _, response := range responses {
		fmt.Println("--------------------------------")
		fmt.Printf("%-20s\t%-10s\n\n%s\n", *response.InstanceID, *response.Status, *response.Output)
	}
	fmt.Println("--------------------------------")
	return nil
}
