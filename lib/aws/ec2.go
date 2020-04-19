package awslib

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// TagFilter is a filter that can be applied to EC2ListInstances
type TagFilter struct {
	Name  string
	Value string
}

// Instace is a result of EC2ListInstances
type Instace struct {
	ID          string
	Name        string
	IP          string
	ServiceType string
	Environment string
	ProjectName string
}

// EC2ListInstances return a list of available instance
func EC2ListInstances(ses *session.Session, tags []TagFilter) ([]Instace, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
				},
			},
		},
	}

	if tags != nil {
		for _, tag := range tags {
			input.Filters = append(input.Filters, &ec2.Filter{
				Name: aws.String("tag:" + tag.Name),
				Values: []*string{
					aws.String(tag.Value),
				},
			})
		}
	}

	// Load session from shared config
	if ses == nil {
		ses = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	}

	// Create new EC2 client
	ec2Svc := ec2.New(ses)
	result, err := ec2Svc.DescribeInstances(input)
	formattedInstances := make([]Instace, 0)
	if err != nil {
		return formattedInstances, err
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			var tagName string
			var serviceName string
			var envName string
			var projectName string
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					tagName = *tag.Value
				}
				if *tag.Key == "ServiceType" {
					serviceName = *tag.Value
				}
				if *tag.Key == "Environment" {
					envName = *tag.Value
				}
				if *tag.Key == "ProjectName" {
					projectName = *tag.Value
				}
			}
			formattedInstances = append(formattedInstances, Instace{
				ID:          *instance.InstanceId,
				IP:          *instance.PrivateIpAddress,
				Name:        tagName,
				ServiceType: serviceName,
				Environment: envName,
				ProjectName: projectName,
			})
		}
	}

	return formattedInstances, nil
}
