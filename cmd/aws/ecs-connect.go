package aws

import (
	"awslib"
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
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
	currentSession := CreateAWSSession(c)

	// List available service
	services, err := awslib.ECSListServices(currentSession, c.String("cluster"))
	if err != nil {
		return cli.Exit("Error during ECS services list: "+err.Error(), -1)
	}
	if len(services) == 0 {
		return cli.Exit("No services found", -1)
	}

	// Build table
	header := fmt.Sprintf(
		"%-40s\t%-6s\t%-6s",
		"Name", "Desired", "Running",
	)
	var options []string
	for _, service := range services {
		options = append(options, fmt.Sprintf(
			"%-40s\t%-6s\t%-6s",
			*service.Name, strconv.FormatInt(*service.DesiredCount, 10), strconv.FormatInt(*service.RunningCount, 10),
		))
	}
	// Ask selection
	serviceSelectedIndex := -1
	prompt := &survey.Select{
		Message:  "Select a service: \n\n  " + header + "\n",
		Options:  options,
		PageSize: 15,
	}
	survey.AskOne(prompt, &serviceSelectedIndex)
	fmt.Println("")

	// Check response
	if serviceSelectedIndex == -1 {
		fmt.Println("\nNo service selected")
		return nil
	}

	// Set service in context
	c.Set("service", *services[serviceSelectedIndex].Name)
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
	currentSession := CreateAWSSession(c)

	// List available service
	tasks, err := awslib.ECSListServiceTasks(currentSession, c.String("cluster"), c.String("service"))
	if err != nil {
		return cli.Exit("Error during ECS tasks list: "+err.Error(), -1)
	}
	if len(tasks) == 0 {
		return cli.Exit("No tasks found for this service", -1)
	}

	// Build table
	header := fmt.Sprintf(
		"%-35s\t%-8s\t%-8s\t%-8s\t%s",
		"Task Definition Family", "Revision",
		"Status", "Health", "Instance ID",
	)
	var options []string
	for _, task := range tasks {
		options = append(options, fmt.Sprintf(
			"%-35s\t%-8s\t%-8s\t%-8s\t%s",
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
	fmt.Println("")

	// Check response
	if taskSelectedIndex == -1 {
		fmt.Println("\nNo task selected")
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
	currentSession := CreateAWSSession(c)

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
		"%-35s\t%-12s\t%8s\t%s",
		"Container Name", "Container ID", "Status", "Image",
	)
	var options []string
	for _, container := range containers {
		options = append(options, fmt.Sprintf(
			"%-35s\t%-12s\t%8s\t%s",
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
	fmt.Println("")

	// Check response
	if containerSelectedIndex == -1 {
		fmt.Println("\nNo container selected")
		return nil
	}

	// Set container ID
	c.Set("container", *containers[containerSelectedIndex].RuntimeId)
	return ECSConnectToContainer(c)
}

// ECSConnectToContainer connect to select container
func ECSConnectToContainer(c *cli.Context) error {
	dockerExecCmd := "sudo docker exec"

	user := c.String("user")
	if len(user) > 0 {
		dockerExecCmd += fmt.Sprintf(" --user %s", user)
	}

	workdir := c.String("workdir")
	if len(workdir) > 0 {
		dockerExecCmd += fmt.Sprintf(" --workdir %s", workdir)
	}

	c.Set("command", fmt.Sprintf("%s -it %s %s", dockerExecCmd, c.String("container"), c.String("command")))
	return SSMStartSession(c)
}
