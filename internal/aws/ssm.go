package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// SSMExecuteCommand execute an SSM command and return command id
func SSMExecuteCommand(ses *session.Session, instanceIDs []string, documentName string, parameters map[string][]*string, comment string) (*string, error) {
	// Load session from shared config
	if ses == nil {
		ses = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	}

	// Create new SSM client
	ssmSvc := ssm.New(ses)
	result, err := ssmSvc.SendCommand(&ssm.SendCommandInput{
		Comment:      &comment,
		InstanceIds:  aws.StringSlice(instanceIDs),
		DocumentName: aws.String(documentName),
		Parameters:   parameters,
	})

	var commandID *string
	if err != nil {
		return commandID, err
	}

	commandID = result.Command.CommandId
	return commandID, nil
}

// SSMCommandResponse is the command's response
type SSMCommandResponse struct {
	Status       *string
	InstanceID   *string
	InstanceName *string
	Output       *string
}

// SSMWaitCommand wait command to end and return output
func SSMWaitCommand(ses *session.Session, commandID *string) ([]SSMCommandResponse, bool, error) {
	// Load session from shared config
	if ses == nil {
		ses = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	}

	// Create new SSM client
	ssmSvc := ssm.New(ses)
	var result *ssm.ListCommandInvocationsOutput
	var err error

	// Wait for command ingestion
	time.Sleep(2 * time.Second)

	// Wait for command status change
	commandStilRunning := true
	for end := true; end; end = commandStilRunning {
		result, err = ssmSvc.ListCommandInvocations(&ssm.ListCommandInvocationsInput{
			CommandId: commandID,
			Details:   aws.Bool(true),
		})

		// Check if all commands still running
		atLeastOneStillRunning := false
		for _, commandInvocation := range result.CommandInvocations {
			if *commandInvocation.Status == "Pending" || *commandInvocation.Status == "InProgress" || *commandInvocation.Status == "Delayed" {
				atLeastOneStillRunning = true
			}
		}

		if !atLeastOneStillRunning {
			commandStilRunning = false
		} else {
			time.Sleep(1 * time.Second)
		}
	}

	var responses []SSMCommandResponse
	if err != nil {
		return responses, false, err
	}

	// Format command responses
	allInvocationSuccess := true
	for _, commandInvocation := range result.CommandInvocations {
		output := ""
		if *commandInvocation.Status != "Success" {
			allInvocationSuccess = false
		}
		for _, commandPlugin := range commandInvocation.CommandPlugins {
			if len(output) > 0 {
				output += "\n"
			}
			output += *commandPlugin.Output
		}
		responses = append(responses, SSMCommandResponse{
			Status:       commandInvocation.Status,
			InstanceID:   commandInvocation.InstanceId,
			InstanceName: commandInvocation.InstanceName,
			Output:       &output,
		})
	}

	return responses, allInvocationSuccess, nil
}

// SSMDocument is the SSM list document response
type SSMDocument struct {
	Name  *string
	Owner *string
}

// SSMListDocuments return a list of SSM document
func SSMListDocuments(ses *session.Session, owner string) ([]SSMDocument, error) {
	if ses == nil {
		ses = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	}

	// Create new SSM client
	ssmSvc := ssm.New(ses)

	// List Documents
	var documents []SSMDocument
	ssmSvc.ListDocumentsPages(&ssm.ListDocumentsInput{
		Filters: []*ssm.DocumentKeyValuesFilter{
			{
				Key: aws.String("Owner"),
				Values: aws.StringSlice([]string{
					owner,
				}),
			},
		},
	}, func(page *ssm.ListDocumentsOutput, lastPage bool) bool {
		for _, documentIdentifier := range page.DocumentIdentifiers {
			documents = append(documents, SSMDocument{
				Name:  documentIdentifier.Name,
				Owner: documentIdentifier.Owner,
			})
		}
		return true // iterate over all pages
	})

	return documents, nil
}

// SSMGetDocumentParameters return a details about SSM document
func SSMGetDocumentParameters(ses *session.Session, documentName string) ([]*ssm.DocumentParameter, error) {
	if ses == nil {
		ses = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	}

	// List Documents
	ssmSvc := ssm.New(ses)
	describeResponse, err := ssmSvc.DescribeDocument(&ssm.DescribeDocumentInput{
		Name: aws.String(documentName),
	})
	var parameters []*ssm.DocumentParameter
	if err != nil {
		return parameters, err
	}

	if describeResponse.Document.Parameters == nil {
		return make([]*ssm.DocumentParameter, 0), nil
	}
	parameters = describeResponse.Document.Parameters

	return parameters, nil
}
