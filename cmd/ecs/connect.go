package ecs

import (
	"awslib"
	"fmt"
	"ssm"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"
)

// NewConnectCommand return "ecs:connect" command
func NewConnectCommand(globalFlags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:   "connect",
		Usage:  "Connect to an ECS Task container",
		Action: Connect,
		Flags: append(globalFlags, []cli.Flag{
			&cli.StringFlag{
				Name:    "cluster",
				Aliases: []string{"c"},
				Usage:   "Cluster Name",
			},
			&cli.StringFlag{
				Name:    "service",
				Aliases: []string{"s"},
				Usage:   "Service name (example: my-service)",
			},
			&cli.StringFlag{
				Name:    "task",
				Aliases: []string{"t"},
				Usage:   "Task ID (example: xxxxxxxxxxxxxxxxxxxx)",
			},
			&cli.StringFlag{
				Name:   "container",
				Hidden: true,
			},
			&cli.StringFlag{
				Name:   "instance",
				Hidden: true,
			},
			&cli.StringFlag{
				Name:    "workdir",
				Aliases: []string{"w"},
				Usage:   "Docker exec 'workdir' parameters (example: /app)",
			},
			&cli.StringFlag{
				Name:    "user",
				Aliases: []string{"u"},
				Usage:   "Docker exec 'user' parameters (example: www-data)",
			},
			&cli.StringFlag{
				Name:  "command",
				Usage: "Use a custom command as entrypoint",
			},
		}...),
	}
}

// Connect connect to an ECS container
func Connect(c *cli.Context) error {
	var err error
	// List ECS clusters
	err = ListClusters(c)
	if err != nil {
		return err
	}

	// List ECS services
	err = ListServices(c)
	if err != nil {
		return err
	}

	// List ECS tasks
	err = ListTasks(c)
	if err != nil {
		return err
	}

	// List ECS container
	err = ListContainer(c)
	if err != nil {
		return err
	}

	// connect to ECS container
	err = ConnectToContainer(c)
	if err != nil {
		return err
	}

	return nil
}

// ListClusters list ECS Clusters
func ListClusters(c *cli.Context) error {
	// Check if service name is provided
	clusterName := c.String("cluster")
	if len(clusterName) != 0 {
		// List services
		return nil
	}

	// Create AWS session
	currentSession := awslib.CreateAWSSession(c, awslib.Config{
		Profile: config.Profile,
		Region:  config.Region,
	})

	// Get cluster name
	cluster := c.String("cluster")
	if len(cluster) == 0 {
		cluster = config.ECS.Cluster
	}

	// List available clusters
	clusters, err := awslib.ECSListClusters(currentSession)
	if err != nil {
		return cli.Exit("Error during ECS clusters list: "+err.Error(), 1)
	}
	if len(clusters) == 0 {
		return cli.Exit("No cluster found", 1)
	}

	// If only one cluster is found select it
	if len(clusters) == 1 {
		fmt.Println("Cluster auto selected: ", *clusters[0].Name)
		c.Set("cluster", *clusters[0].Name)
		return nil
	}

	// Build options
	var options []string
	for _, cluster := range clusters {
		options = append(options, *cluster.Name)
	}
	// Ask selection
	var clusterSelected string
	prompt := &survey.Select{
		Message:  "Select a cluster:",
		Options:  options,
		PageSize: 10,
	}
	survey.AskOne(prompt, &clusterSelected)
	fmt.Println("")

	// Check response
	if len(clusterSelected) == 0 {
		return cli.Exit("No cluster selected", 1)
	}

	// Set service in context
	c.Set("cluster", clusterSelected)
	return nil
}

// ListServices list ECS Services
func ListServices(c *cli.Context) error {
	// Check if service name is provided
	serviceName := c.String("service")
	if len(serviceName) != 0 {
		// Start SSM session
		return nil
	}

	// Create AWS session
	currentSession := awslib.CreateAWSSession(c, awslib.Config{
		Profile: config.Profile,
		Region:  config.Region,
	})

	// Get cluster name
	cluster := c.String("cluster")
	if len(cluster) == 0 {
		cluster = config.ECS.Cluster
	}

	// List available service
	services, err := awslib.ECSListServices(currentSession, cluster)
	if err != nil {
		return cli.Exit("Error during ECS services list: "+err.Error(), 1)
	}
	if len(services) == 0 {
		return cli.Exit("No services found", 1)
	}

	// If only one services is found connect to it
	if len(services) == 1 {
		fmt.Println("Service auto selected: ", *services[0].Name)
		c.Set("service", *services[0].Name)
		return nil
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
		return cli.Exit("No service selected", 1)
	}

	// Check task status
	if *services[serviceSelectedIndex].RunningCount == 0 {
		fmt.Println("Selected service has no running tasks, cannot connect")
		return nil
	}

	// Set service in context
	c.Set("service", *services[serviceSelectedIndex].Name)
	return nil
}

// ListTasks list ECS Tasks
func ListTasks(c *cli.Context) error {
	// Check if task name is provided
	taskID := c.String("task")
	if len(taskID) != 0 {
		// Start SSM session
		return nil
	}

	// Create AWS session
	currentSession := awslib.CreateAWSSession(c, awslib.Config{
		Profile: config.Profile,
		Region:  config.Region,
	})

	// Get cluster name
	cluster := c.String("cluster")
	if len(cluster) == 0 {
		cluster = config.ECS.Cluster
	}

	// List available service
	tasks, err := awslib.ECSListServiceTasks(currentSession, cluster, c.String("service"))
	if err != nil {
		return cli.Exit("Error during ECS tasks list: "+err.Error(), 1)
	}
	if len(tasks) == 0 {
		return cli.Exit("No tasks found for this service", 1)
	}

	// If only one task is found connect to it
	if len(tasks) == 1 {
		fmt.Println("Task auto selected: ", *tasks[0].Arn)
		c.Set("task", *tasks[0].Arn)
		c.Set("instance", *tasks[0].ContainerInstance.Ec2InstanceId)
		return nil
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
		return cli.Exit("No task selected", 1)
	}

	// Check task status
	status := *tasks[taskSelectedIndex].Status
	if status != "RUNNING" {
		fmt.Println("Selected task is in ", status, " status, cannot connect")
		return nil
	}

	// Set task and instance
	c.Set("task", *tasks[taskSelectedIndex].Arn)
	c.Set("instance", *tasks[taskSelectedIndex].ContainerInstance.Ec2InstanceId)
	return nil
}

// ListContainer list ECS Tasks containers
func ListContainer(c *cli.Context) error {
	// Check if container id is provided
	containerID := c.String("container")
	if len(containerID) != 0 {
		// Start SSM session
		return nil
	}

	// Create AWS session
	currentSession := awslib.CreateAWSSession(c, awslib.Config{
		Profile: config.Profile,
		Region:  config.Region,
	})

	// Get cluster name
	cluster := c.String("cluster")
	if len(cluster) == 0 {
		cluster = config.ECS.Cluster
	}

	// List availables service
	containers, err := awslib.ECSListTaskContainers(currentSession, cluster, c.String("task"))
	if err != nil {
		return cli.Exit("Error during ECS task containers list: "+err.Error(), 1)
	}
	if len(containers) == 0 {
		return cli.Exit("No containers found", 1)
	}

	// If only one container is found connect to it
	if len(containers) == 1 {
		fmt.Println("Container auto selected: ", *containers[0].RuntimeId)
		c.Set("container", *containers[0].RuntimeId)
		return nil
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
		return cli.Exit("No container selected", 1)
	}

	// Set container ID
	c.Set("container", *containers[containerSelectedIndex].RuntimeId)
	return nil
}

// ConnectToContainer connect to select container
func ConnectToContainer(c *cli.Context) error {
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
	return ssm.StartSession(c)
}
