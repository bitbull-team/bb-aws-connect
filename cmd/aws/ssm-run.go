package aws

import (
	"awslib"
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/urfave/cli/v2"
)

// NewSSMRunCommand return "ssm:run" command
func NewSSMRunCommand(globalFlags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:      "ssm:run",
		Usage:     "Run command to EC2 instances using a SSM command",
		ArgsUsage: "[command to execute]",
		Action:    SSMRun,
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
			&cli.StringFlag{
				Name:   "command",
				Hidden: true,
			},
			&cli.StringFlag{
				Name:   "document",
				Hidden: true,
			},
		}...),
	}
}

// SSMRun run a command to an EC2 instance using SSM
func SSMRun(c *cli.Context) error {
	var err error
	// Select multiple EC2 instances
	err = SSMSelectInstances(c)
	if err != nil {
		return err
	}

	// Select command to run
	err = SSMSelectCommand(c)
	if err != nil {
		return err
	}

	// Run SSM commands
	c.Set("document", "AWS-RunShellScript")
	commandRows := strings.Split(c.String("command"), "\n")
	err = SSMRunCommands(c, map[string][]*string{
		"commands": aws.StringSlice(commandRows),
	})
	if err != nil {
		return err
	}

	return nil
}

// SSMSelectInstances select multiple instances
func SSMSelectInstances(c *cli.Context) error {
	// Check if instance is provided
	instanceIDs := c.StringSlice("instance")
	if len(instanceIDs) != 0 {
		// Start SSM session
		return nil
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
		return cli.Exit("Error during EC2 instance list: "+err.Error(), 1)
	}
	if len(instances) == 0 {
		return cli.Exit("No instances found", 1)
	}

	// If only one instance is found connect to it
	if len(instances) == 1 {
		fmt.Println("Instace auto selected: ", *instances[0].ID)
		c.Set("instance", *instances[0].ID)
		return nil
	}

	// Build table
	header := fmt.Sprintf("%-20s\t%-15s\t%s\t%s", "Instace ID", "IP address", "Environment", "ServiceType")
	var options []string
	for _, instance := range instances {
		options = append(options, fmt.Sprintf("%-20s\t%-15s\t%-8s\t%s", *instance.ID, *instance.IP, *instance.Environment, *instance.ServiceType))
	}

	// Check if auto select is set
	if c.Bool("auto-select") {
		for _, instance := range instances {
			c.Set("instance", *instance.ID)
		}
		return nil
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
		return cli.Exit("No instances selected", 1)
	}

	// Set instance ids
	for _, instanceSelectedIndex := range instancesSelectedIndex {
		c.Set("instance", *instances[instanceSelectedIndex].ID)
	}

	return nil
}

// SSMSelectCommand select SSM command
func SSMSelectCommand(c *cli.Context) error {
	var command string
	scriptFile := c.String("file")

	// Check command from first argument
	if len(command) == 0 && len(scriptFile) == 0 && c.Args().Present() {
		command = c.Args().First()
	}

	// Check script file
	if len(command) == 0 && len(scriptFile) != 0 {
		// Check command arguments
		if len(scriptFile) != 0 {
			// Read script file
			file, err := os.Open(scriptFile)
			if err != nil {
				return cli.Exit("Cannot open script file "+scriptFile+": "+err.Error(), 1)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				command += scanner.Text() + "\n"
			}

			if err := scanner.Err(); err != nil {
				return cli.Exit("Error reading script file "+scriptFile+": "+err.Error(), 1)
			}
		}
	}

	// Ask command
	if len(command) == 0 {
		prompt := &survey.Multiline{
			Message: "Type command to execute: ",
		}
		survey.AskOne(prompt, &command)
		fmt.Println("")

		if len(command) == 0 {
			return cli.Exit("No command or file arguments provided", 1)
		}
	}

	c.Set("command", command)
	return nil
}

// SSMRunCommands execute command to remote instance
func SSMRunCommands(c *cli.Context, parameters map[string][]*string) error {
	// Create AWS session
	currentSession := CreateAWSSession(c)

	// Execute SSM command
	commandID, err := awslib.SSMExecuteCommand(
		currentSession,
		c.StringSlice("instance"),
		c.String("document"),
		parameters,
		"Executed from bb-cli",
	)
	if err != nil {
		return cli.Exit("Error before SSM command execution: "+err.Error(), 1)
	}

	// Wait until all commands ends
	fmt.Println("Waiting for command id ", *commandID, "..")
	responses, allSuccess, errWait := awslib.SSMWaitCommand(currentSession, commandID)
	if errWait != nil {
		return cli.Exit("Error during SSM command execution: "+errWait.Error(), 1)
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
