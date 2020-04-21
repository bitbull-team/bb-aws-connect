package aws

import (
	"awslib"
	"fmt"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/urfave/cli/v2"
)

// ECSListServices list ECS Services
func ECSListServices(c *cli.Context) error {
	// Check if service name is provided
	serviceName := c.String("service")
	if len(serviceName) != 0 {
		// Start SSM session
		return ECSListTasks(c)
	}

	// Create AWS session
	currentSession := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           c.String("profile"),
		SharedConfigState: session.SharedConfigEnable,
	}))

	// List available service
	services, err := awslib.ECSListServices(currentSession, c.String("cluster"))
	if err != nil {
		return cli.Exit("Error during ECS services list: "+err.Error(), -1)
	}
	if len(services) == 0 {
		return cli.Exit("No services found", -1)
	}

	// Build table
	var options []string
	for _, service := range services {
		options = append(options, fmt.Sprintf("%s", service.Name))
	}

	// Ask selection
	serviceSelected := ""
	prompt := &survey.Select{
		Message:  "Select a service:",
		Options:  options,
		PageSize: 15,
	}
	survey.AskOne(prompt, &serviceSelected)

	// Check response
	serviceName = strings.Split(serviceSelected, "\t")[0]
	if len(serviceName) == 0 {
		fmt.Println("No service selected")
		return nil
	}

	// Set service in context
	c.Set("service", serviceName)
	return ECSListTasks(c)
}

// ECSListTasks list ECS Tasks
func ECSListTasks(c *cli.Context) error {
	// Check if task name is provided
	taskID := c.String("task")
	if len(taskID) != 0 {
		// Start SSM session
		return ECSListContainer(c)
	}

	// Create AWS session
	currentSession := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           c.String("profile"),
		SharedConfigState: session.SharedConfigEnable,
	}))

	// List available service
	tasks, err := awslib.ECSListServiceTasks(currentSession, c.String("cluster"), c.String("service"))
	if err != nil {
		return cli.Exit("Error during ECS tasks list: "+err.Error(), -1)
	}
	if len(tasks) == 0 {
		return cli.Exit("No tasks found", -1)
	}

	// Build table
	header := fmt.Sprintf(
		"%s\t%s\t%s\t%s\t%s",
		"Task Definition Family          ", "Revision",
		"Status  ", "Health  ", "Instance ID",
	)
	var options []string
	for _, task := range tasks {
		options = append(options, fmt.Sprintf(
			"%s\t%s         \t%s  \t%s  \t%s",
			*task.TaskDefinition.Family, strconv.FormatInt(*task.TaskDefinition.Revision, 10),
			*task.Status, *task.HealthStatus, *task.ContainerInstance.Ec2InstanceId,
		))
	}

	// Ask selection
	taskSelectedIndex := -1
	prompt := &survey.Select{
		Message:  "Select a task: \n\n  " + header + "\n",
		Options:  options,
		PageSize: 15,
	}
	survey.AskOne(prompt, &taskSelectedIndex)

	// Check response
	if taskSelectedIndex == -1 {
		fmt.Println("No task selected")
		return nil
	}

	// Set task and instance
	c.Set("task", *tasks[taskSelectedIndex].Arn)
	c.Set("instance", *tasks[taskSelectedIndex].ContainerInstance.Ec2InstanceId)
	return ECSListContainer(c)
}

// ECSListContainer list ECS Tasks containers
func ECSListContainer(c *cli.Context) error {
	// Check if container id is provided
	containerID := c.String("container")
	if len(containerID) != 0 {
		// Start SSM session
		return ECSConnectToContainer(c)
	}

	// Create AWS session
	currentSession := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           c.String("profile"),
		SharedConfigState: session.SharedConfigEnable,
	}))

	// List availables service
	containers, err := awslib.ECSListTaskContainers(currentSession, c.String("cluster"), c.String("task"))
	if err != nil {
		return cli.Exit("Error during ECS task containers list: "+err.Error(), -1)
	}
	if len(containers) == 0 {
		return cli.Exit("No containers found", -1)
	}

	// Build table
	header := fmt.Sprintf(
		"%s\t%s\t%s\t%s",
		"Container Name          ", "Container ID", "Status", "Image",
	)
	var options []string
	for _, container := range containers {
		options = append(options, fmt.Sprintf(
			"%s\t%s\t%s\t%s",
			*container.Name, string(*container.RuntimeId)[0:12], *container.LastStatus, *container.Image,
		))
	}

	// Ask selection
	containerSelectedIndex := -1
	prompt := &survey.Select{
		Message:  "Select a container: \n\n  " + header + "\n",
		Options:  options,
		PageSize: 15,
	}
	survey.AskOne(prompt, &containerSelectedIndex)

	// Check response
	if containerSelectedIndex == -1 {
		fmt.Println("No container selected")
		return nil
	}

	// Set container ID
	c.Set("container", *containers[containerSelectedIndex].RuntimeId)
	return ECSConnectToContainer(c)
}

// ECSConnectToContainer connect to select container
func ECSConnectToContainer(c *cli.Context) error {
	c.Set("command", fmt.Sprintf("sudo docker exec -it %s /bin/bash", c.String("container")))
	return SSMStartSession(c)
}
