package awslib

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// ECSService is a result of EC2ListInstances
type ECSService struct {
	Name string
}

// ECSListServices return a list of Service
func ECSListServices(ses *session.Session, cluster string) ([]ECSService, error) {
	// Load session from shared config
	if ses == nil {
		ses = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	}

	// Create new ECS client and list services
	ecsSvc := ecs.New(ses)
	listResult, listErr := ecsSvc.ListServices(&ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	})
	var formattedServices []ECSService
	if listErr != nil {
		return formattedServices, listErr
	}

	// Retrieve ECS services details
	describeResult, describeErr := ecsSvc.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: listResult.ServiceArns,
	})
	if describeErr != nil {
		return formattedServices, describeErr
	}

	for _, service := range describeResult.Services {
		formattedServices = append(formattedServices, ECSService{
			Name: *service.ServiceName,
		})
	}

	return formattedServices, nil
}

// ECSTask is a result of ECSListServiceTasks
type ECSTask struct {
	Arn               *string
	ContainerInstance *ecs.ContainerInstance
	HealthStatus      *string
	Status            *string
	TaskDefinition    *ecs.TaskDefinition
}

// ECSListServiceTasks return a list of Service's Tasks
func ECSListServiceTasks(ses *session.Session, cluster string, serviceName string) ([]ECSTask, error) {
	// Load session from shared config
	if ses == nil {
		ses = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	}

	// Create new ECS client and list services
	ecsSvc := ecs.New(ses)
	listResult, listErr := ecsSvc.ListTasks(&ecs.ListTasksInput{
		Cluster:     aws.String(cluster),
		ServiceName: aws.String(serviceName),
	})
	var formattedTasks []ECSTask
	if listErr != nil {
		return formattedTasks, listErr
	}

	// Retrieve ECS tasks details
	describeResult, describeErr := ecsSvc.DescribeTasks(&ecs.DescribeTasksInput{
		Cluster: aws.String(cluster),
		Tasks:   listResult.TaskArns,
	})
	if describeErr != nil {
		return formattedTasks, describeErr
	}

	var instanceArns []*string
	var taskDefinitionArns []*string
	for _, task := range describeResult.Tasks {
		instanceArns = append(instanceArns, task.ContainerInstanceArn)
		taskDefinitionArns = append(taskDefinitionArns, task.TaskDefinitionArn)
	}

	// Retrieve ECS cluster instances details
	describeInstancesResult, describeInstancesErr := ecsSvc.DescribeContainerInstances(&ecs.DescribeContainerInstancesInput{
		Cluster:            aws.String(cluster),
		ContainerInstances: instanceArns,
	})
	if describeInstancesErr != nil {
		return formattedTasks, describeInstancesErr
	}

	containerInstances := make(map[string]*ecs.ContainerInstance)
	for _, containerInstance := range describeInstancesResult.ContainerInstances {
		containerInstances[*containerInstance.ContainerInstanceArn] = containerInstance
	}

	// Enrich task
	taskDefinitions := make(map[string]*ecs.TaskDefinition)
	for _, task := range describeResult.Tasks {
		if taskDefinitions[*task.TaskDefinitionArn] == nil {
			// Retrieve ECS task definition details
			describeTaskDefinitionResult, describeTaskDefinitionErr := ecsSvc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
				TaskDefinition: task.TaskDefinitionArn,
			})
			if describeTaskDefinitionErr != nil {
				return formattedTasks, describeTaskDefinitionErr
			}
			taskDefinitions[*task.TaskDefinitionArn] = describeTaskDefinitionResult.TaskDefinition
		}
		formattedTasks = append(formattedTasks, ECSTask{
			Arn:               task.TaskArn,
			ContainerInstance: containerInstances[*task.ContainerInstanceArn],
			HealthStatus:      task.HealthStatus,
			Status:            task.LastStatus,
			TaskDefinition:    taskDefinitions[*task.TaskDefinitionArn],
		})
	}

	return formattedTasks, nil
}

// ECSListTaskContainers return a list of containers
func ECSListTaskContainers(ses *session.Session, cluster string, taskArn string) ([]*ecs.Container, error) {
	// Load session from shared config
	if ses == nil {
		ses = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	}

	// Retrieve ECS services details
	ecsSvc := ecs.New(ses)
	describeResult, describeErr := ecsSvc.DescribeTasks(&ecs.DescribeTasksInput{
		Cluster: aws.String(cluster),
		Tasks: []*string{
			aws.String(taskArn),
		},
	})
	if describeErr != nil {
		return make([]*ecs.Container, 0), describeErr
	}

	return describeResult.Tasks[0].Containers, nil
}