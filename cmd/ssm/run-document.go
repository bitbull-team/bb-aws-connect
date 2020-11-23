package ssm

import (
	"fmt"
	"strings"

	"github.com/bitbull-team/bb-aws-connect/internal/aws"

	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"
)

// NewRunDocumentCommand return "ssm:run:document" command
func NewRunDocumentCommand(globalFlags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:   "run-document",
		Usage:  "Run a SSM document to EC2 instances",
		Action: RunDocument,
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
			&cli.BoolFlag{
				Name:    "auto-select",
				Aliases: []string{"a"},
				Usage:   "Automatically select all instance listed without asking",
			},
			&cli.BoolFlag{
				Name:    "self",
				Usage:   "SSM filter document with owner Self",
				Aliases: []string{"o"},
			},
			&cli.BoolFlag{
				Name:    "private",
				Usage:   "SSM filter document with owner Private",
				Aliases: []string{"t"},
			},
			&cli.StringFlag{
				Name:    "document",
				Usage:   "SSM document name",
				Aliases: []string{"d"},
			},
			&cli.StringSliceFlag{
				Name:    "parameter",
				Usage:   "SSM document parameter",
				Aliases: []string{"m"},
			},
		}...),
	}
}

// RunDocument run a SSM document to an EC2 instance
func RunDocument(c *cli.Context) error {
	var err error

	// Select document command
	err = SelectDocument(c)
	if err != nil {
		return err
	}

	// Check parameters
	var parameters map[string][]*string
	parameters, err = CheckDocumentParameters(c)
	if err != nil {
		return err
	}

	// Select multiple EC2 instances
	err = SelectInstances(c)
	if err != nil {
		return err
	}

	// Run SSM commands
	err = RunCommands(c, parameters)
	if err != nil {
		return err
	}

	return nil
}

// SelectDocument list documents ro run
func SelectDocument(c *cli.Context) error {
	// Check if document is provided
	document := c.String("document")
	if len(document) != 0 {
		return nil
	}

	// Create AWS session
	currentSession := aws.CreateAWSSession(c, aws.Config{
		Profile: globalConfig.Profile,
		Region:  globalConfig.Region,
	})

	// Check owner
	owner := "Amazon"
	if c.Bool("self") {
		owner = "Self"
	}
	if c.Bool("private") {
		owner = "Private"
	}

	// List documents
	documents, err := aws.SSMListDocuments(currentSession, owner)
	if err != nil {
		return cli.Exit("Error during SSM documents list: "+err.Error(), 1)
	}
	if len(documents) == 0 {
		return cli.Exit("No documents found", 1)
	}

	// Build table
	var options []string
	for _, document := range documents {
		if strings.HasPrefix(*document.Name, "arn:aws:ssm") {
			options = append(options, strings.Split(*document.Name, ":document/")[1])
		} else {
			options = append(options, *document.Name)
		}
	}

	// Ask selection
	documentSelectedIndex := -1
	prompt := &survey.Select{
		Message:  "Select a document: ",
		Options:  options,
		PageSize: 10,
	}
	survey.AskOne(prompt, &documentSelectedIndex)
	fmt.Println("")

	// Check response
	if documentSelectedIndex == -1 {
		return cli.Exit("No document selected", 1)
	}

	// Set document name
	c.Set("document", *documents[documentSelectedIndex].Name)
	return nil
}

// CheckDocumentParameters check for parameters
func CheckDocumentParameters(c *cli.Context) (map[string][]*string, error) {
	parametersValues := make(map[string][]*string)
	commandParametersValues := make(map[string][]*string)

	// Create AWS session
	currentSession := aws.CreateAWSSession(c, aws.Config{
		Profile: globalConfig.Profile,
		Region:  globalConfig.Region,
	})

	// List documents
	parameters, err := aws.SSMGetDocumentParameters(currentSession, c.String("document"))
	if err != nil {
		return commandParametersValues, cli.Exit("Error during SSM document parameters list: "+err.Error(), 1)
	}

	// Document has no parameters
	if len(parameters) == 0 {
		return commandParametersValues, nil
	}

	// Parse provided parameters
	for _, parameter := range c.StringSlice("parameter") {
		if len(parameter) == 0 {
			break
		}

		parameterParts := strings.Split(parameter, "=")
		if len(parameterParts) > 1 {
			parametersValues[strings.ToLower(parameterParts[0])] = aws.StringSlice([]string{
				parameterParts[1],
			})
		}
	}

	// Ask for parameters
	for _, parameter := range parameters {
		parameterKey := strings.ToLower(*parameter.Name)
		currentValue := parametersValues[parameterKey]
		if currentValue != nil && len(currentValue) != 0 {
			commandParametersValues[*parameter.Name] = currentValue
			break
		}

		// Format default value
		var defaultValue string
		if parameter.DefaultValue != nil {
			defaultValue = *parameter.DefaultValue
		}

		// Format prompt
		var prompt survey.Prompt
		if *parameter.Type == "String" || *parameter.Type == "Integer" {
			prompt = &survey.Input{
				Message: *parameter.Description,
				Default: defaultValue,
			}
		} else if *parameter.Type == "StringList" {
			prompt = &survey.Multiline{
				Message: *parameter.Description,
				Default: defaultValue,
			}
		} else if *parameter.Type == "Boolean" {
			prompt = &survey.Confirm{
				Message: *parameter.Description,
				Default: defaultValue == "true",
			}
		}

		// Ask selection
		var response string
		survey.AskOne(prompt, &response)
		commandParametersValues[*parameter.Name] = aws.StringSlice([]string{
			response,
		})
		fmt.Println("")
	}

	return commandParametersValues, nil
}
