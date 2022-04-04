package aws

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// ECSCluster is a result of ECSListClusters
type ECSCluster struct {
	Name *string
}

// ECSListClusters return a list of cluster
func ECSListClusters(ses *session.Session) ([]ECSCluster, error) {
	// Load session from shared config
	if ses == nil {
		ses = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	}

	// Create new ECS client and list services
	ecsSvc := ecs.New(ses)
	listResult, listErr := ecsSvc.ListClusters(&ecs.ListClustersInput{
		MaxResults: aws.Int64(100),
	})
	var formattedClusters []ECSCluster
	if listErr != nil {
		return formattedClusters, listErr
	}

	if len(listResult.ClusterArns) == 0 {
		return formattedClusters, nil
	}

	// Collect clusters ARNs
	var clusterArns []*string
	for _, clusterArn := range listResult.ClusterArns {
		clusterArns = append(clusterArns, clusterArn)
	}

	// Retrieve ECS clusters details
	var clusters []*ecs.Cluster
	chunkSize := 10
	for i := 0; i < len(clusterArns); i += chunkSize {
		end := i + chunkSize

		if end > len(clusterArns) {
			end = len(clusterArns)
		}

		describeResult, describeErr := ecsSvc.DescribeClusters(&ecs.DescribeClustersInput{
			Clusters: clusterArns[i:end],
		})
		if describeErr != nil {
			return formattedClusters, describeErr
		}
		clusters = append(clusters, describeResult.Clusters...)
	}

	// Format Clusters
	for _, cluster := range clusters {
		formattedClusters = append(formattedClusters, ECSCluster{
			Name: cluster.ClusterName,
		})
	}

	return formattedClusters, nil
}

// ECSService is a result of EC2ListInstances
type ECSService struct {
	Name         *string
	DesiredCount *int64
	PendingCount *int64
	RunningCount *int64
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
	serviceArns := make([]*string, 0)
	ecsSvc.ListServicesPages(&ecs.ListServicesInput{
		Cluster:    aws.String(cluster),
		MaxResults: aws.Int64(100),
	}, func(page *ecs.ListServicesOutput, lastPage bool) bool {
		for _, serviceArn := range page.ServiceArns {
			serviceArns = append(serviceArns, serviceArn)
		}
		return true // iterate over all pages
	})

	// Setup service channel
	errs := make(chan error, len(serviceArns))
	services := make(chan *ecs.Service, len(serviceArns))

	// Setup wait group
	var waitGroup sync.WaitGroup

	// Retrieve ECS services details
	chunkSize := 10
	for i := 0; i < len(serviceArns); i += chunkSize {
		// Prepare chunk
		end := i + chunkSize
		if end > len(serviceArns) {
			end = len(serviceArns)
		}
		serviceArnsChunk := serviceArns[i:end]

		// Execute parallel deploy
		waitGroup.Add(1)
		go func(sac []*string) {
			defer waitGroup.Done()

			res, err := ecsSvc.DescribeServices(&ecs.DescribeServicesInput{
				Cluster:  aws.String(cluster),
				Services: sac,
			})

			errs <- err

			for _, s := range res.Services {
				services <- s
			}

		}(serviceArnsChunk)
	}

	// Wait until all requests end
	waitGroup.Wait()

	// Close channels
	close(errs)
	close(services)

	// Check error
	for i := 0; i < len(serviceArns); i++ {
		err := <-errs
		if err != nil {
			return nil, err
		}
	}

	// Format service
	var formattedServices []ECSService
	for i := 0; i < len(serviceArns); i++ {
		service := <-services
		formattedServices = append(formattedServices, ECSService{
			Name:         service.ServiceName,
			DesiredCount: service.DesiredCount,
			PendingCount: service.PendingCount,
			RunningCount: service.RunningCount,
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

	if len(listResult.TaskArns) == 0 {
		return formattedTasks, nil
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

	containerInstances := make(map[string]*ecs.ContainerInstance)
	if len(instanceArns) > 0 && instanceArns[0] != nil {
		// Retrieve ECS cluster instances details
		describeInstancesResult, describeInstancesErr := ecsSvc.DescribeContainerInstances(&ecs.DescribeContainerInstancesInput{
			Cluster:            aws.String(cluster),
			ContainerInstances: instanceArns,
		})
		if describeInstancesErr != nil {
			return formattedTasks, describeInstancesErr
		}

		// Elaborate instance map
		for _, containerInstance := range describeInstancesResult.ContainerInstances {
			containerInstances[*containerInstance.ContainerInstanceArn] = containerInstance
		}
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

		var containerInstance *ecs.ContainerInstance
		if task.ContainerInstanceArn != nil {
			containerInstance = containerInstances[*task.ContainerInstanceArn]
		}

		formattedTasks = append(formattedTasks, ECSTask{
			Arn:               task.TaskArn,
			ContainerInstance: containerInstance,
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
